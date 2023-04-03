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

type BotSpellSettingController struct {
	db       *database.DatabaseResolver
	logger   *logrus.Logger
	auditLog *auditlog.UserEvent
}

func NewBotSpellSettingController(
	db *database.DatabaseResolver,
	logger *logrus.Logger,
	auditLog *auditlog.UserEvent,
) *BotSpellSettingController {
	return &BotSpellSettingController{
		db:       db,
		logger:   logger,
		auditLog: auditLog,
	}
}

func (e *BotSpellSettingController) Routes() []*routes.Route {
	return []*routes.Route{
		routes.RegisterRoute(http.MethodGet, "bot_spell_setting/:id", e.getBotSpellSetting, nil),
		routes.RegisterRoute(http.MethodGet, "bot_spell_settings", e.listBotSpellSettings, nil),
		routes.RegisterRoute(http.MethodGet, "bot_spell_settings/count", e.getBotSpellSettingsCount, nil),
		routes.RegisterRoute(http.MethodPut, "bot_spell_setting", e.createBotSpellSetting, nil),
		routes.RegisterRoute(http.MethodDelete, "bot_spell_setting/:id", e.deleteBotSpellSetting, nil),
		routes.RegisterRoute(http.MethodPatch, "bot_spell_setting/:id", e.updateBotSpellSetting, nil),
		routes.RegisterRoute(http.MethodPost, "bot_spell_settings/bulk", e.getBotSpellSettingsBulk, nil),
	}
}

// listBotSpellSettings godoc
// @Id listBotSpellSettings
// @Summary Lists BotSpellSettings
// @Accept json
// @Produce json
// @Tags BotSpellSetting
// @Param includes query string false "Relationships [all] for all [number] for depth of relationships to load or [.] separated relationship names"
// @Param where query string false "Filter on specific fields. Multiple conditions [.] separated Example: col_like_value.col2__val2"
// @Param whereOr query string false "Filter on specific fields (Chained ors). Multiple conditions [.] separated Example: col_like_value.col2__val2"
// @Param groupBy query string false "Group by field. Multiple conditions [.] separated Example: field1.field2"
// @Param limit query string false "Rows to limit in response (Default: 10,000)"
// @Param page query int 0 "Pagination page"
// @Param orderBy query string false "Order by [field]"
// @Param orderDirection query string false "Order by field direction"
// @Param select query string false "Column names [.] separated to fetch specific fields in response"
// @Success 200 {array} models.BotSpellSetting
// @Failure 500 {string} string "Bad query request"
// @Router /bot_spell_settings [get]
func (e *BotSpellSettingController) listBotSpellSettings(c echo.Context) error {
	var results []models.BotSpellSetting
	err := e.db.QueryContext(models.BotSpellSetting{}, c).Find(&results).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err})
	}

	return c.JSON(http.StatusOK, results)
}

// getBotSpellSetting godoc
// @Id getBotSpellSetting
// @Summary Gets BotSpellSetting
// @Accept json
// @Produce json
// @Tags BotSpellSetting
// @Param id path int true "Id"
// @Param includes query string false "Relationships [all] for all [number] for depth of relationships to load or [.] separated relationship names"
// @Param select query string false "Column names [.] separated to fetch specific fields in response"
// @Success 200 {array} models.BotSpellSetting
// @Failure 404 {string} string "Entity not found"
// @Failure 500 {string} string "Cannot find param"
// @Failure 500 {string} string "Bad query request"
// @Router /bot_spell_setting/{id} [get]
func (e *BotSpellSettingController) getBotSpellSetting(c echo.Context) error {
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
	var result models.BotSpellSetting
	query := e.db.QueryContext(models.BotSpellSetting{}, c)
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

// updateBotSpellSetting godoc
// @Id updateBotSpellSetting
// @Summary Updates BotSpellSetting
// @Accept json
// @Produce json
// @Tags BotSpellSetting
// @Param id path int true "Id"
// @Param bot_spell_setting body models.BotSpellSetting true "BotSpellSetting"
// @Success 200 {array} models.BotSpellSetting
// @Failure 404 {string} string "Cannot find entity"
// @Failure 500 {string} string "Error binding to entity"
// @Failure 500 {string} string "Error updating entity"
// @Router /bot_spell_setting/{id} [patch]
func (e *BotSpellSettingController) updateBotSpellSetting(c echo.Context) error {
	request := new(models.BotSpellSetting)
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
	var result models.BotSpellSetting
	query := e.db.QueryContext(models.BotSpellSetting{}, c)
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
		event := fmt.Sprintf("Updated [BotSpellSetting] [%v] fields [%v]", strings.Join(ids, ", "), strings.Join(fieldsUpdated, ", "))
		e.auditLog.LogUserEvent(c, "UPDATE", event)
	}

	return c.JSON(http.StatusOK, request)
}

