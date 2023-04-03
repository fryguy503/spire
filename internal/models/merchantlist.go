package models

import (
	"github.com/volatiletech/null/v8"
)

type Merchantlist struct {
	Merchantid             int         `json:"merchantid" gorm:"Column:merchantid"`
	Slot                   int         `json:"slot" gorm:"Column:slot"`
	Item                   int         `json:"item" gorm:"Column:item"`
	FactionRequired        int16       `json:"faction_required" gorm:"Column:faction_required"`
	LevelRequired          uint8       `json:"level_required" gorm:"Column:level_required"`
	MinStatus              uint8       `json:"min_status" gorm:"Column:min_status"`
	MaxStatus              uint8       `json:"max_status" gorm:"Column:max_status"`
	AltCurrencyCost        uint16      `json:"alt_currency_cost" gorm:"Column:alt_currency_cost"`
	ClassesRequired        int         `json:"classes_required" gorm:"Column:classes_required"`
	Probability            int         `json:"probability" gorm:"Column:probability"`
	BucketName             string      `json:"bucket_name" gorm:"Column:bucket_name"`
	BucketValue            string      `json:"bucket_value" gorm:"Column:bucket_value"`
	BucketComparison       null.Uint8  `json:"bucket_comparison" gorm:"Column:bucket_comparison"`
	MinExpansion           int8        `json:"min_expansion" gorm:"Column:min_expansion"`
	MaxExpansion           int8        `json:"max_expansion" gorm:"Column:max_expansion"`
	ContentFlags           null.String `json:"content_flags" gorm:"Column:content_flags"`
	ContentFlagsDisabled   null.String `json:"content_flags_disabled" gorm:"Column:content_flags_disabled"`
	NpcTypes               []NpcType   `json:"npc_types,omitempty" gorm:"foreignKey:merchant_id;references:merchantid"`
	Items                  []Item      `json:"items,omitempty" gorm:"foreignKey:id;references:item"`
}

func (Merchantlist) TableName() string {
    return "merchantlist"
}

