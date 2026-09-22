package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/EQEmuTools/spire/internal/database"
	"github.com/labstack/echo/v4"
	gocache "github.com/patrickmn/go-cache"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

// Run against either an existing modern or legacy EQEmu database by setting
// SPIRE_INVENTORY_TEST_DSN. This suite only reads existing data; it never creates
// fixtures or runs migrations. A database account with SELECT access is enough.
func TestInventoryKeyringDatabaseIntegration(t *testing.T) {
	dsn := os.Getenv("SPIRE_INVENTORY_TEST_DSN")
	if dsn == "" {
		t.Skip("set SPIRE_INVENTORY_TEST_DSN to run read-only database coverage")
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 gormLogger.Default.LogMode(gormLogger.Silent),
	})
	if err != nil {
		t.Fatal("could not open integration test database")
	}
	pool, err := db.DB()
	if err != nil {
		t.Fatal("could not obtain integration test connection")
	}
	t.Cleanup(func() { _ = pool.Close() })
	resolver := database.NewResolver(database.NewConnections(db, db, nil), nil, nil, gocache.New(gocache.NoExpiration, 0))
	controller := NewInventoryKeyringController(resolver, nil)
	e := echo.New()
	for _, route := range controller.Routes() {
		if route.Method() == http.MethodGet {
			e.Add(route.Method(), "/"+route.Route(), route.Handler())
		}
	}
	get := func(t *testing.T, path string, destination interface{}, fields ...string) {
		t.Helper()
		response := httptest.NewRecorder()
		e.ServeHTTP(response, httptest.NewRequest(http.MethodGet, path, nil))
		if response.Code != http.StatusOK {
			// Do not print response bodies: they may contain player data.
			t.Fatalf("GET returned HTTP %d, want 200", response.Code)
		}
		var payload map[string]json.RawMessage
		if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
			t.Fatal("GET response is not a JSON object")
		}
		for _, field := range fields {
			if _, found := payload[field]; !found {
				t.Fatalf("GET response is missing %s", field)
			}
		}
		if err := json.Unmarshal(response.Body.Bytes(), destination); err != nil {
			t.Fatal("GET response does not match the expected JSON types")
		}
	}
	schema, err := loadInventoryKeyringSchema(db)
	if err != nil {
		t.Fatalf("load inventory schema: %v", err)
	}
	t.Run("query_columns", func(t *testing.T) {
		for _, query := range []struct{ table, selection, join string }{
			{"inventory inv", schema.inventorySelect, "LEFT JOIN items item ON item.id = inv.item_id"},
			{"sharedbank sb", schema.sharedBankSelect, "LEFT JOIN items item ON item.id = sb.item_id"},
			{"inventory_snapshots snap", schema.snapshotSelect, "LEFT JOIN items item ON item.id = snap." + schema.snapshotItemID},
		} {
			rows, err := db.Table(query.table).Select(query.selection).Joins(query.join).Limit(0).Rows()
			if err != nil {
				t.Fatalf("validate %s projection: %v", query.table, err)
			}
			_ = rows.Close()
		}
	})
	t.Run("summary", func(t *testing.T) {
		var summary inventoryKeyringSummary
		get(t, "/inventory-keyring/summary", &summary,
			"characters_with_inventory", "inventory_items", "keyring_characters", "keyring_entries", "snapshot_characters", "snapshot_sets")
		var snapshotCharacters, snapshotSets int64
		if err := db.Raw("SELECT COUNT(DISTINCT " + schema.snapshotCharacterID + ") FROM inventory_snapshots").Scan(&snapshotCharacters).Error; err != nil {
			t.Fatalf("count snapshot characters: %v", err)
		}
		if err := db.Raw("SELECT COUNT(*) FROM (SELECT " + schema.snapshotCharacterID + ", time_index FROM inventory_snapshots GROUP BY " + schema.snapshotCharacterID + ", time_index) snapshots").Scan(&snapshotSets).Error; err != nil {
			t.Fatalf("count snapshot sets: %v", err)
		}
		if summary.SnapshotCharacters != snapshotCharacters || summary.SnapshotSets != snapshotSets {
			t.Fatal("summary snapshot counts do not match the database")
		}
	})
	for _, state := range []string{"all", "inventory", "keyring", "snapshots", "empty"} {
		t.Run("directory_"+state, func(t *testing.T) {
			var page struct {
				Data  []inventoryKeyringCharacterSummary `json:"data"`
				Total int64                              `json:"total"`
				Page  int                                `json:"page"`
				Limit int                                `json:"limit"`
			}
			get(t, "/inventory-keyring/characters?page=1&limit=30&state="+state, &page, "data", "total", "page", "limit")
			if page.Data == nil || page.Page != 1 || page.Limit != 30 || len(page.Data) > 30 || page.Total < int64(len(page.Data)) {
				t.Fatal("directory pagination is invalid")
			}
			for _, character := range page.Data {
				if character.ID <= 0 {
					t.Fatal("directory row has no character ID")
				}
				if (state == "inventory" && character.InventoryCount == 0) ||
					(state == "keyring" && character.KeyCount == 0) ||
					(state == "snapshots" && character.SnapshotCount == 0) ||
					(state == "empty" && (character.InventoryCount != 0 || character.KeyCount != 0)) {
					t.Fatal("directory returned a row that does not match its state filter")
				}
			}
		})
	}

	var fixture struct {
		ID        int `gorm:"column:id"`
		AccountID int `gorm:"column:account_id"`
	}
	base := db.Table("character_data ch").Select("ch.id, ch.account_id").Where("ch.deleted_at IS NULL")
	// Prefer a character that exercises both inventory and snapshot reads, and
	// shared bank when available. Sparse databases still cover the GET routes.
	if err := base.Session(&gorm.Session{}).
		Where("EXISTS (SELECT 1 FROM inventory inv WHERE inv.character_id = ch.id)").
		Where("EXISTS (SELECT 1 FROM inventory_snapshots snap WHERE snap." + schema.snapshotCharacterID + " = ch.id AND snap.time_index > 0)").
		Order("EXISTS (SELECT 1 FROM sharedbank sb WHERE sb.account_id = ch.account_id) DESC, ch.id").
		Limit(1).Scan(&fixture).Error; err != nil {
		t.Fatalf("select combined fixture: %v", err)
	}
	if fixture.ID == 0 {
		if err := base.Session(&gorm.Session{}).
			Order("EXISTS (SELECT 1 FROM inventory inv WHERE inv.character_id = ch.id) DESC, ch.id").
			Limit(1).Scan(&fixture).Error; err != nil {
			t.Fatalf("select fallback fixture: %v", err)
		}
	}
	t.Run("directory_search", func(t *testing.T) {
		query := "1"
		if fixture.ID > 0 {
			query = strconv.Itoa(fixture.ID)
		}
		var page struct {
			Data  []inventoryKeyringCharacterSummary `json:"data"`
			Total int64                              `json:"total"`
		}
		get(t, "/inventory-keyring/characters?page=1&limit=100&q="+query, &page, "data", "total", "page", "limit")
		if fixture.ID > 0 && page.Total == 0 {
			t.Fatal("search by an existing character ID returned no matches")
		}
	})
	t.Run("character_detail_and_record_loaders", func(t *testing.T) {
		if fixture.ID == 0 {
			t.Skip("database has no nondeleted character fixture")
		}
		var detail inventoryKeyringCharacterDetail
		get(t, fmt.Sprintf("/inventory-keyring/character/%d", fixture.ID), &detail,
			"character", "inventory", "keyring", "snapshots", "slots")
		if detail.Character.ID != fixture.ID || detail.Character.AccountID != fixture.AccountID || len(detail.Slots) == 0 {
			t.Fatal("character detail does not match its fixture")
		}
		if int64(len(detail.Inventory)) != detail.Character.InventoryCount || int64(len(detail.Keyring)) != detail.Character.KeyCount {
			t.Fatal("detail inventory or keyring count does not match its summary")
		}
		for _, storage := range []struct {
			name, table, ownerColumn string
			ownerID                  int
		}{
			{inventoryStorageCharacter, "inventory", "character_id", fixture.ID},
			{inventoryStorageSharedBank, "sharedbank", "account_id", fixture.AccountID},
		} {
			t.Run(storage.name, func(t *testing.T) {
				var records []inventoryKeyringRecord
				for _, record := range detail.Inventory {
					if record.StorageKind == storage.name {
						records = append(records, record)
					}
				}
				if len(records) == 0 {
					t.Skip("fixture has no items in this storage")
				}
				record := records[0]
				loaded, err := loadInventoryKeyringRecord(db, fixture.ID, record.SlotID)
				if err != nil {
					t.Fatalf("load existing inventory record: %v", err)
				}
				if loaded.CharacterID != fixture.ID || loaded.AccountID != fixture.AccountID || loaded.StorageKind != storage.name ||
					loaded.SlotID != record.SlotID || loaded.ItemID != record.ItemID || loaded.ItemUniqueID != record.ItemUniqueID {
					t.Fatal("single-record loader does not match character detail")
				}
				assertInventoryIntegrationUniqueIDs(t, db, storage.table, "slot_id", storage.ownerColumn+" = ?", []interface{}{storage.ownerID}, records)
			})
		}
	})
	t.Run("snapshot_detail", func(t *testing.T) {
		if fixture.ID == 0 {
			t.Skip("database has no nondeleted character fixture")
		}
		var timeIndex int64
		if err := db.Table("inventory_snapshots").Select("COALESCE(MAX(time_index), 0)").
			Where(schema.snapshotCharacterID+" = ?", fixture.ID).Scan(&timeIndex).Error; err != nil {
			t.Fatalf("select snapshot fixture: %v", err)
		}
		if timeIndex <= 0 {
			t.Skip("character fixture has no snapshot")
		}
		var snapshot struct {
			TimeIndex int64                            `json:"time_index"`
			Items     []inventoryKeyringSnapshotRecord `json:"items"`
		}
		get(t, fmt.Sprintf("/inventory-keyring/character/%d/snapshot/%d", fixture.ID, timeIndex), &snapshot, "time_index", "items")
		where := schema.snapshotCharacterID + " = ? AND time_index = ?"
		args := []interface{}{fixture.ID, timeIndex}
		var count int64
		if err := db.Table("inventory_snapshots").Where(where, args...).Count(&count).Error; err != nil {
			t.Fatalf("count snapshot items: %v", err)
		}
		if snapshot.TimeIndex != timeIndex || count == 0 || int64(len(snapshot.Items)) != count {
			t.Fatal("snapshot detail does not match its persisted set")
		}
		var expectedItems []struct {
			SlotID    int `gorm:"column:slot_id"`
			ItemID    int `gorm:"column:item_id"`
			CatalogID int `gorm:"column:catalog_id"`
		}
		if err := db.Table("inventory_snapshots snap").
			Select("snap."+schema.snapshotSlotID+" AS slot_id, COALESCE(snap."+schema.snapshotItemID+", 0) AS item_id, COALESCE(item.id, 0) AS catalog_id").
			Joins("LEFT JOIN items item ON item.id = snap."+schema.snapshotItemID).
			Where("snap."+schema.snapshotCharacterID+" = ? AND snap.time_index = ?", fixture.ID, timeIndex).
			Order("snap." + schema.snapshotSlotID).Scan(&expectedItems).Error; err != nil {
			t.Fatalf("read persisted snapshot items: %v", err)
		}
		records := make([]inventoryKeyringRecord, 0, len(snapshot.Items))
		for index, record := range snapshot.Items {
			if record.TimeIndex != timeIndex || len(record.Augments) != 6 {
				t.Fatal("snapshot row is missing its timestamp or augment sockets")
			}
			if index >= len(expectedItems) || record.SlotID != expectedItems[index].SlotID || record.ItemID != expectedItems[index].ItemID || record.Item.ID != expectedItems[index].CatalogID {
				t.Fatal("snapshot slot, instance item, or catalog item ID does not match the database")
			}
			records = append(records, inventoryKeyringRecord{SlotID: record.SlotID, ItemUniqueID: record.ItemUniqueID})
		}
		assertInventoryIntegrationUniqueIDs(t, db, "inventory_snapshots", schema.snapshotSlotID, where, args, records)
	})
}

