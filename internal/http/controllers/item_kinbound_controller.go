package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/EQEmuTools/spire/internal/auditlog"
	"github.com/EQEmuTools/spire/internal/database"
	"github.com/EQEmuTools/spire/internal/http/routes"
	"github.com/EQEmuTools/spire/internal/models"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Kinbound is a Bastion extension. Keep it outside generated item CRUD so
// ordinary item editing still works against databases without this extension.
type ItemKinboundController struct {
	db       *database.Resolver
	auditLog *auditlog.UserEvent
}

type itemKinboundState struct {
	Supported     bool   `json:"supported" gorm:"-"`
	Message       string `json:"message,omitempty" gorm:"-"`
	ID            int    `json:"id"`
	Kinbound      int    `json:"kinbound"`
	Nodrop        int    `json:"nodrop"`
	NoTransfer    int    `json:"notransfer" gorm:"column:notransfer"`
	Epic          bool   `json:"epic" gorm:"-"`
	CanAlways     bool   `json:"can_always" gorm:"-"`
	BlockedReason string `json:"blocked_reason,omitempty" gorm:"-"`
}

type itemKinboundUpdate struct {
	Kinbound         *int `json:"kinbound"`
	ExpectedKinbound *int `json:"expected_kinbound"`
}

func NewItemKinboundController(db *database.Resolver, auditLog *auditlog.UserEvent) *ItemKinboundController {
	return &ItemKinboundController{db: db, auditLog: auditLog}
}

func (k *ItemKinboundController) Routes() []*routes.Route {
	return []*routes.Route{
		routes.RegisterRoute(http.MethodGet, "item-kinbound/:id", k.get, nil),
		routes.RegisterRoute(http.MethodPatch, "item-kinbound/:id", k.update, nil),
	}
}

func itemKinboundSchemaAvailable(db *gorm.DB) (bool, error) {
	var count int64
	err := db.Raw(`SELECT COUNT(*) FROM information_schema.columns
		WHERE table_schema = DATABASE() AND
		((table_name = 'items' AND column_name = 'kinbound') OR
		(table_name = 'item_kinbound_policy' AND column_name IN ('item_id', 'expansion', 'epic')))`).Scan(&count).Error
	return count == 4, err
}

func loadItemKinbound(db *gorm.DB, id int, lock bool) (itemKinboundState, error) {
	state := itemKinboundState{Supported: true}
	query := db.Table("items").Select("id, kinbound, nodrop, notransfer").Where("id = ?", id)
	if lock {
		query = query.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	if err := query.Take(&state).Error; err != nil {
		return state, err
	}
	// Source defaults explicit selections to non-epic when the catalog loads,
	// but any known epic exclusion vetoes them. Expansion never enables an item.
	var policy struct{ Epic int }
	query = db.Table("item_kinbound_policy").Select("epic").Where("item_id = ?", id)
	if lock {
		query = query.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	if err := query.Take(&policy).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return state, err
	}
	state.Epic = policy.Epic == 1
	state.BlockedReason = itemKinboundBlockedReason(state)
	state.CanAlways = state.BlockedReason == ""
	return state, nil
}

func itemKinboundBlockedReason(state itemKinboundState) string {
	if state.Nodrop != 0 {
		return "Always Kinbound requires an item saved as NO TRADE (No Drop enabled)."
	}
	if state.NoTransfer != 0 {
		return "NO TRANSFER items cannot be Always Kinbound."
	}
	if state.Epic {
		return "This item has a known epic exclusion in the Kinbound catalog."
	}
	return ""
}

func validateItemKinboundUpdate(input itemKinboundUpdate) error {
	if input.Kinbound == nil || *input.Kinbound < 0 || *input.Kinbound > 2 {
		return echo.NewHTTPError(http.StatusBadRequest, "Kinbound must be 0 (Unlisted), 1 (Always), or 2 (Never).")
	}
	if input.ExpectedKinbound == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "The previously loaded Kinbound setting is required. Refresh the item.")
	}
	return nil
}

