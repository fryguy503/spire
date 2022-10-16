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

type PerlEventExportSettingController struct {
	db       *database.DatabaseResolver
	logger   *logrus.Logger
	auditLog *auditlog.UserEvent
}

func NewPerlEventExportSettingController(
	db *database.DatabaseResolver,
	logger *logrus.Logger,
	auditLog *auditlog.UserEvent,
) *PerlEventExportSettingController {
	return &PerlEventExportSettingController{
		db:       db,
		logger:   logger,
		auditLog: auditLog,
	}
}

func (e *PerlEventExportSettingController) Routes() []*routes.Route {
	return []*routes.Route{
		routes.RegisterRoute(http.MethodGet, "perl_event_export_setting/:eventId", e.getPerlEventExportSetting, nil),
		routes.RegisterRoute(http.MethodGet, "perl_event_export_settings", e.listPerlEventExportSettings, nil),
		routes.RegisterRoute(http.MethodPut, "perl_event_export_setting", e.createPerlEventExportSetting, nil),
		routes.RegisterRoute(http.MethodDelete, "perl_event_export_setting/:eventId", e.deletePerlEventExportSetting, nil),
		routes.RegisterRoute(http.MethodPatch, "perl_event_export_setting/:eventId", e.updatePerlEventExportSetting, nil),
		routes.RegisterRoute(http.MethodPost, "perl_event_export_settings/bulk", e.getPerlEventExportSettingsBulk, nil),
	}
}

// listPerlEventExportSettings godoc
// @Id listPerlEventExportSettings
// @Summary Lists PerlEventExportSettings
// @Accept json
// @Produce json
// @Tags PerlEventExportSetting
// @Param includes query string false "Relationships [all] for all [number] for depth of relationships to load or [.] separated relationship names "
// @Param where query string false "Filter on specific fields. Multiple conditions [.] separated Example: col_like_value.col2__val2"
// @Param whereOr query string false "Filter on specific fields (Chained ors). Multiple conditions [.] separated Example: col_like_value.col2__val2"
// @Param groupBy query string false "Group by field. Multiple conditions [.] separated Example: field1.field2"
// @Param limit query string false "Rows to limit in response (Default: 10,000)"
// @Param page query int 0 "Pagination page"
// @Param orderBy query string false "Order by [field]"
// @Param orderDirection query string false "Order by field direction"
// @Param select query string false "Column names [.] separated to fetch specific fields in response"
// @Success 200 {array} models.PerlEventExportSetting
// @Failure 500 {string} string "Bad query request"
// @Router /perl_event_export_settings [get]
func (e *PerlEventExportSettingController) listPerlEventExportSettings(c echo.Context) error {
	var results []models.PerlEventExportSetting
	err := e.db.QueryContext(models.PerlEventExportSetting{}, c).Find(&results).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err})
	}

	return c.JSON(http.StatusOK, results)
}

// getPerlEventExportSetting godoc
// @Id getPerlEventExportSetting
// @Summary Gets PerlEventExportSetting
// @Accept json
// @Produce json
// @Tags PerlEventExportSetting
// @Param id path int true "Id"
// @Param includes query string false "Relationships [all] for all [number] for depth of relationships to load or [.] separated relationship names "
// @Param select query string false "Column names [.] separated to fetch specific fields in response"
// @Success 200 {array} models.PerlEventExportSetting
// @Failure 404 {string} string "Entity not found"
// @Failure 500 {string} string "Cannot find param"
// @Failure 500 {string} string "Bad query request"
// @Router /perl_event_export_setting/{id} [get]
func (e *PerlEventExportSettingController) getPerlEventExportSetting(c echo.Context) error {
	var params []interface{}
	var keys []string

	// primary key param
	eventId, err := strconv.Atoi(c.Param("eventId"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Cannot find param [EventId]"})
	}
	params = append(params, eventId)
	keys = append(keys, "event_id = ?")

	// query builder
	var result models.PerlEventExportSetting
	query := e.db.QueryContext(models.PerlEventExportSetting{}, c)
	for i, _ := range keys {
		query = query.Where(keys[i], params[i])
	}

	// grab first entry
	err = query.First(&result).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	// couldn't find entity
	if result.EventId == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Cannot find entity"})
	}

	return c.JSON(http.StatusOK, result)
}

