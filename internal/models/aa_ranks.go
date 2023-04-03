package models

type AaRank struct {
	ID               uint       `json:"id" gorm:"Column:id"`
	UpperHotkeySid   int        `json:"upper_hotkey_sid" gorm:"Column:upper_hotkey_sid"`
	LowerHotkeySid   int        `json:"lower_hotkey_sid" gorm:"Column:lower_hotkey_sid"`
	TitleSid         int        `json:"title_sid" gorm:"Column:title_sid"`
	DescSid          int        `json:"desc_sid" gorm:"Column:desc_sid"`
	Cost             int        `json:"cost" gorm:"Column:cost"`
	LevelReq         int        `json:"level_req" gorm:"Column:level_req"`
	Spell            int        `json:"spell" gorm:"Column:spell"`
	SpellType        int        `json:"spell_type" gorm:"Column:spell_type"`
	RecastTime       int        `json:"recast_time" gorm:"Column:recast_time"`
	Expansion        int        `json:"expansion" gorm:"Column:expansion"`
	PrevId           int        `json:"prev_id" gorm:"Column:prev_id"`
	NextId           int        `json:"next_id" gorm:"Column:next_id"`
	SpellsNew        *SpellsNew `json:"spells_new,omitempty" gorm:"foreignKey:spell;references:id"`
	AaAbility        *AaAbility `json:"aa_ability,omitempty" gorm:"foreignKey:id;references:first_rank_id"`
}

func (AaRank) TableName() string {
    return "aa_ranks"
}

