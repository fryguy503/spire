package crudcontrollers

import (
	"fmt"
	"github.com/Akkadius/spire/internal/auditlog"
	"github.com/Akkadius/spire/internal/database"
	"github.com/Akkadius/spire/internal/http/routes"
	"github.com/Akkadius/spire/internal/models"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
)

type BotSpellsEntryController struct {
	db       *database.DatabaseResolver
	logger   *logrus.Logger
	auditLog *auditlog.UserEvent
}

func NewBotSpellsEntryController(
	db *database.DatabaseResolver,
	logger *logrus.Logger,
	auditLog *auditlog.UserEvent,
) *BotSpellsEntryController {
	return &BotSpellsEntryController{
		db:       db,
		logger:   logger,
		auditLog: auditLog,
	}
}

func (e *BotSpellsEntryController) Routes() []*routes.Route {
	return []*routes.Route{
		routes.RegisterRoute(http.MethodGet, "bot_spells_entry/:id", e.getBotSpellsEntry, nil),
		routes.RegisterRoute(http.MethodGet, "bot_spells_entries", e.listBotSpellsEntries, nil),
		routes.RegisterRoute(http.MethodPut, "bot_spells_entry", e.createBotSpellsEntry, nil),
		routes.RegisterRoute(http.MethodDelete, "bot_spells_entry/:id", e.deleteBotSpellsEntry, nil),
		routes.RegisterRoute(http.MethodPatch, "bot_spells_entry/:id", e.updateBotSpellsEntry, nil),
		routes.RegisterRoute(http.MethodPost, "bot_spells_entries/bulk", e.getBotSpellsEntriesBulk, nil),
	}
}

