package models

import (
	"github.com/volatiletech/null/v8"
)

type Fishing struct {
	ID                     int         `json:"id" gorm:"Column:id"`
	Zoneid                 int         `json:"zoneid" gorm:"Column:zoneid"`
	Itemid                 int         `json:"itemid" gorm:"Column:Itemid"`
	SkillLevel             int16       `json:"skill_level" gorm:"Column:skill_level"`
	Chance                 int16       `json:"chance" gorm:"Column:chance"`
	NpcId                  int         `json:"npc_id" gorm:"Column:npc_id"`
	NpcChance              int         `json:"npc_chance" gorm:"Column:npc_chance"`
	MinExpansion           int8        `json:"min_expansion" gorm:"Column:min_expansion"`
	MaxExpansion           int8        `json:"max_expansion" gorm:"Column:max_expansion"`
	ContentFlags           null.String `json:"content_flags" gorm:"Column:content_flags"`
	ContentFlagsDisabled   null.String `json:"content_flags_disabled" gorm:"Column:content_flags_disabled"`
	Item                   *Item       `json:"item,omitempty" gorm:"foreignKey:Itemid;references:id"`
	Zone                   *Zone       `json:"zone,omitempty" gorm:"foreignKey:zoneid;references:zoneidnumber"`
	NpcType                *NpcType    `json:"npc_type,omitempty" gorm:"foreignKey:npc_id;references:id"`
}

func (Fishing) TableName() string {
    return "fishing"
}