// createBotSpellSetting godoc
// @Id createBotSpellSetting
// @Summary Creates BotSpellSetting
// @Accept json
// @Produce json
// @Param bot_spell_setting body models.BotSpellSetting true "BotSpellSetting"
// @Tags BotSpellSetting
// @Success 200 {array} models.BotSpellSetting
// @Failure 500 {string} string "Error binding to entity"
// @Failure 500 {string} string "Error inserting entity"
// @Router /bot_spell_setting [put]
func (e *BotSpellSettingController) createBotSpellSetting(c echo.Context) error {
	botSpellSetting := new(models.BotSpellSetting)
	if err := c.Bind(botSpellSetting); err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			echo.Map{"error": fmt.Sprintf("Error binding to entity [%v]", err.Error())},
		)
	}

	err := e.db.Get(models.BotSpellSetting{}, c).Model(&models.BotSpellSetting{}).Create(&botSpellSetting).Error
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			echo.Map{"error": fmt.Sprintf("Error inserting entity [%v]", err.Error())},
		)
	}

	// log create event
	if e.db.GetSpireDb() != nil {
		// diff between an empty model and the created
		diff := database.ResultDifference(models.BotSpellSetting{}, botSpellSetting)
		// build fields created
		var fields []string
		for k, v := range diff {
			fields = append(fields, fmt.Sprintf("%v = %v", k, v))
		}
		// record event
		event := fmt.Sprintf("Created [BotSpellSetting] [%v] data [%v]", botSpellSetting.ID, strings.Join(fields, ", "))
		e.auditLog.LogUserEvent(c, "CREATE", event)
	}

	return c.JSON(http.StatusOK, botSpellSetting)
}

// deleteBotSpellSetting godoc
// @Id deleteBotSpellSetting
// @Summary Deletes BotSpellSetting
// @Accept json
// @Produce json
// @Tags BotSpellSetting
// @Param id path int true "id"
// @Success 200 {string} string "Entity deleted successfully"
// @Failure 404 {string} string "Cannot find entity"
// @Failure 500 {string} string "Error binding to entity"
// @Failure 500 {string} string "Error deleting entity"
// @Router /bot_spell_setting/{id} [delete]
func (e *BotSpellSettingController) deleteBotSpellSetting(c echo.Context) error {
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
	var result models.BotSpellSetting
	query := e.db.QueryContext(models.BotSpellSetting{}, c)
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
		event := fmt.Sprintf("Deleted [BotSpellSetting] [%v] keys [%v]", result.ID, strings.Join(ids, ", "))
		e.auditLog.LogUserEvent(c, "DELETE", event)
	}

	return c.JSON(http.StatusOK, echo.Map{"success": "Entity deleted successfully"})
}

// getBotSpellSettingsBulk godoc
// @Id getBotSpellSettingsBulk
// @Summary Gets BotSpellSettings in bulk
// @Accept json
// @Produce json
// @Param Body body BulkFetchByIdsGetRequest true "body"
// @Tags BotSpellSetting
// @Success 200 {array} models.BotSpellSetting
// @Failure 500 {string} string "Bad query request"
// @Router /bot_spell_settings/bulk [post]
func (e *BotSpellSettingController) getBotSpellSettingsBulk(c echo.Context) error {
	var results []models.BotSpellSetting

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

	err := e.db.QueryContext(models.BotSpellSetting{}, c).Find(&results, r.IDs).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, results)
}

// getBotSpellSettingsCount godoc
// @Id getBotSpellSettingsCount
// @Summary Counts BotSpellSettings
// @Accept json
// @Produce json
// @Tags BotSpellSetting
// @Param includes query string false "Relationships [all] for all [number] for depth of relationships to load or [.] separated relationship names"
// @Param where query string false "Filter on specific fields. Multiple conditions [.] separated Example: col_like_value.col2__val2"
// @Param whereOr query string false "Filter on specific fields (Chained ors). Multiple conditions [.] separated Example: col_like_value.col2__val2"
// @Param groupBy query string false "Group by field. Multiple conditions [.] separated Example: field1.field2"
// @Param limit query string false "Rows to limit in response (Default: 10,000)"
// @Param page query int 0 "Pagination page"
// @Param orderBy query string false "Order by [field]"
// @Param orderDirection query string false "Order by field direction"
// @Param select query string false "Column names [.] separated to fetch specific fields in response"
// @Success 200 {array} models.BotSpellSetting
// @Failure 500 {string} string "Bad query request"
// @Router /bot_spell_settings/count [get]
func (e *BotSpellSettingController) getBotSpellSettingsCount(c echo.Context) error {
	var count int64
	err := e.db.QueryContext(models.BotSpellSetting{}, c).Count(&count).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err})
	}

	return c.JSON(http.StatusOK, echo.Map{"count": count})
}