func assertInventoryIntegrationUniqueIDs(t *testing.T, db *gorm.DB, table, slotColumn, where string, args []interface{}, records []inventoryKeyringRecord) {
	t.Helper()
	var hasColumn int64
	if err := db.Raw(`SELECT COUNT(*) FROM information_schema.COLUMNS
		WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = 'item_unique_id'`, table).Scan(&hasColumn).Error; err != nil {
		t.Fatalf("inspect item identity column: %v", err)
	}
	if hasColumn == 0 {
		return
	}
	var identities []struct {
		SlotID       int    `gorm:"column:slot_id"`
		ItemUniqueID string `gorm:"column:item_unique_id"`
	}
	if err := db.Table(table).Select(slotColumn+" AS slot_id, COALESCE(item_unique_id, '') AS item_unique_id").
		Where(where, args...).Scan(&identities).Error; err != nil {
		t.Fatalf("read persisted item identities: %v", err)
	}
	expected := make(map[int]string, len(identities))
	for _, identity := range identities {
		expected[identity.SlotID] = identity.ItemUniqueID
	}
	for _, record := range records {
		identity, found := expected[record.SlotID]
		if !found || identity != record.ItemUniqueID {
			t.Fatal("item_unique_id was not preserved as its original string")
		}
	}
}
