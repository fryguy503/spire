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
	"gorm.io/gorm/clause"
	"net/http"
	"strconv"
	"strings"
)

type FactionBaseDatumController struct {
	db       *database.DatabaseResolver
	logger   *logrus.Logger
	auditLog *auditlog.UserEvent
}

func NewFactionBaseDatumController(
	db *database.DatabaseResolver,
	logger *logrus.Logger,
	auditLog *auditlog.UserEvent,
) *FactionBaseDatumController {
	return &FactionBaseDatumController{
		db:       db,
		logger:   logger,
		auditLog: auditLog,
	}
}

func (e *FactionBaseDatumController) Routes() []*routes.Route {
	return []*routes.Route{
		routes.RegisterRoute(http.MethodGet, "faction_base_datum/:clientFactionId", e.getFactionBaseDatum, nil),
		routes.RegisterRoute(http.MethodGet, "faction_base_data", e.listFactionBaseData, nil),
		routes.RegisterRoute(http.MethodGet, "faction_base_data/count", e.getFactionBaseDataCount, nil),
		routes.RegisterRoute(http.MethodPut, "faction_base_datum", e.createFactionBaseDatum, nil),
		routes.RegisterRoute(http.MethodDelete, "faction_base_datum/:clientFactionId", e.deleteFactionBaseDatum, nil),
		routes.RegisterRoute(http.MethodPatch, "faction_base_datum/:clientFactionId", e.updateFactionBaseDatum, nil),
		routes.RegisterRoute(http.MethodPost, "faction_base_data/bulk", e.getFactionBaseDataBulk, nil),
	}
}

// listFactionBaseData godoc
// @Id listFactionBaseData
// @Summary Lists FactionBaseData
// @Accept json
// @Produce json
// @Tags FactionBaseDatum
// @Param includes query string false "Relationships [all] for all [number] for depth of relationships to load or [.] separated relationship names"
// @Param where query string false "Filter on specific fields. Multiple conditions [.] separated Example: col_like_value.col2__val2"
// @Param whereOr query string false "Filter on specific fields (Chained ors). Multiple conditions [.] separated Example: col_like_value.col2__val2"
// @Param groupBy query string false "Group by field. Multiple conditions [.] separated Example: field1.field2"
// @Param limit query string false "Rows to limit in response (Default: 10,000)"
// @Param page query int 0 "Pagination page"
// @Param orderBy query string false "Order by [field]"
// @Param orderDirection query string false "Order by field direction"
// @Param select query string false "Column names [.] separated to fetch specific fields in response"
// @Success 200 {array} models.FactionBaseDatum
// @Failure 500 {string} string "Bad query request"
// @Router /faction_base_data [get]
func (e *FactionBaseDatumController) listFactionBaseData(c echo.Context) error {
	var results []models.FactionBaseDatum
	err := e.db.QueryContext(models.FactionBaseDatum{}, c).Find(&results).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err})
	}

	return c.JSON(http.StatusOK, results)
}

// getFactionBaseDatum godoc
// @Id getFactionBaseDatum
// @Summary Gets FactionBaseDatum
// @Accept json
// @Produce json
// @Tags FactionBaseDatum
// @Param id path int true "Id"
// @Param includes query string false "Relationships [all] for all [number] for depth of relationships to load or [.] separated relationship names"
// @Param select query string false "Column names [.] separated to fetch specific fields in response"
// @Success 200 {array} models.FactionBaseDatum
// @Failure 404 {string} string "Entity not found"
// @Failure 500 {string} string "Cannot find param"
// @Failure 500 {string} string "Bad query request"
// @Router /faction_base_datum/{id} [get]
func (e *FactionBaseDatumController) getFactionBaseDatum(c echo.Context) error {
	var params []interface{}
	var keys []string

	// primary key param
	clientFactionId, err := strconv.Atoi(c.Param("clientFactionId"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Cannot find param [ClientFactionId]"})
	}
	params = append(params, clientFactionId)
	keys = append(keys, "client_faction_id = ?")

	// query builder
	var result models.FactionBaseDatum
	query := e.db.QueryContext(models.FactionBaseDatum{}, c)
	for i, _ := range keys {
		query = query.Where(keys[i], params[i])
	}

	// grab first entry
	err = query.First(&result).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	// couldn't find entity
	if result.ClientFactionId == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Cannot find entity"})
	}

	return c.JSON(http.StatusOK, result)
}

// updateFactionBaseDatum godoc
// @Id updateFactionBaseDatum
// @Summary Updates FactionBaseDatum
// @Accept json
// @Produce json
// @Tags FactionBaseDatum
// @Param id path int true "Id"
// @Param faction_base_datum body models.FactionBaseDatum true "FactionBaseDatum"
// @Success 200 {array} models.FactionBaseDatum
// @Failure 404 {string} string "Cannot find entity"
// @Failure 500 {string} string "Error binding to entity"
// @Failure 500 {string} string "Error updating entity"
// @Router /faction_base_datum/{id} [patch]
func (e *FactionBaseDatumController) updateFactionBaseDatum(c echo.Context) error {
	request := new(models.FactionBaseDatum)
	if err := c.Bind(request); err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			echo.Map{"error": fmt.Sprintf("Error binding to entity [%v]", err.Error())},
		)
	}

	var params []interface{}
	var keys []string

	// primary key param
	clientFactionId, err := strconv.Atoi(c.Param("clientFactionId"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Cannot find param [ClientFactionId]"})
	}
	params = append(params, clientFactionId)
	keys = append(keys, "client_faction_id = ?")

	// query builder
	var result models.FactionBaseDatum
	query := e.db.QueryContext(models.FactionBaseDatum{}, c)
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
		event := fmt.Sprintf("Updated [FactionBaseDatum] [%v] fields [%v]", strings.Join(ids, ", "), strings.Join(fieldsUpdated, ", "))
		e.auditLog.LogUserEvent(c, "UPDATE", event)
	}

	return c.JSON(http.StatusOK, request)
}