// listBotSpellsEntries godoc
// @Id listBotSpellsEntries
// @Summary Lists BotSpellsEntries
// @Accept json
// @Produce json
// @Tags BotSpellsEntry
// @Param includes query string false "Relationships [all] for all [number] for depth of relationships to load or [.] separated relationship names <h4>Relationships</h4>NpcSpell<br>NpcSpell.BotSpellsEntries<br>NpcSpell.NpcSpell<br>NpcSpell.NpcSpellsEntries<br>NpcSpell.NpcSpellsEntries.SpellsNew<br>NpcSpell.NpcSpellsEntries.SpellsNew.Aura<br>NpcSpell.NpcSpellsEntries.SpellsNew.Aura.SpellsNew<br>NpcSpell.NpcSpellsEntries.SpellsNew.BlockedSpells<br>NpcSpell.NpcSpellsEntries.SpellsNew.BotSpellsEntries<br>NpcSpell.NpcSpellsEntries.SpellsNew.Damageshieldtypes<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.AlternateCurrencies<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.AlternateCurrencies.Item<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.CharacterCorpseItems<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.DiscoveredItems<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Doors<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Doors.Item<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings.Item<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings.NpcType<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings.NpcType.AlternateCurrency<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings.NpcType.AlternateCurrency.Item<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings.NpcType.Loottable<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings.NpcType.Loottable.LoottableEntries<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings.NpcType.Loottable.LoottableEntries.Lootdrop<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings.NpcType.Loottable.LoottableEntries.Lootdrop.LootdropEntries<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings.NpcType.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings.NpcType.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Lootdrop<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings.NpcType.Loottable.LoottableEntries.Lootdrop.LoottableEntries<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings.NpcType.Loottable.LoottableEntries.Loottable<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings.NpcType.Loottable.NpcTypes<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings.NpcType.Merchantlists<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings.NpcType.Merchantlists.Items<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings.NpcType.Merchantlists.NpcTypes<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings.NpcType.NpcEmotes<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings.NpcType.NpcFactions<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings.NpcType.NpcFactions.NpcFactionEntries<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings.NpcType.NpcFactions.NpcFactionEntries.FactionList<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings.NpcType.NpcSpell<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings.NpcType.NpcTypesTint<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings.NpcType.Spawnentries<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings.NpcType.Spawnentries.NpcType<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings.NpcType.Spawnentries.Spawngroup<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings.NpcType.Spawnentries.Spawngroup.Spawn2<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings.NpcType.Spawnentries.Spawngroup.Spawn2.Spawnentries<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings.NpcType.Spawnentries.Spawngroup.Spawn2.Spawngroup<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings.Zone<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Forages<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Forages.Item<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Forages.Zone<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.GroundSpawns<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.GroundSpawns.Zone<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.ItemTicks<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Keyrings<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Item<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LootdropEntries<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Lootdrop<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.LoottableEntries<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.AlternateCurrency<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.AlternateCurrency.Item<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Loottable<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Merchantlists<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Merchantlists.Items<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Merchantlists.NpcTypes<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcEmotes<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcFactions<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcFactions.NpcFactionEntries<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcFactions.NpcFactionEntries.FactionList<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcSpell<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcTypesTint<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Spawnentries<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Spawnentries.NpcType<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Spawnentries.Spawngroup<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Spawnentries.Spawngroup.Spawn2<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Spawnentries.Spawngroup.Spawn2.Spawnentries<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Spawnentries.Spawngroup.Spawn2.Spawngroup<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Merchantlists<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Merchantlists.Items<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Merchantlists.NpcTypes<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Merchantlists.NpcTypes.AlternateCurrency<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Merchantlists.NpcTypes.AlternateCurrency.Item<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Merchantlists.NpcTypes.Loottable<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Merchantlists.NpcTypes.Loottable.LoottableEntries<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Merchantlists.NpcTypes.Loottable.LoottableEntries.Lootdrop<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Merchantlists.NpcTypes.Loottable.LoottableEntries.Lootdrop.LootdropEntries<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Merchantlists.NpcTypes.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Merchantlists.NpcTypes.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Lootdrop<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Merchantlists.NpcTypes.Loottable.LoottableEntries.Lootdrop.LoottableEntries<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Merchantlists.NpcTypes.Loottable.LoottableEntries.Loottable<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Merchantlists.NpcTypes.Loottable.NpcTypes<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Merchantlists.NpcTypes.Merchantlists<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Merchantlists.NpcTypes.NpcEmotes<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Merchantlists.NpcTypes.NpcFactions<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Merchantlists.NpcTypes.NpcFactions.NpcFactionEntries<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Merchantlists.NpcTypes.NpcFactions.NpcFactionEntries.FactionList<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Merchantlists.NpcTypes.NpcSpell<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Merchantlists.NpcTypes.NpcTypesTint<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Merchantlists.NpcTypes.Spawnentries<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Merchantlists.NpcTypes.Spawnentries.NpcType<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Merchantlists.NpcTypes.Spawnentries.Spawngroup<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Merchantlists.NpcTypes.Spawnentries.Spawngroup.Spawn2<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Merchantlists.NpcTypes.Spawnentries.Spawngroup.Spawn2.Spawnentries<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Merchantlists.NpcTypes.Spawnentries.Spawngroup.Spawn2.Spawngroup<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.ObjectContents<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Objects<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Objects.Item<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Objects.Zone<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.StartingItems<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.StartingItems.Item<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.StartingItems.Zone<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.TradeskillRecipeEntries<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.TradeskillRecipeEntries.TradeskillRecipe<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.TributeLevels<br>NpcSpell.NpcSpellsEntries.SpellsNew.NpcSpellsEntries<br>NpcSpell.NpcSpellsEntries.SpellsNew.SpellBuckets<br>NpcSpell.NpcSpellsEntries.SpellsNew.SpellGlobals<br>SpellsNew<br>SpellsNew.Aura<br>SpellsNew.Aura.SpellsNew<br>SpellsNew.BlockedSpells<br>SpellsNew.BotSpellsEntries<br>SpellsNew.Damageshieldtypes<br>SpellsNew.Items<br>SpellsNew.Items.AlternateCurrencies<br>SpellsNew.Items.AlternateCurrencies.Item<br>SpellsNew.Items.CharacterCorpseItems<br>SpellsNew.Items.DiscoveredItems<br>SpellsNew.Items.Doors<br>SpellsNew.Items.Doors.Item<br>SpellsNew.Items.Fishings<br>SpellsNew.Items.Fishings.Item<br>SpellsNew.Items.Fishings.NpcType<br>SpellsNew.Items.Fishings.NpcType.AlternateCurrency<br>SpellsNew.Items.Fishings.NpcType.AlternateCurrency.Item<br>SpellsNew.Items.Fishings.NpcType.Loottable<br>SpellsNew.Items.Fishings.NpcType.Loottable.LoottableEntries<br>SpellsNew.Items.Fishings.NpcType.Loottable.LoottableEntries.Lootdrop<br>SpellsNew.Items.Fishings.NpcType.Loottable.LoottableEntries.Lootdrop.LootdropEntries<br>SpellsNew.Items.Fishings.NpcType.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item<br>SpellsNew.Items.Fishings.NpcType.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Lootdrop<br>SpellsNew.Items.Fishings.NpcType.Loottable.LoottableEntries.Lootdrop.LoottableEntries<br>SpellsNew.Items.Fishings.NpcType.Loottable.LoottableEntries.Loottable<br>SpellsNew.Items.Fishings.NpcType.Loottable.NpcTypes<br>SpellsNew.Items.Fishings.NpcType.Merchantlists<br>SpellsNew.Items.Fishings.NpcType.Merchantlists.Items<br>SpellsNew.Items.Fishings.NpcType.Merchantlists.NpcTypes<br>SpellsNew.Items.Fishings.NpcType.NpcEmotes<br>SpellsNew.Items.Fishings.NpcType.NpcFactions<br>SpellsNew.Items.Fishings.NpcType.NpcFactions.NpcFactionEntries<br>SpellsNew.Items.Fishings.NpcType.NpcFactions.NpcFactionEntries.FactionList<br>SpellsNew.Items.Fishings.NpcType.NpcSpell<br>SpellsNew.Items.Fishings.NpcType.NpcSpell.BotSpellsEntries<br>SpellsNew.Items.Fishings.NpcType.NpcSpell.NpcSpell<br>SpellsNew.Items.Fishings.NpcType.NpcSpell.NpcSpellsEntries<br>SpellsNew.Items.Fishings.NpcType.NpcSpell.NpcSpellsEntries.SpellsNew<br>SpellsNew.Items.Fishings.NpcType.NpcTypesTint<br>SpellsNew.Items.Fishings.NpcType.Spawnentries<br>SpellsNew.Items.Fishings.NpcType.Spawnentries.NpcType<br>SpellsNew.Items.Fishings.NpcType.Spawnentries.Spawngroup<br>SpellsNew.Items.Fishings.NpcType.Spawnentries.Spawngroup.Spawn2<br>SpellsNew.Items.Fishings.NpcType.Spawnentries.Spawngroup.Spawn2.Spawnentries<br>SpellsNew.Items.Fishings.NpcType.Spawnentries.Spawngroup.Spawn2.Spawngroup<br>SpellsNew.Items.Fishings.Zone<br>SpellsNew.Items.Forages<br>SpellsNew.Items.Forages.Item<br>SpellsNew.Items.Forages.Zone<br>SpellsNew.Items.GroundSpawns<br>SpellsNew.Items.GroundSpawns.Zone<br>SpellsNew.Items.ItemTicks<br>SpellsNew.Items.Keyrings<br>SpellsNew.Items.LootdropEntries<br>SpellsNew.Items.LootdropEntries.Item<br>SpellsNew.Items.LootdropEntries.Lootdrop<br>SpellsNew.Items.LootdropEntries.Lootdrop.LootdropEntries<br>SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries<br>SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Lootdrop<br>SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable<br>SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.LoottableEntries<br>SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes<br>SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.AlternateCurrency<br>SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.AlternateCurrency.Item<br>SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Loottable<br>SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Merchantlists<br>SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Merchantlists.Items<br>SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Merchantlists.NpcTypes<br>SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcEmotes<br>SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcFactions<br>SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcFactions.NpcFactionEntries<br>SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcFactions.NpcFactionEntries.FactionList<br>SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcSpell<br>SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcSpell.BotSpellsEntries<br>SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcSpell.NpcSpell<br>SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcSpell.NpcSpellsEntries<br>SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew<br>SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcTypesTint<br>SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Spawnentries<br>SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Spawnentries.NpcType<br>SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Spawnentries.Spawngroup<br>SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Spawnentries.Spawngroup.Spawn2<br>SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Spawnentries.Spawngroup.Spawn2.Spawnentries<br>SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Spawnentries.Spawngroup.Spawn2.Spawngroup<br>SpellsNew.Items.Merchantlists<br>SpellsNew.Items.Merchantlists.Items<br>SpellsNew.Items.Merchantlists.NpcTypes<br>SpellsNew.Items.Merchantlists.NpcTypes.AlternateCurrency<br>SpellsNew.Items.Merchantlists.NpcTypes.AlternateCurrency.Item<br>SpellsNew.Items.Merchantlists.NpcTypes.Loottable<br>SpellsNew.Items.Merchantlists.NpcTypes.Loottable.LoottableEntries<br>SpellsNew.Items.Merchantlists.NpcTypes.Loottable.LoottableEntries.Lootdrop<br>SpellsNew.Items.Merchantlists.NpcTypes.Loottable.LoottableEntries.Lootdrop.LootdropEntries<br>SpellsNew.Items.Merchantlists.NpcTypes.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item<br>SpellsNew.Items.Merchantlists.NpcTypes.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Lootdrop<br>SpellsNew.Items.Merchantlists.NpcTypes.Loottable.LoottableEntries.Lootdrop.LoottableEntries<br>SpellsNew.Items.Merchantlists.NpcTypes.Loottable.LoottableEntries.Loottable<br>SpellsNew.Items.Merchantlists.NpcTypes.Loottable.NpcTypes<br>SpellsNew.Items.Merchantlists.NpcTypes.Merchantlists<br>SpellsNew.Items.Merchantlists.NpcTypes.NpcEmotes<br>SpellsNew.Items.Merchantlists.NpcTypes.NpcFactions<br>SpellsNew.Items.Merchantlists.NpcTypes.NpcFactions.NpcFactionEntries<br>SpellsNew.Items.Merchantlists.NpcTypes.NpcFactions.NpcFactionEntries.FactionList<br>SpellsNew.Items.Merchantlists.NpcTypes.NpcSpell<br>SpellsNew.Items.Merchantlists.NpcTypes.NpcSpell.BotSpellsEntries<br>SpellsNew.Items.Merchantlists.NpcTypes.NpcSpell.NpcSpell<br>SpellsNew.Items.Merchantlists.NpcTypes.NpcSpell.NpcSpellsEntries<br>SpellsNew.Items.Merchantlists.NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew<br>SpellsNew.Items.Merchantlists.NpcTypes.NpcTypesTint<br>SpellsNew.Items.Merchantlists.NpcTypes.Spawnentries<br>SpellsNew.Items.Merchantlists.NpcTypes.Spawnentries.NpcType<br>SpellsNew.Items.Merchantlists.NpcTypes.Spawnentries.Spawngroup<br>SpellsNew.Items.Merchantlists.NpcTypes.Spawnentries.Spawngroup.Spawn2<br>SpellsNew.Items.Merchantlists.NpcTypes.Spawnentries.Spawngroup.Spawn2.Spawnentries<br>SpellsNew.Items.Merchantlists.NpcTypes.Spawnentries.Spawngroup.Spawn2.Spawngroup<br>SpellsNew.Items.ObjectContents<br>SpellsNew.Items.Objects<br>SpellsNew.Items.Objects.Item<br>SpellsNew.Items.Objects.Zone<br>SpellsNew.Items.StartingItems<br>SpellsNew.Items.StartingItems.Item<br>SpellsNew.Items.StartingItems.Zone<br>SpellsNew.Items.TradeskillRecipeEntries<br>SpellsNew.Items.TradeskillRecipeEntries.TradeskillRecipe<br>SpellsNew.Items.TributeLevels<br>SpellsNew.NpcSpellsEntries<br>SpellsNew.NpcSpellsEntries.SpellsNew<br>SpellsNew.SpellBuckets<br>SpellsNew.SpellGlobals"
// @Param where query string false "Filter on specific fields. Multiple conditions [.] separated Example: col_like_value.col2__val2"
// @Param whereOr query string false "Filter on specific fields (Chained ors). Multiple conditions [.] separated Example: col_like_value.col2__val2"
// @Param groupBy query string false "Group by field. Multiple conditions [.] separated Example: field1.field2"
// @Param limit query string false "Rows to limit in response (Default: 10,000)"
// @Param page query int 0 "Pagination page"
// @Param orderBy query string false "Order by [field]"
// @Param orderDirection query string false "Order by field direction"
// @Param select query string false "Column names [.] separated to fetch specific fields in response"
// @Success 200 {array} models.BotSpellsEntry
// @Failure 500 {string} string "Bad query request"
// @Router /bot_spells_entries [get]
func (e *BotSpellsEntryController) listBotSpellsEntries(c echo.Context) error {
	var results []models.BotSpellsEntry
	err := e.db.QueryContext(models.BotSpellsEntry{}, c).Find(&results).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err})
	}

	return c.JSON(http.StatusOK, results)
}

