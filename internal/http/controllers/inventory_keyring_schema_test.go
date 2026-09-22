package controllers

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func TestInventoryKeyringSnapshotSchemaIncludesEmbeddedFields(t *testing.T) {
	parsed, err := schema.Parse(&inventoryKeyringRawSnapshot{}, &sync.Map{}, schema.NamingStrategy{})
	if err != nil {
		t.Fatal(err)
	}
	for _, column := range []string{"time_index", "character_id", "slot_id", "item_id", "id", "item_unique_id"} {
		if parsed.FieldsByDBName[column] == nil {
			t.Errorf("snapshot scan silently omits %s; mapped columns: %v", column, parsed.DBNames)
		}
	}
}

func TestInventoryKeyringSchemaCompatibility(t *testing.T) {
	for _, generation := range []string{"current", "legacy", "mixed"} {
		t.Run(generation, func(t *testing.T) {
			fixture := newInventoryKeyringSchemaFixture(generation)
			schema, err := loadInventoryKeyringSchema(fixture.open(t))
			if err != nil {
				t.Fatal(err)
			}
			character, slot, item := "character_id", "slot_id", "item_id"
			if generation == "legacy" {
				character, slot, item = "charid", "slotid", "itemid"
			}
			if schema.snapshotCharacterID != character || schema.snapshotSlotID != slot || schema.snapshotItemID != item {
				t.Fatalf("snapshot character/slot/item = %q/%q/%q, want %q/%q/%q", schema.snapshotCharacterID, schema.snapshotSlotID, schema.snapshotItemID, character, slot, item)
			}
			projection := normalizeInventoryKeyringSQL(schema.snapshotSelect)
			for _, fragment := range []string{
				"snap." + character + " AS character_id",
				"snap." + slot + " AS slot_id",
				"COALESCE(snap." + item + ", 0) AS item_id",
				"CONCAT('Unknown item #', snap." + item + ")",
			} {
				if !strings.Contains(projection, fragment) {
					t.Errorf("snapshot projection missing %q: %s", fragment, projection)
				}
			}
			for _, mapping := range inventoryKeyringSnapshotColumnMappings {
				source := mapping.current
				if generation == "legacy" {
					source = mapping.legacy
				}
				if !strings.Contains(projection, "snap."+source) {
					t.Errorf("snapshot projection missing %s", source)
				}
				if generation != "legacy" && strings.Contains(projection, "snap."+mapping.legacy) {
					t.Errorf("current snapshot projection still references legacy %s", mapping.legacy)
				}
			}
			for alias, selected := range map[string]string{"inv": schema.inventorySelect, "sb": schema.sharedBankSelect, "snap": schema.snapshotSelect} {
				selected = normalizeInventoryKeyringSQL(selected)
				guid, uniqueID := "0 AS guid", "COALESCE("+alias+".item_unique_id, '') AS item_unique_id"
				if generation != "current" {
					guid = "COALESCE(" + alias + ".guid, 0) AS guid"
				}
				if generation == "legacy" {
					uniqueID = "'' AS item_unique_id"
				}
				for _, fragment := range []string{guid, uniqueID} {
					if !strings.Contains(selected, fragment) {
						t.Errorf("%s identity projection missing %q: %s", alias, fragment, selected)
					}
				}
				if generation == "current" && strings.Contains(selected, alias+".guid") {
					t.Errorf("current %s projection references removed guid column", alias)
				}
			}
		})
	}
}

func TestInventoryKeyringSchemaOptionalIdentity(t *testing.T) {
	fixture := newInventoryKeyringSchemaFixture("current")
	for table := range fixture.columns {
		fixture.removeColumn(table, "item_unique_id")
	}
	schema, err := loadInventoryKeyringSchema(fixture.open(t))
	if err != nil {
		t.Fatal(err)
	}
	for _, selected := range []string{schema.inventorySelect, schema.sharedBankSelect, schema.snapshotSelect} {
		selected = normalizeInventoryKeyringSQL(selected)
		if !strings.Contains(selected, "0 AS guid") || !strings.Contains(selected, "'' AS item_unique_id") {
			t.Fatalf("missing optional identity defaults: %s", selected)
		}
	}
}