// updatePerlEventExportSetting godoc
// @Id updatePerlEventExportSetting
// @Summary Updates PerlEventExportSetting
// @Accept json
// @Produce json
// @Tags PerlEventExportSetting
// @Param id path int true "Id"
// @Param perl_event_export_setting body models.PerlEventExportSetting true "PerlEventExportSetting"
// @Success 200 {array} models.PerlEventExportSetting
// @Failure 404 {string} string "Cannot find entity"
// @Failure 500 {string} string "Error binding to entity"
// @Failure 500 {string} string "Error updating entity"
// @Router /perl_event_export_setting/{id} [patch]
func (e *PerlEventExportSettingController) updatePerlEventExportSetting(c echo.Context) error {
	request := new(models.PerlEventExportSetting)
	if err := c.Bind(request); err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			echo.Map{"error": fmt.Sprintf("Error binding to entity [%v]", err.Error())},
		)
	}

	var params []interface{}
	var keys []string

	// primary key param
	eventId, err := strconv.Atoi(c.Param("eventId"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Cannot find param [EventId]"})
	}
	params = append(params, eventId)
	keys = append(keys, "event_id = ?")

	// query builder
	var result models.PerlEventExportSetting
	query := e.db.QueryContext(models.PerlEventExportSetting{}, c)
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
		event := fmt.Sprintf("Updated [PerlEventExportSetting] [%v] fields [%v]", strings.Join(ids, ", "), strings.Join(fieldsUpdated, ", "))
		e.auditLog.LogUserEvent(c, "UPDATE", event)
	}

	return c.JSON(http.StatusOK, request)
}

// createPerlEventExportSetting godoc
// @Id createPerlEventExportSetting
// @Summary Creates PerlEventExportSetting
// @Accept json
// @Produce json
// @Param perl_event_export_setting body models.PerlEventExportSetting true "PerlEventExportSetting"
// @Tags PerlEventExportSetting
// @Success 200 {array} models.PerlEventExportSetting
// @Failure 500 {string} string "Error binding to entity"
// @Failure 500 {string} string "Error inserting entity"
// @Router /perl_event_export_setting [put]
func (e *PerlEventExportSettingController) createPerlEventExportSetting(c echo.Context) error {
	perlEventExportSetting := new(models.PerlEventExportSetting)
	if err := c.Bind(perlEventExportSetting); err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			echo.Map{"error": fmt.Sprintf("Error binding to entity [%v]", err.Error())},
		)
	}

	err := e.db.Get(models.PerlEventExportSetting{}, c).Model(&models.PerlEventExportSetting{}).Create(&perlEventExportSetting).Error
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			echo.Map{"error": fmt.Sprintf("Error inserting entity [%v]", err.Error())},
		)
	}

	// log create event
	if e.db.GetSpireDb() != nil {
		// diff between an empty model and the created
		diff := database.ResultDifference(models.PerlEventExportSetting{}, perlEventExportSetting)
		// build fields created
		var fields []string
		for k, v := range diff {
			fields = append(fields, fmt.Sprintf("%v = %v", k, v))
		}
		// record event
		event := fmt.Sprintf("Created [PerlEventExportSetting] [%v] data [%v]", perlEventExportSetting.EventId, strings.Join(fields, ", "))
		e.auditLog.LogUserEvent(c, "CREATE", event)
	}

	return c.JSON(http.StatusOK, perlEventExportSetting)
}

// deletePerlEventExportSetting godoc
// @Id deletePerlEventExportSetting
// @Summary Deletes PerlEventExportSetting
// @Accept json
// @Produce json
// @Tags PerlEventExportSetting
// @Param id path int true "eventId"
// @Success 200 {string} string "Entity deleted successfully"
// @Failure 404 {string} string "Cannot find entity"
// @Failure 500 {string} string "Error binding to entity"
// @Failure 500 {string} string "Error deleting entity"
// @Router /perl_event_export_setting/{id} [delete]
func (e *PerlEventExportSettingController) deletePerlEventExportSetting(c echo.Context) error {
	var params []interface{}
	var keys []string

	// primary key param
	eventId, err := strconv.Atoi(c.Param("eventId"))
	if err != nil {
		e.logger.Error(err)
	}
	params = append(params, eventId)
	keys = append(keys, "event_id = ?")

	// query builder
	var result models.PerlEventExportSetting
	query := e.db.QueryContext(models.PerlEventExportSetting{}, c)
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
		event := fmt.Sprintf("Deleted [PerlEventExportSetting] [%v] keys [%v]", result.EventId, strings.Join(ids, ", "))
		e.auditLog.LogUserEvent(c, "DELETE", event)
	}

	return c.JSON(http.StatusOK, echo.Map{"success": "Entity deleted successfully"})
}

// getPerlEventExportSettingsBulk godoc
// @Id getPerlEventExportSettingsBulk
// @Summary Gets PerlEventExportSettings in bulk
// @Accept json
// @Produce json
// @Param Body body BulkFetchByIdsGetRequest true "body"
// @Tags PerlEventExportSetting
// @Success 200 {array} models.PerlEventExportSetting
// @Failure 500 {string} string "Bad query request"
// @Router /perl_event_export_settings/bulk [post]
func (e *PerlEventExportSettingController) getPerlEventExportSettingsBulk(c echo.Context) error {
	var results []models.PerlEventExportSetting

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

	err := e.db.QueryContext(models.PerlEventExportSetting{}, c).Find(&results, r.IDs).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, results)
}