// getBotSpellsEntry godoc
// @Id getBotSpellsEntry
// @Summary Gets BotSpellsEntry
// @Accept json
// @Produce json
// @Tags BotSpellsEntry
// @Param id path int true "Id"
// @Param includes query string false "Relationships [all] for all [number] for depth of relationships to load or [.] separated relationship names <h4>Relationships</h4>NpcSpell<br>NpcSpell.BotSpellsEntries<br>NpcSpell.NpcSpell<br>NpcSpell.NpcSpellsEntries<br>NpcSpell.NpcSpellsEntries.SpellsNew<br>NpcSpell.NpcSpellsEntries.SpellsNew.Aura<br>NpcSpell.NpcSpellsEntries.SpellsNew.Aura.SpellsNew<br>NpcSpell.NpcSpellsEntries.SpellsNew.BlockedSpells<br>NpcSpell.NpcSpellsEntries.SpellsNew.BotSpellsEntries<br>NpcSpell.NpcSpellsEntries.SpellsNew.Damageshieldtypes<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.AlternateCurrencies<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.AlternateCurrencies.Item<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.CharacterCorpseItems<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.DiscoveredItems<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Doors<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Doors.Item<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings.Item<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings.NpcType<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings.NpcType.AlternateCurrency<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings.NpcType.AlternateCurrency.Item<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings.NpcType.Loottable<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings.NpcType.Loottable.LoottableEntries<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings.NpcType.Loottable.LoottableEntries.Lootdrop<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings.NpcType.Loottable.LoottableEntries.Lootdrop.LootdropEntries<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings.NpcType.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings.NpcType.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Lootdrop<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings.NpcType.Loottable.LoottableEntries.Lootdrop.LoottableEntries<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings.NpcType.Loottable.LoottableEntries.Loottable<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings.NpcType.Loottable.NpcTypes<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings.NpcType.Merchantlists<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings.NpcType.Merchantlists.Items<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings.NpcType.Merchantlists.NpcTypes<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings.NpcType.NpcEmotes<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings.NpcType.NpcFactions<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings.NpcType.NpcFactions.NpcFactionEntries<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings.NpcType.NpcFactions.NpcFactionEntries.FactionList<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings.NpcType.NpcSpell<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings.NpcType.NpcTypesTint<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings.NpcType.Spawnentries<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings.NpcType.Spawnentries.NpcType<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings.NpcType.Spawnentries.Spawngroup<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings.NpcType.Spawnentries.Spawngroup.Spawn2<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings.NpcType.Spawnentries.Spawngroup.Spawn2.Spawnentries<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings.NpcType.Spawnentries.Spawngroup.Spawn2.Spawngroup<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Fishings.Zone<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Forages<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Forages.Item<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Forages.Zone<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.GroundSpawns<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.GroundSpawns.Zone<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.ItemTicks<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Keyrings<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Item<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LootdropEntries<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Lootdrop<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.LoottableEntries<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.AlternateCurrency<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.AlternateCurrency.Item<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Loottable<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Merchantlists<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Merchantlists.Items<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Merchantlists.NpcTypes<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcEmotes<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcFactions<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcFactions.NpcFactionEntries<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcFactions.NpcFactionEntries.FactionList<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcSpell<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcTypesTint<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Spawnentries<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Spawnentries.NpcType<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Spawnentries.Spawngroup<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Spawnentries.Spawngroup.Spawn2<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Spawnentries.Spawngroup.Spawn2.Spawnentries<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Spawnentries.Spawngroup.Spawn2.Spawngroup<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Merchantlists<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Merchantlists.Items<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Merchantlists.NpcTypes<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Merchantlists.NpcTypes.AlternateCurrency<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Merchantlists.NpcTypes.AlternateCurrency.Item<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Merchantlists.NpcTypes.Loottable<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Merchantlists.NpcTypes.Loottable.LoottableEntries<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Merchantlists.NpcTypes.Loottable.LoottableEntries.Lootdrop<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Merchantlists.NpcTypes.Loottable.LoottableEntries.Lootdrop.LootdropEntries<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Merchantlists.NpcTypes.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Merchantlists.NpcTypes.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Lootdrop<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Merchantlists.NpcTypes.Loottable.LoottableEntries.Lootdrop.LoottableEntries<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Merchantlists.NpcTypes.Loottable.LoottableEntries.Loottable<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Merchantlists.NpcTypes.Loottable.NpcTypes<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Merchantlists.NpcTypes.Merchantlists<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Merchantlists.NpcTypes.NpcEmotes<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Merchantlists.NpcTypes.NpcFactions<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Merchantlists.NpcTypes.NpcFactions.NpcFactionEntries<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Merchantlists.NpcTypes.NpcFactions.NpcFactionEntries.FactionList<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Merchantlists.NpcTypes.NpcSpell<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Merchantlists.NpcTypes.NpcTypesTint<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Merchantlists.NpcTypes.Spawnentries<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Merchantlists.NpcTypes.Spawnentries.NpcType<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Merchantlists.NpcTypes.Spawnentries.Spawngroup<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Merchantlists.NpcTypes.Spawnentries.Spawngroup.Spawn2<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Merchantlists.NpcTypes.Spawnentries.Spawngroup.Spawn2.Spawnentries<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Merchantlists.NpcTypes.Spawnentries.Spawngroup.Spawn2.Spawngroup<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.ObjectContents<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Objects<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Objects.Item<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.Objects.Zone<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.StartingItems<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.StartingItems.Item<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.StartingItems.Zone<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.TradeskillRecipeEntries<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.TradeskillRecipeEntries.TradeskillRecipe<br>NpcSpell.NpcSpellsEntries.SpellsNew.Items.TributeLevels<br>NpcSpell.NpcSpellsEntries.SpellsNew.NpcSpellsEntries<br>NpcSpell.NpcSpellsEntries.SpellsNew.SpellBuckets<br>NpcSpell.NpcSpellsEntries.SpellsNew.SpellGlobals<br>SpellsNew<br>SpellsNew.Aura<br>SpellsNew.Aura.SpellsNew<br>SpellsNew.BlockedSpells<br>SpellsNew.BotSpellsEntries<br>SpellsNew.Damageshieldtypes<br>SpellsNew.Items<br>SpellsNew.Items.AlternateCurrencies<br>SpellsNew.Items.AlternateCurrencies.Item<br>SpellsNew.Items.CharacterCorpseItems<br>SpellsNew.Items.DiscoveredItems<br>SpellsNew.Items.Doors<br>SpellsNew.Items.Doors.Item<br>SpellsNew.Items.Fishings<br>SpellsNew.Items.Fishings.Item<br>SpellsNew.Items.Fishings.NpcType<br>SpellsNew.Items.Fishings.NpcType.AlternateCurrency<br>SpellsNew.Items.Fishings.NpcType.AlternateCurrency.Item<br>SpellsNew.Items.Fishings.NpcType.Loottable<br>SpellsNew.Items.Fishings.NpcType.Loottable.LoottableEntries<br>SpellsNew.Items.Fishings.NpcType.Loottable.LoottableEntries.Lootdrop<br>SpellsNew.Items.Fishings.NpcType.Loottable.LoottableEntries.Lootdrop.LootdropEntries<br>SpellsNew.Items.Fishings.NpcType.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item<br>SpellsNew.Items.Fishings.NpcType.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Lootdrop<br>SpellsNew.Items.Fishings.NpcType.Loottable.LoottableEntries.Lootdrop.LoottableEntries<br>SpellsNew.Items.Fishings.NpcType.Loottable.LoottableEntries.Loottable<br>SpellsNew.Items.Fishings.NpcType.Loottable.NpcTypes<br>SpellsNew.Items.Fishings.NpcType.Merchantlists<br>SpellsNew.Items.Fishings.NpcType.Merchantlists.Items<br>SpellsNew.Items.Fishings.NpcType.Merchantlists.NpcTypes<br>SpellsNew.Items.Fishings.NpcType.NpcEmotes<br>SpellsNew.Items.Fishings.NpcType.NpcFactions<br>SpellsNew.Items.Fishings.NpcType.NpcFactions.NpcFactionEntries<br>SpellsNew.Items.Fishings.NpcType.NpcFactions.NpcFactionEntries.FactionList<br>SpellsNew.Items.Fishings.NpcType.NpcSpell<br>SpellsNew.Items.Fishings.NpcType.NpcSpell.BotSpellsEntries<br>SpellsNew.Items.Fishings.NpcType.NpcSpell.NpcSpell<br>SpellsNew.Items.Fishings.NpcType.NpcSpell.NpcSpellsEntries<br>SpellsNew.Items.Fishings.NpcType.NpcSpell.NpcSpellsEntries.SpellsNew<br>SpellsNew.Items.Fishings.NpcType.NpcTypesTint<br>SpellsNew.Items.Fishings.NpcType.Spawnentries<br>SpellsNew.Items.Fishings.NpcType.Spawnentries.NpcType<br>SpellsNew.Items.Fishings.NpcType.Spawnentries.Spawngroup<br>SpellsNew.Items.Fishings.NpcType.Spawnentries.Spawngroup.Spawn2<br>SpellsNew.Items.Fishings.NpcType.Spawnentries.Spawngroup.Spawn2.Spawnentries<br>SpellsNew.Items.Fishings.NpcType.Spawnentries.Spawngroup.Spawn2.Spawngroup<br>SpellsNew.Items.Fishings.Zone<br>SpellsNew.Items.Forages<br>SpellsNew.Items.Forages.Item<br>SpellsNew.Items.Forages.Zone<br>SpellsNew.Items.GroundSpawns<br>SpellsNew.Items.GroundSpawns.Zone<br>SpellsNew.Items.ItemTicks<br>SpellsNew.Items.Keyrings<br>SpellsNew.Items.LootdropEntries<br>SpellsNew.Items.LootdropEntries.Item<br>SpellsNew.Items.LootdropEntries.Lootdrop<br>SpellsNew.Items.LootdropEntries.Lootdrop.LootdropEntries<br>SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries<br>SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Lootdrop<br>SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable<br>SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.LoottableEntries<br>SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes<br>SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.AlternateCurrency<br>SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.AlternateCurrency.Item<br>SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Loottable<br>SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Merchantlists<br>SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Merchantlists.Items<br>SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Merchantlists.NpcTypes<br>SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcEmotes<br>SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcFactions<br>SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcFactions.NpcFactionEntries<br>SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcFactions.NpcFactionEntries.FactionList<br>SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcSpell<br>SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcSpell.BotSpellsEntries<br>SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcSpell.NpcSpell<br>SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcSpell.NpcSpellsEntries<br>SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew<br>SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.NpcTypesTint<br>SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Spawnentries<br>SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Spawnentries.NpcType<br>SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Spawnentries.Spawngroup<br>SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Spawnentries.Spawngroup.Spawn2<br>SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Spawnentries.Spawngroup.Spawn2.Spawnentries<br>SpellsNew.Items.LootdropEntries.Lootdrop.LoottableEntries.Loottable.NpcTypes.Spawnentries.Spawngroup.Spawn2.Spawngroup<br>SpellsNew.Items.Merchantlists<br>SpellsNew.Items.Merchantlists.Items<br>SpellsNew.Items.Merchantlists.NpcTypes<br>SpellsNew.Items.Merchantlists.NpcTypes.AlternateCurrency<br>SpellsNew.Items.Merchantlists.NpcTypes.AlternateCurrency.Item<br>SpellsNew.Items.Merchantlists.NpcTypes.Loottable<br>SpellsNew.Items.Merchantlists.NpcTypes.Loottable.LoottableEntries<br>SpellsNew.Items.Merchantlists.NpcTypes.Loottable.LoottableEntries.Lootdrop<br>SpellsNew.Items.Merchantlists.NpcTypes.Loottable.LoottableEntries.Lootdrop.LootdropEntries<br>SpellsNew.Items.Merchantlists.NpcTypes.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Item<br>SpellsNew.Items.Merchantlists.NpcTypes.Loottable.LoottableEntries.Lootdrop.LootdropEntries.Lootdrop<br>SpellsNew.Items.Merchantlists.NpcTypes.Loottable.LoottableEntries.Lootdrop.LoottableEntries<br>SpellsNew.Items.Merchantlists.NpcTypes.Loottable.LoottableEntries.Loottable<br>SpellsNew.Items.Merchantlists.NpcTypes.Loottable.NpcTypes<br>SpellsNew.Items.Merchantlists.NpcTypes.Merchantlists<br>SpellsNew.Items.Merchantlists.NpcTypes.NpcEmotes<br>SpellsNew.Items.Merchantlists.NpcTypes.NpcFactions<br>SpellsNew.Items.Merchantlists.NpcTypes.NpcFactions.NpcFactionEntries<br>SpellsNew.Items.Merchantlists.NpcTypes.NpcFactions.NpcFactionEntries.FactionList<br>SpellsNew.Items.Merchantlists.NpcTypes.NpcSpell<br>SpellsNew.Items.Merchantlists.NpcTypes.NpcSpell.BotSpellsEntries<br>SpellsNew.Items.Merchantlists.NpcTypes.NpcSpell.NpcSpell<br>SpellsNew.Items.Merchantlists.NpcTypes.NpcSpell.NpcSpellsEntries<br>SpellsNew.Items.Merchantlists.NpcTypes.NpcSpell.NpcSpellsEntries.SpellsNew<br>SpellsNew.Items.Merchantlists.NpcTypes.NpcTypesTint<br>SpellsNew.Items.Merchantlists.NpcTypes.Spawnentries<br>SpellsNew.Items.Merchantlists.NpcTypes.Spawnentries.NpcType<br>SpellsNew.Items.Merchantlists.NpcTypes.Spawnentries.Spawngroup<br>SpellsNew.Items.Merchantlists.NpcTypes.Spawnentries.Spawngroup.Spawn2<br>SpellsNew.Items.Merchantlists.NpcTypes.Spawnentries.Spawngroup.Spawn2.Spawnentries<br>SpellsNew.Items.Merchantlists.NpcTypes.Spawnentries.Spawngroup.Spawn2.Spawngroup<br>SpellsNew.Items.ObjectContents<br>SpellsNew.Items.Objects<br>SpellsNew.Items.Objects.Item<br>SpellsNew.Items.Objects.Zone<br>SpellsNew.Items.StartingItems<br>SpellsNew.Items.StartingItems.Item<br>SpellsNew.Items.StartingItems.Zone<br>SpellsNew.Items.TradeskillRecipeEntries<br>SpellsNew.Items.TradeskillRecipeEntries.TradeskillRecipe<br>SpellsNew.Items.TributeLevels<br>SpellsNew.NpcSpellsEntries<br>SpellsNew.NpcSpellsEntries.SpellsNew<br>SpellsNew.SpellBuckets<br>SpellsNew.SpellGlobals"
// @Param select query string false "Column names [.] separated to fetch specific fields in response"
// @Success 200 {array} models.BotSpellsEntry
// @Failure 404 {string} string "Entity not found"
// @Failure 500 {string} string "Cannot find param"
// @Failure 500 {string} string "Bad query request"
// @Router /bot_spells_entry/{id} [get]
func (e *BotSpellsEntryController) getBotSpellsEntry(c echo.Context) error {
	var params []interface{}
	var keys []string

	// primary key param
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Cannot find param [Id]"})
	}
	params = append(params, id)
	keys = append(keys, "id = ?")

	// query builder
	var result models.BotSpellsEntry
	query := e.db.QueryContext(models.BotSpellsEntry{}, c)
	for i, _ := range keys {
		query = query.Where(keys[i], params[i])
	}

	// grab first entry
	err = query.First(&result).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	// couldn't find entity
	if result.ID == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Cannot find entity"})
	}

	return c.JSON(http.StatusOK, result)
}