func TestInventoryKeyringSchemaIdentityIsResolvedPerTable(t *testing.T) {
	fixture := newInventoryKeyringSchemaFixture("current")
	fixture.removeColumn("sharedbank", "item_unique_id")
	fixture.columns["sharedbank"] = append(fixture.columns["sharedbank"], "guid")
	fixture.removeColumn("inventory_snapshots", "item_unique_id")
	schema, err := loadInventoryKeyringSchema(fixture.open(t))
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct{ projection, guid, uniqueID string }{
		{schema.inventorySelect, "0 AS guid", "COALESCE(inv.item_unique_id, '') AS item_unique_id"},
		{schema.sharedBankSelect, "COALESCE(sb.guid, 0) AS guid", "'' AS item_unique_id"},
		{schema.snapshotSelect, "0 AS guid", "'' AS item_unique_id"},
	} {
		projection := normalizeInventoryKeyringSQL(test.projection)
		if !strings.Contains(projection, test.guid) || !strings.Contains(projection, test.uniqueID) {
			t.Fatalf("identity columns from another table were reused: %s", projection)
		}
	}
}

func TestInventoryKeyringSchemaRejectsIncompleteMetadata(t *testing.T) {
	for _, scenario := range []string{"missing snapshot slot", "missing snapshot table", "metadata failure"} {
		t.Run(scenario, func(t *testing.T) {
			fixture := newInventoryKeyringSchemaFixture("current")
			switch scenario {
			case "missing snapshot slot":
				fixture.removeColumn("inventory_snapshots", "slot_id")
			case "missing snapshot table":
				delete(fixture.columns, "inventory_snapshots")
			case "metadata failure":
				fixture.queryErr = errors.New("metadata access denied")
			}
			_, err := loadInventoryKeyringSchema(fixture.open(t))
			if err == nil {
				t.Fatal("incomplete or unavailable schema unexpectedly accepted")
			}
			if fixture.queryErr != nil && !errors.Is(err, fixture.queryErr) {
				t.Fatalf("metadata error was not preserved: %v", err)
			}
		})
	}
}

func TestInventoryKeyringHydrationPreservesTextIdentity(t *testing.T) {
	for _, identity := range []string{"000012345678901234567890", "peq:4f823d6a-f92a-4f97-a430-815e12b02c6d"} {
		t.Run(identity, func(t *testing.T) {
			row := InventoryKeyringRawInventory{CharacterID: 42, SlotID: 0, ItemID: 100, GUID: 123, ItemUniqueID: identity}
			inventory, err := hydrateInventoryKeyringRows(nil, []InventoryKeyringRawInventory{row})
			if err != nil {
				t.Fatal(err)
			}
			if len(inventory) != 1 || inventory[0].ItemUniqueID != identity || inventory[0].GUID != 123 {
				t.Fatalf("inventory hydration altered identity: %#v", inventory)
			}
			snapshots, err := hydrateInventoryKeyringSnapshots(nil, []inventoryKeyringRawSnapshot{{TimeIndex: 12345, InventoryKeyringRawInventory: row}})
			if err != nil {
				t.Fatal(err)
			}
			if len(snapshots) != 1 || snapshots[0].ItemUniqueID != identity || snapshots[0].GUID != 123 {
				t.Fatalf("snapshot hydration altered identity: %#v", snapshots)
			}
			for _, record := range []interface{}{inventory[0], snapshots[0]} {
				encoded, err := json.Marshal(record)
				if err != nil {
					t.Fatal(err)
				}
				if !strings.Contains(string(encoded), `"item_unique_id":"`+identity+`"`) {
					t.Fatalf("JSON identity was not preserved as text: %s", encoded)
				}
			}
		})
	}
}

var inventoryKeyringSnapshotColumnMappings = []struct{ current, legacy string }{
	{"character_id", "charid"}, {"slot_id", "slotid"}, {"item_id", "itemid"},
	{"augment_one", "augslot1"}, {"augment_two", "augslot2"}, {"augment_three", "augslot3"},
	{"augment_four", "augslot4"}, {"augment_five", "augslot5"}, {"augment_six", "augslot6"},
	{"ornament_icon", "ornamenticon"}, {"ornament_idfile", "ornamentidfile"},
}