// createFactionBaseDatum godoc
// @Id createFactionBaseDatum
// @Summary Creates FactionBaseDatum
// @Accept json
// @Produce json
// @Param faction_base_datum body models.FactionBaseDatum true "FactionBaseDatum"
// @Tags FactionBaseDatum
// @Success 200 {array} models.FactionBaseDatum
// @Failure 500 {string} string "Error binding to entity"
// @Failure 500 {string} string "Error inserting entity"
// @Router /faction_base_datum [put]
func (e *FactionBaseDatumController) createFactionBaseDatum(c echo.Context) error {
	factionBaseDatum := new(models.FactionBaseDatum)
	if err := c.Bind(factionBaseDatum); err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			echo.Map{"error": fmt.Sprintf("Error binding to entity [%v]", err.Error())},
		)
	}

	db := e.db.Get(models.FactionBaseDatum{}, c).Model(&models.FactionBaseDatum{})

	// save associations
	if c.QueryParam("save_associations") != "true" {
		db = db.Omit(clause.Associations)
	}

	err := db.Create(&factionBaseDatum).Error
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			echo.Map{"error": fmt.Sprintf("Error inserting entity [%v]", err.Error())},
		)
	}

	// log create event
	if e.db.GetSpireDb() != nil {
		// diff between an empty model and the created
		diff := database.ResultDifference(models.FactionBaseDatum{}, factionBaseDatum)
		// build fields created
		var fields []string
		for k, v := range diff {
			fields = append(fields, fmt.Sprintf("%v = %v", k, v))
		}
		// record event
		event := fmt.Sprintf("Created [FactionBaseDatum] [%v] data [%v]", factionBaseDatum.ClientFactionId, strings.Join(fields, ", "))
		e.auditLog.LogUserEvent(c, "CREATE", event)
	}

	return c.JSON(http.StatusOK, factionBaseDatum)
}

// deleteFactionBaseDatum godoc
// @Id deleteFactionBaseDatum
// @Summary Deletes FactionBaseDatum
// @Accept json
// @Produce json
// @Tags FactionBaseDatum
// @Param id path int true "clientFactionId"
// @Success 200 {string} string "Entity deleted successfully"
// @Failure 404 {string} string "Cannot find entity"
// @Failure 500 {string} string "Error binding to entity"
// @Failure 500 {string} string "Error deleting entity"
// @Router /faction_base_datum/{id} [delete]
func (e *FactionBaseDatumController) deleteFactionBaseDatum(c echo.Context) error {
	var params []interface{}
	var keys []string

	// primary key param
	clientFactionId, err := strconv.Atoi(c.Param("clientFactionId"))
	if err != nil {
		e.logger.Error(err)
	}
	params = append(params, clientFactionId)
	keys = append(keys, "client_faction_id = ?")

	// query builder
	var result models.FactionBaseDatum
	query := e.db.QueryContext(models.FactionBaseDatum{}, c)
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
		event := fmt.Sprintf("Deleted [FactionBaseDatum] [%v] keys [%v]", result.ClientFactionId, strings.Join(ids, ", "))
		e.auditLog.LogUserEvent(c, "DELETE", event)
	}

	return c.JSON(http.StatusOK, echo.Map{"success": "Entity deleted successfully"})
}

// getFactionBaseDataBulk godoc
// @Id getFactionBaseDataBulk
// @Summary Gets FactionBaseData in bulk
// @Accept json
// @Produce json
// @Param Body body BulkFetchByIdsGetRequest true "body"
// @Tags FactionBaseDatum
// @Success 200 {array} models.FactionBaseDatum
// @Failure 500 {string} string "Bad query request"
// @Router /faction_base_data/bulk [post]
func (e *FactionBaseDatumController) getFactionBaseDataBulk(c echo.Context) error {
	var results []models.FactionBaseDatum

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

	err := e.db.QueryContext(models.FactionBaseDatum{}, c).Find(&results, r.IDs).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, results)
}

// getFactionBaseDataCount godoc
// @Id getFactionBaseDataCount
// @Summary Counts FactionBaseData
// @Accept json
// @Produce json
// @Tags FactionBaseDatum
// @Param includes query string false "Relationships [all] for all [number] for depth of relationships to load or [.] separated relationship names"
// @Param where query string false "Filter on specific fields. Multiple conditions [.] separated Example: col_like_value.col2__val2"
// @Param whereOr query string false "Filter on specific fields (Chained ors). Multiple conditions [.] separated Example: col_like_value.col2__val2"
// @Param groupBy query string false "Group by field. Multiple conditions [.] separated Example: field1.field2"
// @Param limit query string false "Rows to limit in response (Default: 10,000)"
// @Param page query int 0 "Pagination page"
// @Param orderBy query string false "Order by [field]"
// @Param orderDirection query string false "Order by field direction"
// @Param select query string false "Column names [.] separated to fetch specific fields in response"
// @Success 200 {array} models.FactionBaseDatum
// @Failure 500 {string} string "Bad query request"
// @Router /faction_base_data/count [get]
func (e *FactionBaseDatumController) getFactionBaseDataCount(c echo.Context) error {
	var count int64
	err := e.db.QueryContext(models.FactionBaseDatum{}, c).Count(&count).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err})
	}

	return c.JSON(http.StatusOK, echo.Map{"count": count})
}