// updateBotSpellsEntry godoc
// @Id updateBotSpellsEntry
// @Summary Updates BotSpellsEntry
// @Accept json
// @Produce json
// @Tags BotSpellsEntry
// @Param id path int true "Id"
// @Param bot_spells_entry body models.BotSpellsEntry true "BotSpellsEntry"
// @Success 200 {array} models.BotSpellsEntry
// @Failure 404 {string} string "Cannot find entity"
// @Failure 500 {string} string "Error binding to entity"
// @Failure 500 {string} string "Error updating entity"
// @Router /bot_spells_entry/{id} [patch]
func (e *BotSpellsEntryController) updateBotSpellsEntry(c echo.Context) error {
	request := new(models.BotSpellsEntry)
	if err := c.Bind(request); err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			echo.Map{"error": fmt.Sprintf("Error binding to entity [%v]", err.Error())},
		)
	}

	var params []interface{}
	var keys []string

	// primary key param
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Cannot find param [Id]"})
	}
	params = append(params, id)
	keys = append(keys, "id = ?")

	// query builder
	var result models.BotSpellsEntry
	query := e.db.QueryContext(models.BotSpellsEntry{}, c)
	for i, _ := range keys {
		query = query.Where(keys[i], params[i])
	}

	// grab first entry
	err = query.First(&result).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("Cannot find entity [%s]", err.Error())})
	}

	// save top-level using only changes
	diff := database.ResultDifference(result, request)
	err = query.Session(&gorm.Session{FullSaveAssociations: false}).Updates(diff).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("Error updating entity [%v]", err.Error())})
	}

	// log update event
	if e.db.GetSpireDb() != nil && len(diff) > 0 {
		// build ids
		var ids []string
		for i, _ := range keys {
			param := fmt.Sprintf("%v", params[i])
			ids = append(ids, fmt.Sprintf("%v", strings.ReplaceAll(keys[i], "?", param)))
		}
		// build fields updated
		var fieldsUpdated []string
		for k, v := range diff {
			fieldsUpdated = append(fieldsUpdated, fmt.Sprintf("%v = %v", k, v))
		}
		// record event
		event := fmt.Sprintf("Updated [BotSpellsEntry] [%v] fields [%v]", strings.Join(ids, ", "), strings.Join(fieldsUpdated, ", "))
		e.auditLog.LogUserEvent(c, "UPDATE", event)
	}

	return c.JSON(http.StatusOK, request)
}