func normalizeInventoryKeyringSQL(query string) string {
	return strings.Join(strings.Fields(strings.ReplaceAll(query, "`", "")), " ")
}

// The fixture permits only database-scoped metadata reads. Unexpected SQL or
// any write fails, so schema discovery cannot silently fall back after errors.
type inventoryKeyringSchemaFixture struct {
	columns  map[string][]string
	queryErr error
}

func newInventoryKeyringSchemaFixture(generation string) *inventoryKeyringSchemaFixture {
	fixture := &inventoryKeyringSchemaFixture{columns: map[string][]string{}}
	for _, table := range []string{"inventory", "sharedbank", "inventory_snapshots"} {
		fixture.columns[table] = []string{"charges", "color", "instnodrop", "custom_data", "ornament_hero_model"}
		if generation != "legacy" {
			fixture.columns[table] = append(fixture.columns[table], "item_unique_id")
		}
		if generation != "current" {
			fixture.columns[table] = append(fixture.columns[table], "guid")
		}
		for _, mapping := range inventoryKeyringSnapshotColumnMappings {
			if table != "inventory_snapshots" || generation != "legacy" {
				fixture.columns[table] = append(fixture.columns[table], mapping.current)
			}
			if table == "inventory_snapshots" && generation != "current" {
				fixture.columns[table] = append(fixture.columns[table], mapping.legacy)
			}
		}
	}
	fixture.columns["sharedbank"] = append(fixture.columns["sharedbank"], "account_id")
	fixture.columns["inventory_snapshots"] = append(fixture.columns["inventory_snapshots"], "time_index")
	return fixture
}

func (f *inventoryKeyringSchemaFixture) removeColumn(table, column string) {
	columns := f.columns[table][:0]
	for _, name := range f.columns[table] {
		if name != column {
			columns = append(columns, name)
		}
	}
	f.columns[table] = columns
}

func (f *inventoryKeyringSchemaFixture) open(t *testing.T) *gorm.DB {
	t.Helper()
	pool := sql.OpenDB(f)
	t.Cleanup(func() { _ = pool.Close() })
	db, err := gorm.Open(mysql.New(mysql.Config{Conn: pool, SkipInitializeWithVersion: true}), &gorm.Config{DisableAutomaticPing: true})
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func (f *inventoryKeyringSchemaFixture) Connect(context.Context) (driver.Conn, error) { return f, nil }
func (f *inventoryKeyringSchemaFixture) Driver() driver.Driver                        { return f }
func (f *inventoryKeyringSchemaFixture) Open(string) (driver.Conn, error)             { return f, nil }
func (*inventoryKeyringSchemaFixture) Close() error                                   { return nil }
func (*inventoryKeyringSchemaFixture) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("unexpected prepared statement")
}
func (*inventoryKeyringSchemaFixture) Begin() (driver.Tx, error) {
	return nil, errors.New("unexpected transaction")
}
func (f *inventoryKeyringSchemaFixture) QueryContext(_ context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	normalized := strings.ToLower(normalizeInventoryKeyringSQL(query))
	if !strings.Contains(normalized, "information_schema.columns") || !strings.Contains(normalized, "table_schema = database()") {
		return nil, errors.New("unexpected or unscoped schema query: " + query)
	}
	for _, table := range []string{"inventory", "sharedbank", "inventory_snapshots"} {
		if !strings.Contains(normalized+fmt.Sprint(args), table) {
			return nil, errors.New("schema query omitted table " + table)
		}
	}
	if f.queryErr != nil {
		return nil, f.queryErr
	}
	rows := &inventoryKeyringSchemaRows{}
	for table, columns := range f.columns {
		for _, column := range columns {
			rows.values = append(rows.values, []driver.Value{table, column})
		}
	}
	return rows, nil
}

type inventoryKeyringSchemaRows struct{ values [][]driver.Value }

func (*inventoryKeyringSchemaRows) Columns() []string { return []string{"table_name", "column_name"} }
func (*inventoryKeyringSchemaRows) Close() error      { return nil }
func (r *inventoryKeyringSchemaRows) Next(dest []driver.Value) error {
	if len(r.values) == 0 {
		return io.EOF
	}
	copy(dest, r.values[0])
	r.values = r.values[1:]
	return nil
}