func (Merchantlist) Relationships() []string {
    return []string{
		"Items",
		"Items.AlternateCurrencies",
		"Items.AlternateCurrencies.Item",
		"Items.CharacterCorpseItems",
		"Items.DiscoveredItems",
		"Items.Doors",
		"Items.Doors.Item",
		"Items.Fishings",
		"Items.Fishings.Item",
		"Items.Fishings.NpcType",
		"Items.Fishings.NpcType.AlternateCurrency",
		"Items.Fishings.NpcType.AlternateCurrency.Item",
		"Items.Fishings.NpcType.Loottable",
		"Items.Fishings.NpcType.Loottable.LoottableEntries",
		"Items.Fishings.NpcType.Loottable.LoottableEntries.Lootdrop",
		"Items.Fishings.NpcType.Loottable.LoottableEntries.Lootdrop.LootdropEntries",
		"Items.Fishings.NpcType.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item",
		"Items.Fishings.NpcType.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Lootdrop",
		"Items.Fishings.NpcType.Loottable.LoottableEntries.Lootdrop.LoottableEntries",
		"Items.Fishings.NpcType.Loottable.LoottableEntries.Loottable",
		"Items.Fishings.NpcType.Loottable.NpcTypes",
		"Items.Fishings.NpcType.Merchantlists",
		"Items.Fishings.NpcType.NpcEmotes",
		"Items.Fishings.NpcType.NpcFactions",
		"Items.Fishings.NpcType.NpcFactions.NpcFactionEntries",
		"Items.Fishings.NpcType.NpcFactions.NpcFactionEntries.FactionList",
		"Items.Fishings.NpcType.NpcSpell",
		"Items.Fishings.NpcType.NpcSpell.BotSpellsEntries",
		"Items.Fishings.NpcType.NpcSpell.BotSpellsEntries.NpcSpell",
		"Items.Fishings.NpcType.NpcSpell.BotSpellsEntries.SpellsNew",
		"Items.Fishings.NpcType.NpcSpell.BotSpellsEntries.SpellsNew.Aura",
		"Items.Fishings.NpcType.NpcSpell.BotSpellsEntries.SpellsNew.Aura.SpellsNew",
		"Items.Fishings.NpcType.NpcSpell.BotSpellsEntries.SpellsNew.BlockedSpells",
		"Items.Fishings.NpcType.NpcSpell.BotSpellsEntries.SpellsNew.BotSpellsEntries",
		"Items.Fishings.NpcType.NpcSpell.BotSpellsEntries.SpellsNew.Damageshieldtypes",
		"Items.Fishings.NpcType.NpcSpell.BotSpellsEntries.SpellsNew.Items",
		"Items.Fishings.NpcType.NpcSpell.BotSpellsEntries.SpellsNew.NpcSpellsEntries",
		"Items.Fishings.NpcType.NpcSpell.BotSpellsEntries.SpellsNew.NpcSpellsEntries.SpellsNew",
		"Items.Fishings.NpcType.NpcSpell.BotSpellsEntries.SpellsNew.SpellBuckets",
		"Items.Fishings.NpcType.NpcSpell.BotSpellsEntries.SpellsNew.SpellGlobals",
		"Items.Fishings.NpcType.NpcSpell.NpcSpell",
		"Items.Fishings.NpcType.NpcSpell.NpcSpellsEntries",
		"Items.Fishings.NpcType.NpcSpell.NpcSpellsEntries.SpellsNew",
		"Items.Fishings.NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.Aura",
		"Items.Fishings.NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.Aura.SpellsNew",
		"Items.Fishings.NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.BlockedSpells",
		"Items.Fishings.NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.BotSpellsEntries",
		"Items.Fishings.NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.BotSpellsEntries.NpcSpell",
		"Items.Fishings.NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.BotSpellsEntries.SpellsNew",
		"Items.Fishings.NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.Damageshieldtypes",
		"Items.Fishings.NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.Items",
		"Items.Fishings.NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.NpcSpellsEntries",
		"Items.Fishings.NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.SpellBuckets",
		"Items.Fishings.NpcType.NpcSpell.NpcSpellsEntries.SpellsNew.SpellGlobals",
		"Items.Fishings.NpcType.NpcTypesTint",
		"Items.Fishings.NpcType.Spawnentries",
		"Items.Fishings.NpcType.Spawnentries.NpcType",
		"Items.Fishings.NpcType.Spawnentries.Spawngroup",
		"Items.Fishings.NpcType.Spawnentries.Spawngroup.Spawn2",
		"Items.Fishings.NpcType.Spawnentries.Spawngroup.Spawn2.Spawnentries",
		"Items.Fishings.NpcType.Spawnentries.Spawngroup.Spawn2.Spawngroup",
		"Items.Fishings.Zone",
		"Items.Forages",
		"Items.Forages.Item",
		"Items.Forages.Zone",
		"Items.GroundSpawns",
		"Items.GroundSpawns.Zone",
		"Items.ItemTicks",
		"Items.Keyrings",
		"Items.LootdropEntries",
		"Items.LootdropEntries.Item",
		"Items.LootdropEntries.Lootdrop",
		"Items.LootdropEntries.Lootdrop.LootdropEntries",
		"Items.LootdropEntries.Lootdrop.LoottableEntries",
		"Items.LootdropEntries.Lootdrop.LoottableEntries.Lootdrop",
		"Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable",
		"Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.LoottableEntries",
		"Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes",
		"Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.AlternateCurrency",
		"Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.AlternateCurrency.Item",
		"Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Loottable",
		"Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Merchantlists",
		"Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcEmotes",
		"Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcFactions",
		"Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcFactions.NpcFactionEntries",
		"Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcFactions.NpcFactionEntries.FactionList",
		"Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcSpell",
		"Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcSpell.BotSpellsEntries",
		"Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcSpell.BotSpellsEntries.NpcSpell",
		"Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcSpell.BotSpellsEntries.SpellsNew",
		"Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcSpell.BotSpellsEntries.SpellsNew.Aura",
		"Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcSpell.BotSpellsEntries.SpellsNew.Aura.SpellsNew",
		"Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcSpell.BotSpellsEntries.SpellsNew.BlockedSpells",
		"Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcSpell.BotSpellsEntries.SpellsNew.BotSpellsEntries",
		"Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcSpell.BotSpellsEntries.SpellsNew.Damageshieldtypes",
		"Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcSpell.BotSpellsEntries.SpellsNew.Items",
		"Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcSpell.BotSpellsEntries.SpellsNew.NpcSpellsEntries",
		"Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcSpell.BotSpellsEntries.SpellsNew.NpcSpellsEntries.SpellsNew",
		"Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcSpell.BotSpellsEntries.SpellsNew.SpellBuckets",
		"Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcSpell.BotSpellsEntries.SpellsNew.SpellGlobals",
		"Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcSpell.NpcSpell",
		"Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcSpell.NpcSpellsEntries",
		"Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew",
		"Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Aura",
		"Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Aura.SpellsNew",
		"Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.BlockedSpells",
		"Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.BotSpellsEntries",
		"Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.BotSpellsEntries.NpcSpell",
		"Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.BotSpellsEntries.SpellsNew",
		"Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Damageshieldtypes",
		"Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Items",
		"Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.NpcSpellsEntries",
		"Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.SpellBuckets",
		"Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.SpellGlobals",
		"Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcTypesTint",
		"Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Spawnentries",
		"Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Spawnentries.NpcType",
		"Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Spawnentries.Spawngroup",
		"Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Spawnentries.Spawngroup.Spawn2",
		"Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Spawnentries.Spawngroup.Spawn2.Spawnentries",
		"Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Spawnentries.Spawngroup.Spawn2.Spawngroup",
		"Items.Merchantlists",
		"Items.ObjectContents",
		"Items.Objects",
		"Items.Objects.Item",
		"Items.Objects.Zone",
		"Items.StartingItems",
		"Items.StartingItems.Item",
		"Items.StartingItems.Zone",
		"Items.TradeskillRecipeEntries",
		"Items.TradeskillRecipeEntries.TradeskillRecipe",
		"Items.TributeLevels",
		"NpcTypes",
		"NpcTypes.AlternateCurrency",
		"NpcTypes.AlternateCurrency.Item",
		"NpcTypes.AlternateCurrency.Item.AlternateCurrencies",
		"NpcTypes.AlternateCurrency.Item.CharacterCorpseItems",
		"NpcTypes.AlternateCurrency.Item.DiscoveredItems",
		"NpcTypes.AlternateCurrency.Item.Doors",
		"NpcTypes.AlternateCurrency.Item.Doors.Item",
		"NpcTypes.AlternateCurrency.Item.Fishings",
		"NpcTypes.AlternateCurrency.Item.Fishings.Item",
		"NpcTypes.AlternateCurrency.Item.Fishings.NpcType",
		"NpcTypes.AlternateCurrency.Item.Fishings.Zone",
		"NpcTypes.AlternateCurrency.Item.Forages",
		"NpcTypes.AlternateCurrency.Item.Forages.Item",
		"NpcTypes.AlternateCurrency.Item.Forages.Zone",
		"NpcTypes.AlternateCurrency.Item.GroundSpawns",
		"NpcTypes.AlternateCurrency.Item.GroundSpawns.Zone",
		"NpcTypes.AlternateCurrency.Item.ItemTicks",
		"NpcTypes.AlternateCurrency.Item.Keyrings",
		"NpcTypes.AlternateCurrency.Item.LootdropEntries",
		"NpcTypes.AlternateCurrency.Item.LootdropEntries.Item",
		"NpcTypes.AlternateCurrency.Item.LootdropEntries.Lootdrop",
		"NpcTypes.AlternateCurrency.Item.LootdropEntries.Lootdrop.LootdropEntries",
		"NpcTypes.AlternateCurrency.Item.LootdropEntries.Lootdrop.LoottableEntries",
		"NpcTypes.AlternateCurrency.Item.LootdropEntries.Lootdrop.LoottableEntries.Lootdrop",
		"NpcTypes.AlternateCurrency.Item.LootdropEntries.Lootdrop.LoottableEntries.Loottable",
		"NpcTypes.AlternateCurrency.Item.LootdropEntries.Lootdrop.LoottableEntries.Loottable.LoottableEntries",
		"NpcTypes.AlternateCurrency.Item.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes",
		"NpcTypes.AlternateCurrency.Item.Merchantlists",
		"NpcTypes.AlternateCurrency.Item.ObjectContents",
		"NpcTypes.AlternateCurrency.Item.Objects",
		"NpcTypes.AlternateCurrency.Item.Objects.Item",
		"NpcTypes.AlternateCurrency.Item.Objects.Zone",
		"NpcTypes.AlternateCurrency.Item.StartingItems",
		"NpcTypes.AlternateCurrency.Item.StartingItems.Item",
		"NpcTypes.AlternateCurrency.Item.StartingItems.Zone",
		"NpcTypes.AlternateCurrency.Item.TradeskillRecipeEntries",
		"NpcTypes.AlternateCurrency.Item.TradeskillRecipeEntries.TradeskillRecipe",
		"NpcTypes.AlternateCurrency.Item.TributeLevels",
		"NpcTypes.Loottable",
		"NpcTypes.Loottable.LoottableEntries",
		"NpcTypes.Loottable.LoottableEntries.Lootdrop",
		"NpcTypes.Loottable.LoottableEntries.Lootdrop.LootdropEntries",
		"NpcTypes.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item",
		"NpcTypes.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item.AlternateCurrencies",
		"NpcTypes.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item.AlternateCurrencies.Item",
		"NpcTypes.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item.CharacterCorpseItems",
		"NpcTypes.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item.DiscoveredItems",
		"NpcTypes.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item.Doors",
		"NpcTypes.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item.Doors.Item",
		"NpcTypes.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item.Fishings",
		"NpcTypes.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item.Fishings.Item",
		"NpcTypes.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item.Fishings.NpcType",
		"NpcTypes.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item.Fishings.Zone",
		"NpcTypes.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item.Forages",
		"NpcTypes.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item.Forages.Item",
		"NpcTypes.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item.Forages.Zone",
		"NpcTypes.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item.GroundSpawns",
		"NpcTypes.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item.GroundSpawns.Zone",
		"NpcTypes.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item.ItemTicks",
		"NpcTypes.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item.Keyrings",
		"NpcTypes.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item.LootdropEntries",
		"NpcTypes.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item.Merchantlists",
		"NpcTypes.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item.ObjectContents",
		"NpcTypes.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item.Objects",
		"NpcTypes.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item.Objects.Item",
		"NpcTypes.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item.Objects.Zone",
		"NpcTypes.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item.StartingItems",
		"NpcTypes.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item.StartingItems.Item",
		"NpcTypes.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item.StartingItems.Zone",
		"NpcTypes.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item.TradeskillRecipeEntries",
		"NpcTypes.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item.TradeskillRecipeEntries.TradeskillRecipe",
		"NpcTypes.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item.TributeLevels",
		"NpcTypes.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Lootdrop",
		"NpcTypes.Loottable.LoottableEntries.Lootdrop.LoottableEntries",
		"NpcTypes.Loottable.LoottableEntries.Loottable",
		"NpcTypes.Loottable.NpcTypes",
		"NpcTypes.Merchantlists",
		"NpcTypes.NpcEmotes",
		"NpcTypes.NpcFactions",
		"NpcTypes.NpcFactions.NpcFactionEntries",
		"NpcTypes.NpcFactions.NpcFactionEntries.FactionList",
		"NpcTypes.NpcSpell",
		"NpcTypes.NpcSpell.BotSpellsEntries",
		"NpcTypes.NpcSpell.BotSpellsEntries.NpcSpell",
		"NpcTypes.NpcSpell.BotSpellsEntries.SpellsNew",
		"NpcTypes.NpcSpell.BotSpellsEntries.SpellsNew.Aura",
		"NpcTypes.NpcSpell.BotSpellsEntries.SpellsNew.Aura.SpellsNew",
		"NpcTypes.NpcSpell.BotSpellsEntries.SpellsNew.BlockedSpells",
		"NpcTypes.NpcSpell.BotSpellsEntries.SpellsNew.BotSpellsEntries",
		"NpcTypes.NpcSpell.BotSpellsEntries.SpellsNew.Damageshieldtypes",
		"NpcTypes.NpcSpell.BotSpellsEntries.SpellsNew.Items",
		"NpcTypes.NpcSpell.BotSpellsEntries.SpellsNew.Items.AlternateCurrencies",
		"NpcTypes.NpcSpell.BotSpellsEntries.SpellsNew.Items.AlternateCurrencies.Item",
		"NpcTypes.NpcSpell.BotSpellsEntries.SpellsNew.Items.CharacterCorpseItems",
		"NpcTypes.NpcSpell.BotSpellsEntries.SpellsNew.Items.DiscoveredItems",
		"NpcTypes.NpcSpell.BotSpellsEntries.SpellsNew.Items.Doors",
		"NpcTypes.NpcSpell.BotSpellsEntries.SpellsNew.Items.Doors.Item",
		"NpcTypes.NpcSpell.BotSpellsEntries.SpellsNew.Items.Fishings",
		"NpcTypes.NpcSpell.BotSpellsEntries.SpellsNew.Items.Fishings.Item",
		"NpcTypes.NpcSpell.BotSpellsEntries.SpellsNew.Items.Fishings.NpcType",
		"NpcTypes.NpcSpell.BotSpellsEntries.SpellsNew.Items.Fishings.Zone",
		"NpcTypes.NpcSpell.BotSpellsEntries.SpellsNew.Items.Forages",
		"NpcTypes.NpcSpell.BotSpellsEntries.SpellsNew.Items.Forages.Item",
		"NpcTypes.NpcSpell.BotSpellsEntries.SpellsNew.Items.Forages.Zone",
		"NpcTypes.NpcSpell.BotSpellsEntries.SpellsNew.Items.GroundSpawns",
		"NpcTypes.NpcSpell.BotSpellsEntries.SpellsNew.Items.GroundSpawns.Zone",
		"NpcTypes.NpcSpell.BotSpellsEntries.SpellsNew.Items.ItemTicks",
		"NpcTypes.NpcSpell.BotSpellsEntries.SpellsNew.Items.Keyrings",
		"NpcTypes.NpcSpell.BotSpellsEntries.SpellsNew.Items.LootdropEntries",
		"NpcTypes.NpcSpell.BotSpellsEntries.SpellsNew.Items.LootdropEntries.Item",
		"NpcTypes.NpcSpell.BotSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop",
		"NpcTypes.NpcSpell.BotSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LootdropEntries",
		"NpcTypes.NpcSpell.BotSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries",
		"NpcTypes.NpcSpell.BotSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Lootdrop",
		"NpcTypes.NpcSpell.BotSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable",
		"NpcTypes.NpcSpell.BotSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.LoottableEntries",
		"NpcTypes.NpcSpell.BotSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes",
		"NpcTypes.NpcSpell.BotSpellsEntries.SpellsNew.Items.Merchantlists",
		"NpcTypes.NpcSpell.BotSpellsEntries.SpellsNew.Items.ObjectContents",
		"NpcTypes.NpcSpell.BotSpellsEntries.SpellsNew.Items.Objects",
		"NpcTypes.NpcSpell.BotSpellsEntries.SpellsNew.Items.Objects.Item",
		"NpcTypes.NpcSpell.BotSpellsEntries.SpellsNew.Items.Objects.Zone",
		"NpcTypes.NpcSpell.BotSpellsEntries.SpellsNew.Items.StartingItems",
		"NpcTypes.NpcSpell.BotSpellsEntries.SpellsNew.Items.StartingItems.Item",
		"NpcTypes.NpcSpell.BotSpellsEntries.SpellsNew.Items.StartingItems.Zone",
		"NpcTypes.NpcSpell.BotSpellsEntries.SpellsNew.Items.TradeskillRecipeEntries",
		"NpcTypes.NpcSpell.BotSpellsEntries.SpellsNew.Items.TradeskillRecipeEntries.TradeskillRecipe",
		"NpcTypes.NpcSpell.BotSpellsEntries.SpellsNew.Items.TributeLevels",
		"NpcTypes.NpcSpell.BotSpellsEntries.SpellsNew.NpcSpellsEntries",
		"NpcTypes.NpcSpell.BotSpellsEntries.SpellsNew.NpcSpellsEntries.SpellsNew",
		"NpcTypes.NpcSpell.BotSpellsEntries.SpellsNew.SpellBuckets",
		"NpcTypes.NpcSpell.BotSpellsEntries.SpellsNew.SpellGlobals",
		"NpcTypes.NpcSpell.NpcSpell",
		"NpcTypes.NpcSpell.NpcSpellsEntries",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Aura",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Aura.SpellsNew",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.BlockedSpells",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.BotSpellsEntries",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.BotSpellsEntries.NpcSpell",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.BotSpellsEntries.SpellsNew",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Damageshieldtypes",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Items",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Items.AlternateCurrencies",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Items.AlternateCurrencies.Item",
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
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.LoottableEntries",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Items.Merchantlists",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Items.ObjectContents",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Items.Objects",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Items.Objects.Item",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Items.Objects.Zone",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Items.StartingItems",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Items.StartingItems.Item",
		"NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew.Items.StartingItems.Zone",
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

func (Merchantlist) Connection() string {
    return "eqemu_content"
}