// createBotSpellsEntry godoc
// @Id createBotSpellsEntry
// @Summary Creates BotSpellsEntry
// @Accept json
// @Produce json
// @Param bot_spells_entry body models.BotSpellsEntry true "BotSpellsEntry"
// @Tags BotSpellsEntry
// @Success 200 {array} models.BotSpellsEntry
// @Failure 500 {string} string "Error binding to entity"
// @Failure 500 {string} string "Error inserting entity"
// @Router /bot_spells_entry [put]
func (e *BotSpellsEntryController) createBotSpellsEntry(c echo.Context) error {
	botSpellsEntry := new(models.BotSpellsEntry)
	if err := c.Bind(botSpellsEntry); err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			echo.Map{"error": fmt.Sprintf("Error binding to entity [%v]", err.Error())},
		)
	}

	err := e.db.Get(models.BotSpellsEntry{}, c).Model(&models.BotSpellsEntry{}).Create(&botSpellsEntry).Error
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			echo.Map{"error": fmt.Sprintf("Error inserting entity [%v]", err.Error())},
		)
	}

	// log create event
	if e.db.GetSpireDb() != nil {
		// diff between an empty model and the created
		diff := database.ResultDifference(models.BotSpellsEntry{}, botSpellsEntry)
		// build fields created
		var fields []string
		for k, v := range diff {
			fields = append(fields, fmt.Sprintf("%v = %v", k, v))
		}
		// record event
		event := fmt.Sprintf("Created [BotSpellsEntry] [%v] data [%v]", botSpellsEntry.ID, strings.Join(fields, ", "))
		e.auditLog.LogUserEvent(c, "CREATE", event)
	}

	return c.JSON(http.StatusOK, botSpellsEntry)
}

