package models

import (
	"github.com/volatiletech/null/v8"
)

type Loottable struct {
	ID                     uint             `json:"id" gorm:"Column:id"`
	Name                   string           `json:"name" gorm:"Column:name"`
	Mincash                uint             `json:"mincash" gorm:"Column:mincash"`
	Maxcash                uint             `json:"maxcash" gorm:"Column:maxcash"`
	Avgcoin                uint             `json:"avgcoin" gorm:"Column:avgcoin"`
	Done                   int8             `json:"done" gorm:"Column:done"`
	MinExpansion           int8             `json:"min_expansion" gorm:"Column:min_expansion"`
	MaxExpansion           int8             `json:"max_expansion" gorm:"Column:max_expansion"`
	ContentFlags           null.String      `json:"content_flags" gorm:"Column:content_flags"`
	ContentFlagsDisabled   null.String      `json:"content_flags_disabled" gorm:"Column:content_flags_disabled"`
	LoottableEntries       []LoottableEntry `json:"loottable_entries,omitempty" gorm:"foreignKey:loottable_id;references:id"`
	NpcTypes               []NpcType        `json:"npc_types,omitempty" gorm:"foreignKey:loottable_id;references:id"`
}

func (Loottable) TableName() string {
    return "loottable"
}

func (Loottable) Relationships() []string {
    return []string{
		"LoottableEntries",
		"LoottableEntries.Lootdrop",
		"LoottableEntries.Lootdrop.LootdropEntries",
		"LoottableEntries.Lootdrop.LootdropEntries.Item",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.AlternateCurrencies",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.CharacterCorpseItems",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.DiscoveredItems",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Doors",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Doors.Item",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Fishings",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Fishings.Item",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Fishings.NpcType",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Fishings.NpcType.AlternateCurrency",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Fishings.NpcType.Loottable",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Fishings.NpcType.Merchantlists",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Fishings.NpcType.Merchantlists.Items",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Fishings.NpcType.Merchantlists.NpcTypes",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Fishings.NpcType.NpcEmotes",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Fishings.NpcType.NpcFactions",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Fishings.NpcType.NpcFactions.NpcFactionEntries",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Fishings.NpcType.NpcFactions.NpcFactionEntries.FactionList",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Fishings.NpcType.NpcSpell",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Fishings.NpcType.NpcSpell.NpcSpellsEntries",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Fishings.NpcType.NpcSpell.NpcSpellsEntries.SpellsNew",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Fishings.NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.Aura",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Fishings.NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.Aura.SpellsNew",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Fishings.NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.BlockedSpells",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Fishings.NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.Damageshieldtypes",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Fishings.NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.Items",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Fishings.NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.NpcSpellsEntries",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Fishings.NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.SpellBuckets",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Fishings.NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.SpellGlobals",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Fishings.NpcType.NpcTypesTint",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Fishings.NpcType.Spawnentries",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Fishings.NpcType.Spawnentries.NpcType",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Fishings.NpcType.Spawnentries.Spawngroup",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Fishings.NpcType.Spawnentries.Spawngroup.Spawn2",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Fishings.NpcType.Spawnentries.Spawngroup.Spawn2.Spawnentries",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Fishings.NpcType.Spawnentries.Spawngroup.Spawn2.Spawngroup",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Fishings.Zone",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Forages",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Forages.Item",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Forages.Zone",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.GroundSpawns",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.GroundSpawns.Zone",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.ItemTicks",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Keyrings",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.LootdropEntries",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Merchantlists",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Merchantlists.Items",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Merchantlists.NpcTypes",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Merchantlists.NpcTypes.AlternateCurrency",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Merchantlists.NpcTypes.Loottable",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Merchantlists.NpcTypes.Merchantlists",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Merchantlists.NpcTypes.NpcEmotes",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Merchantlists.NpcTypes.NpcFactions",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Merchantlists.NpcTypes.NpcFactions.NpcFactionEntries",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Merchantlists.NpcTypes.NpcFactions.NpcFactionEntries.FactionList",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Merchantlists.NpcTypes.NpcSpell",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Merchantlists.NpcTypes.NpcSpell.NpcSpellsEntries",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Merchantlists.NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Merchantlists.NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Aura",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Merchantlists.NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Aura.SpellsNew",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Merchantlists.NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.BlockedSpells",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Merchantlists.NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Damageshieldtypes",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Merchantlists.NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Items",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Merchantlists.NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.NpcSpellsEntries",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Merchantlists.NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.SpellBuckets",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Merchantlists.NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.SpellGlobals",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Merchantlists.NpcTypes.NpcTypesTint",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Merchantlists.NpcTypes.Spawnentries",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Merchantlists.NpcTypes.Spawnentries.NpcType",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Merchantlists.NpcTypes.Spawnentries.Spawngroup",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Merchantlists.NpcTypes.Spawnentries.Spawngroup.Spawn2",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Merchantlists.NpcTypes.Spawnentries.Spawngroup.Spawn2.Spawnentries",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Merchantlists.NpcTypes.Spawnentries.Spawngroup.Spawn2.Spawngroup",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.ObjectContents",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Objects",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Objects.Item",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Objects.Zone",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.StartingItems",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.StartingItems.Item",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.StartingItems.Zone",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Tasks",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Tasks.TaskActivities",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Tasks.TaskActivities.Goallists",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Tasks.TaskActivities.NpcType",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Tasks.TaskActivities.NpcType.AlternateCurrency",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Tasks.TaskActivities.NpcType.Loottable",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Tasks.TaskActivities.NpcType.Merchantlists",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Tasks.TaskActivities.NpcType.Merchantlists.Items",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Tasks.TaskActivities.NpcType.Merchantlists.NpcTypes",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Tasks.TaskActivities.NpcType.NpcEmotes",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Tasks.TaskActivities.NpcType.NpcFactions",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Tasks.TaskActivities.NpcType.NpcFactions.NpcFactionEntries",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Tasks.TaskActivities.NpcType.NpcFactions.NpcFactionEntries.FactionList",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Tasks.TaskActivities.NpcType.NpcSpell",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Tasks.TaskActivities.NpcType.NpcSpell.NpcSpellsEntries",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Tasks.TaskActivities.NpcType.NpcSpell.NpcSpellsEntries.SpellsNew",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Tasks.TaskActivities.NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.Aura",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Tasks.TaskActivities.NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.Aura.SpellsNew",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Tasks.TaskActivities.NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.BlockedSpells",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Tasks.TaskActivities.NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.Damageshieldtypes",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Tasks.TaskActivities.NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.Items",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Tasks.TaskActivities.NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.NpcSpellsEntries",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Tasks.TaskActivities.NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.SpellBuckets",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Tasks.TaskActivities.NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.SpellGlobals",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Tasks.TaskActivities.NpcType.NpcTypesTint",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Tasks.TaskActivities.NpcType.Spawnentries",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Tasks.TaskActivities.NpcType.Spawnentries.NpcType",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Tasks.TaskActivities.NpcType.Spawnentries.Spawngroup",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Tasks.TaskActivities.NpcType.Spawnentries.Spawngroup.Spawn2",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Tasks.TaskActivities.NpcType.Spawnentries.Spawngroup.Spawn2.Spawnentries",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Tasks.TaskActivities.NpcType.Spawnentries.Spawngroup.Spawn2.Spawngroup",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.Tasks.Tasksets",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.TradeskillRecipeEntries",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.TradeskillRecipeEntries.TradeskillRecipe",
		"LoottableEntries.Lootdrop.LootdropEntries.Item.TributeLevels",
		"LoottableEntries.Lootdrop.LootdropEntries.Lootdrop",
		"LoottableEntries.Lootdrop.LoottableEntries",
		"LoottableEntries.Loottable",
		"NpcTypes",
		"NpcTypes.AlternateCurrency",
		"NpcTypes.Loottable",
		"NpcTypes.Merchantlists",
		"NpcTypes.Merchantlists.Items",
		"NpcTypes.Merchantlists.Items.AlternateCurrencies",
		"NpcTypes.Merchantlists.Items.CharacterCorpseItems",
		"NpcTypes.Merchantlists.Items.DiscoveredItems",
		"NpcTypes.Merchantlists.Items.Doors",
		"NpcTypes.Merchantlists.Items.Doors.Item",
		"NpcTypes.Merchantlists.Items.Fishings",
		"NpcTypes.Merchantlists.Items.Fishings.Item",
		"NpcTypes.Merchantlists.Items.Fishings.NpcType",
		"NpcTypes.Merchantlists.Items.Fishings.Zone",
		"NpcTypes.Merchantlists.Items.Forages",
		"NpcTypes.Merchantlists.Items.Forages.Item",
		"NpcTypes.Merchantlists.Items.Forages.Zone",
		"NpcTypes.Merchantlists.Items.GroundSpawns",
		"NpcTypes.Merchantlists.Items.GroundSpawns.Zone",
		"NpcTypes.Merchantlists.Items.ItemTicks",
		"NpcTypes.Merchantlists.Items.Keyrings",
		"NpcTypes.Merchantlists.Items.LootdropEntries",
		"NpcTypes.Merchantlists.Items.LootdropEntries.Item",
		"NpcTypes.Merchantlists.Items.LootdropEntries.Lootdrop",
		"NpcTypes.Merchantlists.Items.LootdropEntries.Lootdrop.LootdropEntries",
		"NpcTypes.Merchantlists.Items.LootdropEntries.Lootdrop.LoottableEntries",
		"NpcTypes.Merchantlists.Items.LootdropEntries.Lootdrop.LoottableEntries.Lootdrop",
		"NpcTypes.Merchantlists.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable",
		"NpcTypes.Merchantlists.Items.Merchantlists",
		"NpcTypes.Merchantlists.Items.ObjectContents",
		"NpcTypes.Merchantlists.Items.Objects",
		"NpcTypes.Merchantlists.Items.Objects.Item",
		"NpcTypes.Merchantlists.Items.Objects.Zone",
		"NpcTypes.Merchantlists.Items.StartingItems",
		"NpcTypes.Merchantlists.Items.StartingItems.Item",
		"NpcTypes.Merchantlists.Items.StartingItems.Zone",
		"NpcTypes.Merchantlists.Items.Tasks",
		"NpcTypes.Merchantlists.Items.Tasks.TaskActivities",
		"NpcTypes.Merchantlists.Items.Tasks.TaskActivities.Goallists",
		"NpcTypes.Merchantlists.Items.Tasks.TaskActivities.NpcType",
		"NpcTypes.Merchantlists.Items.Tasks.Tasksets",
		"NpcTypes.Merchantlists.Items.TradeskillRecipeEntries",
		"NpcTypes.Merchantlists.Items.TradeskillRecipeEntries.TradeskillRecipe",
		"NpcTypes.Merchantlists.Items.TributeLevels",
		"NpcTypes.Merchantlists.NpcTypes",
		"NpcTypes.NpcEmotes",
		"NpcTypes.NpcFactions",
		"NpcTypes.NpcFactions.NpcFactionEntries",
		"NpcTypes.NpcFactions.NpcFactionEntries.FactionList",
		"NpcTypes.NpcSpell",
		"NpcTypes.NpcSpell.NpcSpellsEntries",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Aura",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Aura.SpellsNew",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.BlockedSpells",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Damageshieldtypes",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Items",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Items.AlternateCurrencies",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Items.CharacterCorpseItems",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Items.DiscoveredItems",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Items.Doors",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Items.Doors.Item",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings.Item",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings.NpcType",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings.Zone",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Items.Forages",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Items.Forages.Item",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Items.Forages.Zone",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Items.GroundSpawns",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Items.GroundSpawns.Zone",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Items.ItemTicks",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Items.Keyrings",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Item",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LootdropEntries",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Lootdrop",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Items.Merchantlists",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Items.Merchantlists.Items",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Items.Merchantlists.NpcTypes",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Items.ObjectContents",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Items.Objects",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Items.Objects.Item",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Items.Objects.Zone",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Items.StartingItems",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Items.StartingItems.Item",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Items.StartingItems.Zone",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Items.Tasks",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Items.Tasks.TaskActivities",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Items.Tasks.TaskActivities.Goallists",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Items.Tasks.TaskActivities.NpcType",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Items.Tasks.Tasksets",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Items.TradeskillRecipeEntries",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Items.TradeskillRecipeEntries.TradeskillRecipe",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Items.TributeLevels",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.NpcSpellsEntries",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.SpellBuckets",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.SpellGlobals",
		"NpcTypes.NpcTypesTint",
		"NpcTypes.Spawnentries",
		"NpcTypes.Spawnentries.NpcType",
		"NpcTypes.Spawnentries.Spawngroup",
		"NpcTypes.Spawnentries.Spawngroup.Spawn2",
		"NpcTypes.Spawnentries.Spawngroup.Spawn2.Spawnentries",
		"NpcTypes.Spawnentries.Spawngroup.Spawn2.Spawngroup",
	}
}

func (Loottable) Connection() string {
    return "eqemu_content"
}
