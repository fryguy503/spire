package controllers

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// Snapshot columns were renamed by EQEmu migration 9324. Resolve against the
// selected connection on each request, so switching servers cannot reuse another
// database's schema. Identifiers below come only from these fixed column names.
type inventoryKeyringSchema struct {
	snapshotCharacterID string
	snapshotSlotID      string
	snapshotItemID      string
	snapshotSelect      string
	inventorySelect     string
	sharedBankSelect    string
}

func loadInventoryKeyringSchema(db *gorm.DB) (inventoryKeyringSchema, error) {
	var columns []struct {
		TableName  string `gorm:"column:table_name"`
		ColumnName string `gorm:"column:column_name"`
	}
	if err := db.Raw(`SELECT TABLE_NAME AS table_name, COLUMN_NAME AS column_name
		FROM information_schema.COLUMNS
		WHERE TABLE_SCHEMA = DATABASE()
		AND TABLE_NAME IN ('inventory', 'sharedbank', 'inventory_snapshots')`).Scan(&columns).Error; err != nil {
		return inventoryKeyringSchema{}, err
	}
	tables := map[string]map[string]bool{}
	for _, column := range columns {
		if tables[column.TableName] == nil {
			tables[column.TableName] = map[string]bool{}
		}
		tables[column.TableName][column.ColumnName] = true
	}
	for _, table := range []string{"inventory", "sharedbank", "inventory_snapshots"} {
		if len(tables[table]) == 0 {
			return inventoryKeyringSchema{}, fmt.Errorf("inventory and keyring schema: missing table %s", table)
		}
	}

	snapshotColumns := map[string]string{}
	replacements := make([]string, 0, 22)
	for _, names := range [][2]string{
		{"character_id", "charid"}, {"slot_id", "slotid"}, {"item_id", "itemid"},
		{"augment_one", "augslot1"}, {"augment_two", "augslot2"}, {"augment_three", "augslot3"},
		{"augment_four", "augslot4"}, {"augment_five", "augslot5"}, {"augment_six", "augslot6"},
		{"ornament_icon", "ornamenticon"}, {"ornament_idfile", "ornamentidfile"},
	} {
		column := names[0]
		if !tables["inventory_snapshots"][column] {
			column = names[1]
			if !tables["inventory_snapshots"][column] {
				return inventoryKeyringSchema{}, fmt.Errorf("inventory and keyring schema: inventory_snapshots is missing %s or %s", names[0], names[1])
			}
		}
		snapshotColumns[names[0]] = column
		replacements = append(replacements, "snap."+names[1], "snap."+column)
	}

	// A modern item_unique_id is a string, not the old numeric GUID. Preserve
	// each separately, and never select an identity column absent from this table.
	identitySelect := func(selection, table, alias string) string {
		guid := "0 AS guid"
		if tables[table]["guid"] {
			guid = "COALESCE(" + alias + ".guid, 0) AS guid"
		}
		uniqueID := "'' AS item_unique_id"
		if tables[table]["item_unique_id"] {
			uniqueID = "COALESCE(" + alias + ".item_unique_id, '') AS item_unique_id"
		}
		return strings.Replace(selection, "COALESCE("+alias+".guid, 0) AS guid", guid+",\n\t"+uniqueID, 1)
	}
	return inventoryKeyringSchema{
		snapshotCharacterID: snapshotColumns["character_id"],
		snapshotSlotID:      snapshotColumns["slot_id"],
		snapshotItemID:      snapshotColumns["item_id"],
		snapshotSelect:      identitySelect(strings.NewReplacer(replacements...).Replace(inventoryKeyringSnapshotSelect), "inventory_snapshots", "snap"),
		inventorySelect:     identitySelect(inventoryKeyringInventorySelect, "inventory", "inv"),
		sharedBankSelect:    identitySelect(inventoryKeyringSharedBankSelect, "sharedbank", "sb"),
	}, nil
}