// deleteBotSpellsEntry godoc
// @Id deleteBotSpellsEntry
// @Summary Deletes BotSpellsEntry
// @Accept json
// @Produce json
// @Tags BotSpellsEntry
// @Param id path int true "id"
// @Success 200 {string} string "Entity deleted successfully"
// @Failure 404 {string} string "Cannot find entity"
// @Failure 500 {string} string "Error binding to entity"
// @Failure 500 {string} string "Error deleting entity"
// @Router /bot_spells_entry/{id} [delete]
func (e *BotSpellsEntryController) deleteBotSpellsEntry(c echo.Context) error {
	var params []interface{}
	var keys []string

	// primary key param
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		e.logger.Error(err)
	}
	params = append(params, id)
	keys = append(keys, "id = ?")

	// query builder
	var result models.BotSpellsEntry
	query := e.db.QueryContext(models.BotSpellsEntry{}, c)
	for i, _ := range keys {
		query = query.Where(keys[i], params[i])
	}

	// grab first entry
	err = query.First(&result).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	err = query.Limit(10000).Delete(&result).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Error deleting entity"})
	}

	// log delete event
	if e.db.GetSpireDb() != nil {
		// build ids
		var ids []string
		for i, _ := range keys {
			param := fmt.Sprintf("%v", params[i])
			ids = append(ids, fmt.Sprintf("%v", strings.ReplaceAll(keys[i], "?", param)))
		}
		// record event
		event := fmt.Sprintf("Deleted [BotSpellsEntry] [%v] keys [%v]", result.ID, strings.Join(ids, ", "))
		e.auditLog.LogUserEvent(c, "DELETE", event)
	}

	return c.JSON(http.StatusOK, echo.Map{"success": "Entity deleted successfully"})
}

// getBotSpellsEntriesBulk godoc
// @Id getBotSpellsEntriesBulk
// @Summary Gets BotSpellsEntries in bulk
// @Accept json
// @Produce json
// @Param Body body BulkFetchByIdsGetRequest true "body"
// @Tags BotSpellsEntry
// @Success 200 {array} models.BotSpellsEntry
// @Failure 500 {string} string "Bad query request"
// @Router /bot_spells_entries/bulk [post]
func (e *BotSpellsEntryController) getBotSpellsEntriesBulk(c echo.Context) error {
	var results []models.BotSpellsEntry

	r := new(BulkFetchByIdsGetRequest)
	if err := c.Bind(r); err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			echo.Map{"error": fmt.Sprintf("Error binding to bulk request: [%v]", err.Error())},
		)
	}

	if len(r.IDs) == 0 {
		return c.JSON(
			http.StatusOK,
			echo.Map{"error": fmt.Sprintf("Missing request field data 'ids'")},
		)
	}

	err := e.db.QueryContext(models.BotSpellsEntry{}, c).Find(&results, r.IDs).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, results)
}