func (AaRank) Relationships() []string {
    return []string{
		"AaAbility",
		"SpellsNew",
		"SpellsNew.Aura",
		"SpellsNew.Aura.SpellsNew",
		"SpellsNew.BlockedSpells",
		"SpellsNew.BotSpellsEntries",
		"SpellsNew.BotSpellsEntries.NpcSpell",
		"SpellsNew.BotSpellsEntries.NpcSpell.BotSpellsEntries",
		"SpellsNew.BotSpellsEntries.NpcSpell.NpcSpell",
		"SpellsNew.BotSpellsEntries.NpcSpell.NpcSpellsEntries",
		"SpellsNew.BotSpellsEntries.NpcSpell.NpcSpellsEntries.SpellsNew",
		"SpellsNew.BotSpellsEntries.SpellsNew",
		"SpellsNew.Damageshieldtypes",
		"SpellsNew.Items",
		"SpellsNew.Items.AlternateCurrencies",
		"SpellsNew.Items.AlternateCurrencies.Item",
		"SpellsNew.Items.CharacterCorpseItems",
		"SpellsNew.Items.DiscoveredItems",
		"SpellsNew.Items.Doors",
		"SpellsNew.Items.Doors.Item",
		"SpellsNew.Items.Fishings",
		"SpellsNew.Items.Fishings.Item",
		"SpellsNew.Items.Fishings.NpcType",
		"SpellsNew.Items.Fishings.NpcType.AlternateCurrency",
		"SpellsNew.Items.Fishings.NpcType.AlternateCurrency.Item",
		"SpellsNew.Items.Fishings.NpcType.Loottable",
		"SpellsNew.Items.Fishings.NpcType.Loottable.LoottableEntries",
		"SpellsNew.Items.Fishings.NpcType.Loottable.LoottableEntries.Lootdrop",
		"SpellsNew.Items.Fishings.NpcType.Loottable.LoottableEntries.Lootdrop.LootdropEntries",
		"SpellsNew.Items.Fishings.NpcType.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item",
		"SpellsNew.Items.Fishings.NpcType.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Lootdrop",
		"SpellsNew.Items.Fishings.NpcType.Loottable.LoottableEntries.Lootdrop.LoottableEntries",
		"SpellsNew.Items.Fishings.NpcType.Loottable.LoottableEntries.Loottable",
		"SpellsNew.Items.Fishings.NpcType.Loottable.NpcTypes",
		"SpellsNew.Items.Fishings.NpcType.Merchantlists",
		"SpellsNew.Items.Fishings.NpcType.Merchantlists.Items",
		"SpellsNew.Items.Fishings.NpcType.Merchantlists.NpcTypes",
		"SpellsNew.Items.Fishings.NpcType.NpcEmotes",
		"SpellsNew.Items.Fishings.NpcType.NpcFactions",
		"SpellsNew.Items.Fishings.NpcType.NpcFactions.NpcFactionEntries",
		"SpellsNew.Items.Fishings.NpcType.NpcFactions.NpcFactionEntries.FactionList",
		"SpellsNew.Items.Fishings.NpcType.NpcSpell",
		"SpellsNew.Items.Fishings.NpcType.NpcSpell.BotSpellsEntries",
		"SpellsNew.Items.Fishings.NpcType.NpcSpell.BotSpellsEntries.NpcSpell",
		"SpellsNew.Items.Fishings.NpcType.NpcSpell.BotSpellsEntries.SpellsNew",
		"SpellsNew.Items.Fishings.NpcType.NpcSpell.NpcSpell",
		"SpellsNew.Items.Fishings.NpcType.NpcSpell.NpcSpellsEntries",
		"SpellsNew.Items.Fishings.NpcType.NpcSpell.NpcSpellsEntries.SpellsNew",
		"SpellsNew.Items.Fishings.NpcType.NpcTypesTint",
		"SpellsNew.Items.Fishings.NpcType.Spawnentries",
		"SpellsNew.Items.Fishings.NpcType.Spawnentries.NpcType",
		"SpellsNew.Items.Fishings.NpcType.Spawnentries.Spawngroup",
		"SpellsNew.Items.Fishings.NpcType.Spawnentries.Spawngroup.Spawn2",
		"SpellsNew.Items.Fishings.NpcType.Spawnentries.Spawngroup.Spawn2.Spawnentries",
		"SpellsNew.Items.Fishings.NpcType.Spawnentries.Spawngroup.Spawn2.Spawngroup",
		"SpellsNew.Items.Fishings.Zone",
		"SpellsNew.Items.Forages",
		"SpellsNew.Items.Forages.Item",
		"SpellsNew.Items.Forages.Zone",
		"SpellsNew.Items.GroundSpawns",
		"SpellsNew.Items.GroundSpawns.Zone",
		"SpellsNew.Items.ItemTicks",
		"SpellsNew.Items.Keyrings",
		"SpellsNew.Items.LootdropEntries",
		"SpellsNew.Items.LootdropEntries.Item",
		"SpellsNew.Items.LootdropEntries.Lootdrop",
		"SpellsNew.Items.LootdropEntries.Lootdrop.LootdropEntries",
		"SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries",
		"SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Lootdrop",
		"SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable",
		"SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.LoottableEntries",
		"SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes",
		"SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.AlternateCurrency",
		"SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.AlternateCurrency.Item",
		"SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Loottable",
		"SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Merchantlists",
		"SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Merchantlists.Items",
		"SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Merchantlists.NpcTypes",
		"SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcEmotes",
		"SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcFactions",
		"SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcFactions.NpcFactionEntries",
		"SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcFactions.NpcFactionEntries.FactionList",
		"SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcSpell",
		"SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcSpell.BotSpellsEntries",
		"SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcSpell.BotSpellsEntries.NpcSpell",
		"SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcSpell.BotSpellsEntries.SpellsNew",
		"SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcSpell.NpcSpell",
		"SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcSpell.NpcSpellsEntries",
		"SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew",
		"SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcTypesTint",
		"SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Spawnentries",
		"SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Spawnentries.NpcType",
		"SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Spawnentries.Spawngroup",
		"SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Spawnentries.Spawngroup.Spawn2",
		"SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Spawnentries.Spawngroup.Spawn2.Spawnentries",
		"SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Spawnentries.Spawngroup.Spawn2.Spawngroup",
		"SpellsNew.Items.Merchantlists",
		"SpellsNew.Items.Merchantlists.Items",
		"SpellsNew.Items.Merchantlists.NpcTypes",
		"SpellsNew.Items.Merchantlists.NpcTypes.AlternateCurrency",
		"SpellsNew.Items.Merchantlists.NpcTypes.AlternateCurrency.Item",
		"SpellsNew.Items.Merchantlists.NpcTypes.Loottable",
		"SpellsNew.Items.Merchantlists.NpcTypes.Loottable.LoottableEntries",
		"SpellsNew.Items.Merchantlists.NpcTypes.Loottable.LoottableEntries.Lootdrop",
		"SpellsNew.Items.Merchantlists.NpcTypes.Loottable.LoottableEntries.Lootdrop.LootdropEntries",
		"SpellsNew.Items.Merchantlists.NpcTypes.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item",
		"SpellsNew.Items.Merchantlists.NpcTypes.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Lootdrop",
		"SpellsNew.Items.Merchantlists.NpcTypes.Loottable.LoottableEntries.Lootdrop.LoottableEntries",
		"SpellsNew.Items.Merchantlists.NpcTypes.Loottable.LoottableEntries.Loottable",
		"SpellsNew.Items.Merchantlists.NpcTypes.Loottable.NpcTypes",
		"SpellsNew.Items.Merchantlists.NpcTypes.Merchantlists",
		"SpellsNew.Items.Merchantlists.NpcTypes.NpcEmotes",
		"SpellsNew.Items.Merchantlists.NpcTypes.NpcFactions",
		"SpellsNew.Items.Merchantlists.NpcTypes.NpcFactions.NpcFactionEntries",
		"SpellsNew.Items.Merchantlists.NpcTypes.NpcFactions.NpcFactionEntries.FactionList",
		"SpellsNew.Items.Merchantlists.NpcTypes.NpcSpell",
		"SpellsNew.Items.Merchantlists.NpcTypes.NpcSpell.BotSpellsEntries",
		"SpellsNew.Items.Merchantlists.NpcTypes.NpcSpell.BotSpellsEntries.NpcSpell",
		"SpellsNew.Items.Merchantlists.NpcTypes.NpcSpell.BotSpellsEntries.SpellsNew",
		"SpellsNew.Items.Merchantlists.NpcTypes.NpcSpell.NpcSpell",
		"SpellsNew.Items.Merchantlists.NpcTypes.NpcSpell.NpcSpellsEntries",
		"SpellsNew.Items.Merchantlists.NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew",
		"SpellsNew.Items.Merchantlists.NpcTypes.NpcTypesTint",
		"SpellsNew.Items.Merchantlists.NpcTypes.Spawnentries",
		"SpellsNew.Items.Merchantlists.NpcTypes.Spawnentries.NpcType",
		"SpellsNew.Items.Merchantlists.NpcTypes.Spawnentries.Spawngroup",
		"SpellsNew.Items.Merchantlists.NpcTypes.Spawnentries.Spawngroup.Spawn2",
		"SpellsNew.Items.Merchantlists.NpcTypes.Spawnentries.Spawngroup.Spawn2.Spawnentries",
		"SpellsNew.Items.Merchantlists.NpcTypes.Spawnentries.Spawngroup.Spawn2.Spawngroup",
		"SpellsNew.Items.ObjectContents",
		"SpellsNew.Items.Objects",
		"SpellsNew.Items.Objects.Item",
		"SpellsNew.Items.Objects.Zone",
		"SpellsNew.Items.StartingItems",
		"SpellsNew.Items.StartingItems.Item",
		"SpellsNew.Items.StartingItems.Zone",
		"SpellsNew.Items.TradeskillRecipeEntries",
		"SpellsNew.Items.TradeskillRecipeEntries.TradeskillRecipe",
		"SpellsNew.Items.TributeLevels",
		"SpellsNew.NpcSpellsEntries",
		"SpellsNew.NpcSpellsEntries.SpellsNew",
		"SpellsNew.SpellBuckets",
		"SpellsNew.SpellGlobals",
	}
}

func (AaRank) Connection() string {
    return "eqemu_content"
}