func (Fishing) Relationships() []string {
    return []string{
		"Item",
		"Item.AlternateCurrencies",
		"Item.CharacterCorpseItems",
		"Item.DiscoveredItems",
		"Item.Doors",
		"Item.Doors.Item",
		"Item.Fishings",
		"Item.Forages",
		"Item.Forages.Item",
		"Item.Forages.Zone",
		"Item.GroundSpawns",
		"Item.GroundSpawns.Zone",
		"Item.ItemTicks",
		"Item.Keyrings",
		"Item.LootdropEntries",
		"Item.LootdropEntries.Item",
		"Item.LootdropEntries.Lootdrop",
		"Item.LootdropEntries.Lootdrop.LootdropEntries",
		"Item.LootdropEntries.Lootdrop.LoottableEntries",
		"Item.LootdropEntries.Lootdrop.LoottableEntries.Lootdrop",
		"Item.LootdropEntries.Lootdrop.LoottableEntries.Loottable",
		"Item.LootdropEntries.Lootdrop.LoottableEntries.Loottable.LoottableEntries",
		"Item.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes",
		"Item.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.AlternateCurrency",
		"Item.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Loottable",
		"Item.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Merchantlists",
		"Item.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Merchantlists.Items",
		"Item.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Merchantlists.NpcTypes",
		"Item.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcEmotes",
		"Item.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcFactions",
		"Item.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcFactions.NpcFactionEntries",
		"Item.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcFactions.NpcFactionEntries.FactionList",
		"Item.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcSpell",
		"Item.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcSpell.NpcSpellsEntries",
		"Item.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew",
		"Item.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Aura",
		"Item.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Aura.SpellsNew",
		"Item.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.BlockedSpells",
		"Item.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Damageshieldtypes",
		"Item.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Items",
		"Item.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.NpcSpellsEntries",
		"Item.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.SpellBuckets",
		"Item.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.SpellGlobals",
		"Item.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcTypesTint",
		"Item.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Spawnentries",
		"Item.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Spawnentries.NpcType",
		"Item.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Spawnentries.Spawngroup",
		"Item.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Spawnentries.Spawngroup.Spawn2",
		"Item.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Spawnentries.Spawngroup.Spawn2.Spawnentries",
		"Item.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Spawnentries.Spawngroup.Spawn2.Spawngroup",
		"Item.Merchantlists",
		"Item.Merchantlists.Items",
		"Item.Merchantlists.NpcTypes",
		"Item.Merchantlists.NpcTypes.AlternateCurrency",
		"Item.Merchantlists.NpcTypes.Loottable",
		"Item.Merchantlists.NpcTypes.Loottable.LoottableEntries",
		"Item.Merchantlists.NpcTypes.Loottable.LoottableEntries.Lootdrop",
		"Item.Merchantlists.NpcTypes.Loottable.LoottableEntries.Lootdrop.LootdropEntries",
		"Item.Merchantlists.NpcTypes.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item",
		"Item.Merchantlists.NpcTypes.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Lootdrop",
		"Item.Merchantlists.NpcTypes.Loottable.LoottableEntries.Lootdrop.LoottableEntries",
		"Item.Merchantlists.NpcTypes.Loottable.LoottableEntries.Loottable",
		"Item.Merchantlists.NpcTypes.Loottable.NpcTypes",
		"Item.Merchantlists.NpcTypes.Merchantlists",
		"Item.Merchantlists.NpcTypes.NpcEmotes",
		"Item.Merchantlists.NpcTypes.NpcFactions",
		"Item.Merchantlists.NpcTypes.NpcFactions.NpcFactionEntries",
		"Item.Merchantlists.NpcTypes.NpcFactions.NpcFactionEntries.FactionList",
		"Item.Merchantlists.NpcTypes.NpcSpell",
		"Item.Merchantlists.NpcTypes.NpcSpell.NpcSpellsEntries",
		"Item.Merchantlists.NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew",
		"Item.Merchantlists.NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Aura",
		"Item.Merchantlists.NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Aura.SpellsNew",
		"Item.Merchantlists.NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.BlockedSpells",
		"Item.Merchantlists.NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Damageshieldtypes",
		"Item.Merchantlists.NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Items",
		"Item.Merchantlists.NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.NpcSpellsEntries",
		"Item.Merchantlists.NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.SpellBuckets",
		"Item.Merchantlists.NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.SpellGlobals",
		"Item.Merchantlists.NpcTypes.NpcTypesTint",
		"Item.Merchantlists.NpcTypes.Spawnentries",
		"Item.Merchantlists.NpcTypes.Spawnentries.NpcType",
		"Item.Merchantlists.NpcTypes.Spawnentries.Spawngroup",
		"Item.Merchantlists.NpcTypes.Spawnentries.Spawngroup.Spawn2",
		"Item.Merchantlists.NpcTypes.Spawnentries.Spawngroup.Spawn2.Spawnentries",
		"Item.Merchantlists.NpcTypes.Spawnentries.Spawngroup.Spawn2.Spawngroup",
		"Item.ObjectContents",
		"Item.Objects",
		"Item.Objects.Item",
		"Item.Objects.Zone",
		"Item.StartingItems",
		"Item.StartingItems.Item",
		"Item.StartingItems.Zone",
		"Item.Tasks",
		"Item.Tasks.TaskActivities",
		"Item.Tasks.TaskActivities.Goallists",
		"Item.Tasks.TaskActivities.NpcType",
		"Item.Tasks.TaskActivities.NpcType.AlternateCurrency",
		"Item.Tasks.TaskActivities.NpcType.Loottable",
		"Item.Tasks.TaskActivities.NpcType.Loottable.LoottableEntries",
		"Item.Tasks.TaskActivities.NpcType.Loottable.LoottableEntries.Lootdrop",
		"Item.Tasks.TaskActivities.NpcType.Loottable.LoottableEntries.Lootdrop.LootdropEntries",
		"Item.Tasks.TaskActivities.NpcType.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item",
		"Item.Tasks.TaskActivities.NpcType.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Lootdrop",
		"Item.Tasks.TaskActivities.NpcType.Loottable.LoottableEntries.Lootdrop.LoottableEntries",
		"Item.Tasks.TaskActivities.NpcType.Loottable.LoottableEntries.Loottable",
		"Item.Tasks.TaskActivities.NpcType.Loottable.NpcTypes",
		"Item.Tasks.TaskActivities.NpcType.Merchantlists",
		"Item.Tasks.TaskActivities.NpcType.Merchantlists.Items",
		"Item.Tasks.TaskActivities.NpcType.Merchantlists.NpcTypes",
		"Item.Tasks.TaskActivities.NpcType.NpcEmotes",
		"Item.Tasks.TaskActivities.NpcType.NpcFactions",
		"Item.Tasks.TaskActivities.NpcType.NpcFactions.NpcFactionEntries",
		"Item.Tasks.TaskActivities.NpcType.NpcFactions.NpcFactionEntries.FactionList",
		"Item.Tasks.TaskActivities.NpcType.NpcSpell",
		"Item.Tasks.TaskActivities.NpcType.NpcSpell.NpcSpellsEntries",
		"Item.Tasks.TaskActivities.NpcType.NpcSpell.NpcSpellsEntries.SpellsNew",
		"Item.Tasks.TaskActivities.NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.Aura",
		"Item.Tasks.TaskActivities.NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.Aura.SpellsNew",
		"Item.Tasks.TaskActivities.NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.BlockedSpells",
		"Item.Tasks.TaskActivities.NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.Damageshieldtypes",
		"Item.Tasks.TaskActivities.NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.Items",
		"Item.Tasks.TaskActivities.NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.NpcSpellsEntries",
		"Item.Tasks.TaskActivities.NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.SpellBuckets",
		"Item.Tasks.TaskActivities.NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.SpellGlobals",
		"Item.Tasks.TaskActivities.NpcType.NpcTypesTint",
		"Item.Tasks.TaskActivities.NpcType.Spawnentries",
		"Item.Tasks.TaskActivities.NpcType.Spawnentries.NpcType",
		"Item.Tasks.TaskActivities.NpcType.Spawnentries.Spawngroup",
		"Item.Tasks.TaskActivities.NpcType.Spawnentries.Spawngroup.Spawn2",
		"Item.Tasks.TaskActivities.NpcType.Spawnentries.Spawngroup.Spawn2.Spawnentries",
		"Item.Tasks.TaskActivities.NpcType.Spawnentries.Spawngroup.Spawn2.Spawngroup",
		"Item.Tasks.Tasksets",
		"Item.TradeskillRecipeEntries",
		"Item.TradeskillRecipeEntries.TradeskillRecipe",
		"Item.TributeLevels",
		"NpcType",
		"NpcType.AlternateCurrency",
		"NpcType.Loottable",
		"NpcType.Loottable.LoottableEntries",
		"NpcType.Loottable.LoottableEntries.Lootdrop",
		"NpcType.Loottable.LoottableEntries.Lootdrop.LootdropEntries",
		"NpcType.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item",
		"NpcType.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item.AlternateCurrencies",
		"NpcType.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item.CharacterCorpseItems",
		"NpcType.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item.DiscoveredItems",
		"NpcType.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item.Doors",
		"NpcType.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item.Doors.Item",
		"NpcType.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item.Fishings",
		"NpcType.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item.Forages",
		"NpcType.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item.Forages.Item",
		"NpcType.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item.Forages.Zone",
		"NpcType.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item.GroundSpawns",
		"NpcType.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item.GroundSpawns.Zone",
		"NpcType.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item.ItemTicks",
		"NpcType.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item.Keyrings",
		"NpcType.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item.LootdropEntries",
		"NpcType.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item.Merchantlists",
		"NpcType.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item.Merchantlists.Items",
		"NpcType.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item.Merchantlists.NpcTypes",
		"NpcType.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item.ObjectContents",
		"NpcType.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item.Objects",
		"NpcType.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item.Objects.Item",
		"NpcType.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item.Objects.Zone",
		"NpcType.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item.StartingItems",
		"NpcType.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item.StartingItems.Item",
		"NpcType.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item.StartingItems.Zone",
		"NpcType.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item.Tasks",
		"NpcType.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item.Tasks.TaskActivities",
		"NpcType.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item.Tasks.TaskActivities.Goallists",
		"NpcType.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item.Tasks.TaskActivities.NpcType",
		"NpcType.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item.Tasks.Tasksets",
		"NpcType.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item.TradeskillRecipeEntries",
		"NpcType.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item.TradeskillRecipeEntries.TradeskillRecipe",
		"NpcType.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item.TributeLevels",
		"NpcType.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Lootdrop",
		"NpcType.Loottable.LoottableEntries.Lootdrop.LoottableEntries",
		"NpcType.Loottable.LoottableEntries.Loottable",
		"NpcType.Loottable.NpcTypes",
		"NpcType.Merchantlists",
		"NpcType.Merchantlists.Items",
		"NpcType.Merchantlists.Items.AlternateCurrencies",
		"NpcType.Merchantlists.Items.CharacterCorpseItems",
		"NpcType.Merchantlists.Items.DiscoveredItems",
		"NpcType.Merchantlists.Items.Doors",
		"NpcType.Merchantlists.Items.Doors.Item",
		"NpcType.Merchantlists.Items.Fishings",
		"NpcType.Merchantlists.Items.Forages",
		"NpcType.Merchantlists.Items.Forages.Item",
		"NpcType.Merchantlists.Items.Forages.Zone",
		"NpcType.Merchantlists.Items.GroundSpawns",
		"NpcType.Merchantlists.Items.GroundSpawns.Zone",
		"NpcType.Merchantlists.Items.ItemTicks",
		"NpcType.Merchantlists.Items.Keyrings",
		"NpcType.Merchantlists.Items.LootdropEntries",
		"NpcType.Merchantlists.Items.LootdropEntries.Item",
		"NpcType.Merchantlists.Items.LootdropEntries.Lootdrop",
		"NpcType.Merchantlists.Items.LootdropEntries.Lootdrop.LootdropEntries",
		"NpcType.Merchantlists.Items.LootdropEntries.Lootdrop.LoottableEntries",
		"NpcType.Merchantlists.Items.LootdropEntries.Lootdrop.LoottableEntries.Lootdrop",
		"NpcType.Merchantlists.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable",
		"NpcType.Merchantlists.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.LoottableEntries",
		"NpcType.Merchantlists.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes",
		"NpcType.Merchantlists.Items.Merchantlists",
		"NpcType.Merchantlists.Items.ObjectContents",
		"NpcType.Merchantlists.Items.Objects",
		"NpcType.Merchantlists.Items.Objects.Item",
		"NpcType.Merchantlists.Items.Objects.Zone",
		"NpcType.Merchantlists.Items.StartingItems",
		"NpcType.Merchantlists.Items.StartingItems.Item",
		"NpcType.Merchantlists.Items.StartingItems.Zone",
		"NpcType.Merchantlists.Items.Tasks",
		"NpcType.Merchantlists.Items.Tasks.TaskActivities",
		"NpcType.Merchantlists.Items.Tasks.TaskActivities.Goallists",
		"NpcType.Merchantlists.Items.Tasks.TaskActivities.NpcType",
		"NpcType.Merchantlists.Items.Tasks.Tasksets",
		"NpcType.Merchantlists.Items.TradeskillRecipeEntries",
		"NpcType.Merchantlists.Items.TradeskillRecipeEntries.TradeskillRecipe",
		"NpcType.Merchantlists.Items.TributeLevels",
		"NpcType.Merchantlists.NpcTypes",
		"NpcType.NpcEmotes",
		"NpcType.NpcFactions",
		"NpcType.NpcFactions.NpcFactionEntries",
		"NpcType.NpcFactions.NpcFactionEntries.FactionList",
		"NpcType.NpcSpell",
		"NpcType.NpcSpell.NpcSpellsEntries",
		"NpcType.NpcSpell.NpcSpellsEntries.SpellsNew",
		"NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.Aura",
		"NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.Aura.SpellsNew",
		"NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.BlockedSpells",
		"NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.Damageshieldtypes",
		"NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.Items",
		"NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.Items.AlternateCurrencies",
		"NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.Items.CharacterCorpseItems",
		"NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.Items.DiscoveredItems",
		"NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.Items.Doors",
		"NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.Items.Doors.Item",
		"NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings",
		"NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.Items.Forages",
		"NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.Items.Forages.Item",
		"NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.Items.Forages.Zone",
		"NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.Items.GroundSpawns",
		"NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.Items.GroundSpawns.Zone",
		"NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.Items.ItemTicks",
		"NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.Items.Keyrings",
		"NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries",
		"NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Item",
		"NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop",
		"NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LootdropEntries",
		"NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries",
		"NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Lootdrop",
		"NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable",
		"NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.LoottableEntries",
		"NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes",
		"NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.Items.Merchantlists",
		"NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.Items.Merchantlists.Items",
		"NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.Items.Merchantlists.NpcTypes",
		"NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.Items.ObjectContents",
		"NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.Items.Objects",
		"NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.Items.Objects.Item",
		"NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.Items.Objects.Zone",
		"NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.Items.StartingItems",
		"NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.Items.StartingItems.Item",
		"NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.Items.StartingItems.Zone",
		"NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.Items.Tasks",
		"NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.Items.Tasks.TaskActivities",
		"NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.Items.Tasks.TaskActivities.Goallists",
		"NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.Items.Tasks.TaskActivities.NpcType",
		"NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.Items.Tasks.Tasksets",
		"NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.Items.TradeskillRecipeEntries",
		"NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.Items.TradeskillRecipeEntries.TradeskillRecipe",
		"NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.Items.TributeLevels",
		"NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.NpcSpellsEntries",
		"NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.SpellBuckets",
		"NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.SpellGlobals",
		"NpcType.NpcTypesTint",
		"NpcType.Spawnentries",
		"NpcType.Spawnentries.NpcType",
		"NpcType.Spawnentries.Spawngroup",
		"NpcType.Spawnentries.Spawngroup.Spawn2",
		"NpcType.Spawnentries.Spawngroup.Spawn2.Spawnentries",
		"NpcType.Spawnentries.Spawngroup.Spawn2.Spawngroup",
		"Zone",
	}
}

func (Fishing) Connection() string {
    return "eqemu_content"
}