func itemKinboundID(c echo.Context) (int, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		return 0, echo.NewHTTPError(http.StatusBadRequest, "A positive item ID is required.")
	}
	return id, nil
}

func itemKinboundError(c echo.Context, err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Save this item before editing Kinbound."})
	}
	var httpError *echo.HTTPError
	if errors.As(err, &httpError) {
		return c.JSON(httpError.Code, echo.Map{"error": httpError.Message})
	}
	return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
}

func (k *ItemKinboundController) connection(c echo.Context) (*gorm.DB, error) {
	db := k.db.Get(models.Item{}, c)
	if db == nil {
		return nil, echo.NewHTTPError(http.StatusServiceUnavailable, "The item database is unavailable.")
	}
	return db.WithContext(c.Request().Context()), nil
}

func (k *ItemKinboundController) get(c echo.Context) error {
	id, err := itemKinboundID(c)
	if err != nil {
		return itemKinboundError(c, err)
	}
	db, err := k.connection(c)
	if err != nil {
		return itemKinboundError(c, err)
	}
	available, err := itemKinboundSchemaAvailable(db)
	if err != nil {
		return itemKinboundError(c, err)
	}
	if !available {
		return c.JSON(http.StatusOK, itemKinboundState{ID: id, Message: "Kinbound is unavailable on this database. Install the Source Kinbound catalog and item flag migrations (9397 and 9398)."})
	}
	state, err := loadItemKinbound(db, id, false)
	if err != nil {
		return itemKinboundError(c, err)
	}
	return c.JSON(http.StatusOK, state)
}

func (k *ItemKinboundController) update(c echo.Context) error {
	id, err := itemKinboundID(c)
	if err != nil {
		return itemKinboundError(c, err)
	}
	var input itemKinboundUpdate
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid Kinbound payload."})
	}
	if err := validateItemKinboundUpdate(input); err != nil {
		return itemKinboundError(c, err)
	}
	db, err := k.connection(c)
	if err != nil {
		return itemKinboundError(c, err)
	}
	available, err := itemKinboundSchemaAvailable(db)
	if err != nil {
		return itemKinboundError(c, err)
	}
	if !available {
		return c.JSON(http.StatusConflict, echo.Map{"error": "Kinbound migrations 9397 and 9398 are required on this database."})
	}
	var state itemKinboundState
	var auditID uint
	err = db.Transaction(func(tx *gorm.DB) error {
		var err error
		state, err = loadItemKinbound(tx, id, true)
		if err != nil {
			return err
		}
		if state.Kinbound != *input.ExpectedKinbound {
			return echo.NewHTTPError(http.StatusConflict, "The Kinbound setting changed since it was loaded. Refresh before saving.")
		}
		if *input.Kinbound == 1 && !state.CanAlways {
			return echo.NewHTTPError(http.StatusBadRequest, state.BlockedReason)
		}
		if state.Kinbound == *input.Kinbound {
			return nil
		}
		// Persist the audit before touching item data, including on older item
		// tables whose storage engine may not support transaction rollback.
		auditID, err = k.auditLog.LogEditorEvent(c, "ITEM_KINBOUND_UPDATE",
			fmt.Sprintf("Item %d Kinbound changed from %d to %d", id, state.Kinbound, *input.Kinbound))
		if err != nil {
			return fmt.Errorf("record Kinbound audit: %w", err)
		}
		result := tx.Table("items").Where("id = ? AND kinbound = ? AND nodrop = ? AND notransfer = ?", id, state.Kinbound, state.Nodrop, state.NoTransfer).
			Update("kinbound", *input.Kinbound)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return echo.NewHTTPError(http.StatusConflict, "The saved item flags changed. Refresh before saving Kinbound.")
		}
		state.Kinbound = *input.Kinbound
		return nil
	})
	if err != nil {
		if auditID > 0 {
			_ = k.db.GetSpireDb().Table("spire_user_event_log").Where("id = ?", auditID).Delete(nil).Error
		}
		return itemKinboundError(c, err)
	}
	return c.JSON(http.StatusOK, state)
}
