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

type LfguildController struct {
	db       *database.DatabaseResolver
	logger   *logrus.Logger
	auditLog *auditlog.UserEvent
}

func NewLfguildController(
	db *database.DatabaseResolver,
	logger *logrus.Logger,
	auditLog *auditlog.UserEvent,
) *LfguildController {
	return &LfguildController{
		db:       db,
		logger:   logger,
		auditLog: auditLog,
	}
}

func (e *LfguildController) Routes() []*routes.Route {
	return []*routes.Route{
		routes.RegisterRoute(http.MethodGet, "lfguild/:typeId", e.getLfguild, nil),
		routes.RegisterRoute(http.MethodGet, "lfguilds", e.listLfguilds, nil),
		routes.RegisterRoute(http.MethodGet, "lfguilds/count", e.getLfguildsCount, nil),
		routes.RegisterRoute(http.MethodPut, "lfguild", e.createLfguild, nil),
		routes.RegisterRoute(http.MethodDelete, "lfguild/:typeId", e.deleteLfguild, nil),
		routes.RegisterRoute(http.MethodPatch, "lfguild/:typeId", e.updateLfguild, nil),
		routes.RegisterRoute(http.MethodPost, "lfguilds/bulk", e.getLfguildsBulk, nil),
	}
}

// listLfguilds godoc
// @Id listLfguilds
// @Summary Lists Lfguilds
// @Accept json
// @Produce json
// @Tags Lfguild
// @Param includes query string false "Relationships [all] for all [number] for depth of relationships to load or [.] separated relationship names"
// @Param where query string false "Filter on specific fields. Multiple conditions [.] separated Example: col_like_value.col2__val2"
// @Param whereOr query string false "Filter on specific fields (Chained ors). Multiple conditions [.] separated Example: col_like_value.col2__val2"
// @Param groupBy query string false "Group by field. Multiple conditions [.] separated Example: field1.field2"
// @Param limit query string false "Rows to limit in response (Default: 10,000)"
// @Param page query int 0 "Pagination page"
// @Param orderBy query string false "Order by [field]"
// @Param orderDirection query string false "Order by field direction"
// @Param select query string false "Column names [.] separated to fetch specific fields in response"
// @Success 200 {array} models.Lfguild
// @Failure 500 {string} string "Bad query request"
// @Router /lfguilds [get]
func (e *LfguildController) listLfguilds(c echo.Context) error {
	var results []models.Lfguild
	err := e.db.QueryContext(models.Lfguild{}, c).Find(&results).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err})
	}

	return c.JSON(http.StatusOK, results)
}

// getLfguild godoc
// @Id getLfguild
// @Summary Gets Lfguild
// @Accept json
// @Produce json
// @Tags Lfguild
// @Param id path int true "Id"
// @Param includes query string false "Relationships [all] for all [number] for depth of relationships to load or [.] separated relationship names"
// @Param select query string false "Column names [.] separated to fetch specific fields in response"
// @Success 200 {array} models.Lfguild
// @Failure 404 {string} string "Entity not found"
// @Failure 500 {string} string "Cannot find param"
// @Failure 500 {string} string "Bad query request"
// @Router /lfguild/{id} [get]
func (e *LfguildController) getLfguild(c echo.Context) error {
	var params []interface{}
	var keys []string

	// primary key param
	typeId, err := strconv.Atoi(c.Param("typeId"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Cannot find param [TypeId]"})
	}
	params = append(params, typeId)
	keys = append(keys, "type = ?")

	// key param [name] position [2] type [varchar]
	if len(c.QueryParam("name")) > 0 {
		nameParam, err := strconv.Atoi(c.QueryParam("name"))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("Error parsing query param [name] err [%s]", err.Error())})
		}

		params = append(params, nameParam)
		keys = append(keys, "name = ?")
	}

	// query builder
	var result models.Lfguild
	query := e.db.QueryContext(models.Lfguild{}, c)
	for i, _ := range keys {
		query = query.Where(keys[i], params[i])
	}

	// grab first entry
	err = query.First(&result).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	// couldn't find entity
	if result.Type == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Cannot find entity"})
	}

	return c.JSON(http.StatusOK, result)
}

// updateLfguild godoc
// @Id updateLfguild
// @Summary Updates Lfguild
// @Accept json
// @Produce json
// @Tags Lfguild
// @Param id path int true "Id"
// @Param lfguild body models.Lfguild true "Lfguild"
// @Success 200 {array} models.Lfguild
// @Failure 404 {string} string "Cannot find entity"
// @Failure 500 {string} string "Error binding to entity"
// @Failure 500 {string} string "Error updating entity"
// @Router /lfguild/{id} [patch]
func (e *LfguildController) updateLfguild(c echo.Context) error {
	request := new(models.Lfguild)
	if err := c.Bind(request); err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			echo.Map{"error": fmt.Sprintf("Error binding to entity [%v]", err.Error())},
		)
	}

	var params []interface{}
	var keys []string

	// primary key param
	typeId, err := strconv.Atoi(c.Param("typeId"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Cannot find param [TypeId]"})
	}
	params = append(params, typeId)
	keys = append(keys, "type = ?")

	// key param [name] position [2] type [varchar]
	if len(c.QueryParam("name")) > 0 {
		nameParam, err := strconv.Atoi(c.QueryParam("name"))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("Error parsing query param [name] err [%s]", err.Error())})
		}

		params = append(params, nameParam)
		keys = append(keys, "name = ?")
	}

	// query builder
	var result models.Lfguild
	query := e.db.QueryContext(models.Lfguild{}, c)
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
		event := fmt.Sprintf("Updated [Lfguild] [%v] fields [%v]", strings.Join(ids, ", "), strings.Join(fieldsUpdated, ", "))
		e.auditLog.LogUserEvent(c, "UPDATE", event)
	}

	return c.JSON(http.StatusOK, request)
}

// createLfguild godoc
// @Id createLfguild
// @Summary Creates Lfguild
// @Accept json
// @Produce json
// @Param lfguild body models.Lfguild true "Lfguild"
// @Tags Lfguild
// @Success 200 {array} models.Lfguild
// @Failure 500 {string} string "Error binding to entity"
// @Failure 500 {string} string "Error inserting entity"
// @Router /lfguild [put]
func (e *LfguildController) createLfguild(c echo.Context) error {
	lfguild := new(models.Lfguild)
	if err := c.Bind(lfguild); err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			echo.Map{"error": fmt.Sprintf("Error binding to entity [%v]", err.Error())},
		)
	}

	db := e.db.Get(models.Lfguild{}, c).Model(&models.Lfguild{})

	// save associations
	if c.QueryParam("save_associations") != "true" {
		db = db.Omit(clause.Associations)
	}

	err := db.Create(&lfguild).Error
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			echo.Map{"error": fmt.Sprintf("Error inserting entity [%v]", err.Error())},
		)
	}

	// log create event
	if e.db.GetSpireDb() != nil {
		// diff between an empty model and the created
		diff := database.ResultDifference(models.Lfguild{}, lfguild)
		// build fields created
		var fields []string
		for k, v := range diff {
			fields = append(fields, fmt.Sprintf("%v = %v", k, v))
		}
		// record event
		event := fmt.Sprintf("Created [Lfguild] [%v] data [%v]", lfguild.Type, strings.Join(fields, ", "))
		e.auditLog.LogUserEvent(c, "CREATE", event)
	}

	return c.JSON(http.StatusOK, lfguild)
}

// deleteLfguild godoc
// @Id deleteLfguild
// @Summary Deletes Lfguild
// @Accept json
// @Produce json
// @Tags Lfguild
// @Param id path int true "typeId"
// @Success 200 {string} string "Entity deleted successfully"
// @Failure 404 {string} string "Cannot find entity"
// @Failure 500 {string} string "Error binding to entity"
// @Failure 500 {string} string "Error deleting entity"
// @Router /lfguild/{id} [delete]
func (e *LfguildController) deleteLfguild(c echo.Context) error {
	var params []interface{}
	var keys []string

	// primary key param
	typeId, err := strconv.Atoi(c.Param("typeId"))
	if err != nil {
		e.logger.Error(err)
	}
	params = append(params, typeId)
	keys = append(keys, "type = ?")

	// key param [name] position [2] type [varchar]
	if len(c.QueryParam("name")) > 0 {
		nameParam, err := strconv.Atoi(c.QueryParam("name"))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("Error parsing query param [name] err [%s]", err.Error())})
		}

		params = append(params, nameParam)
		keys = append(keys, "name = ?")
	}

	// query builder
	var result models.Lfguild
	query := e.db.QueryContext(models.Lfguild{}, c)
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
		event := fmt.Sprintf("Deleted [Lfguild] [%v] keys [%v]", result.Type, strings.Join(ids, ", "))
		e.auditLog.LogUserEvent(c, "DELETE", event)
	}

	return c.JSON(http.StatusOK, echo.Map{"success": "Entity deleted successfully"})
}

// getLfguildsBulk godoc
// @Id getLfguildsBulk
// @Summary Gets Lfguilds in bulk
// @Accept json
// @Produce json
// @Param Body body BulkFetchByIdsGetRequest true "body"
// @Tags Lfguild
// @Success 200 {array} models.Lfguild
// @Failure 500 {string} string "Bad query request"
// @Router /lfguilds/bulk [post]
func (e *LfguildController) getLfguildsBulk(c echo.Context) error {
	var results []models.Lfguild

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

	err := e.db.QueryContext(models.Lfguild{}, c).Find(&results, r.IDs).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, results)
}

// getLfguildsCount godoc
// @Id getLfguildsCount
// @Summary Counts Lfguilds
// @Accept json
// @Produce json
// @Tags Lfguild
// @Param includes query string false "Relationships [all] for all [number] for depth of relationships to load or [.] separated relationship names"
// @Param where query string false "Filter on specific fields. Multiple conditions [.] separated Example: col_like_value.col2__val2"
// @Param whereOr query string false "Filter on specific fields (Chained ors). Multiple conditions [.] separated Example: col_like_value.col2__val2"
// @Param groupBy query string false "Group by field. Multiple conditions [.] separated Example: field1.field2"
// @Param limit query string false "Rows to limit in response (Default: 10,000)"
// @Param page query int 0 "Pagination page"
// @Param orderBy query string false "Order by [field]"
// @Param orderDirection query string false "Order by field direction"
// @Param select query string false "Column names [.] separated to fetch specific fields in response"
// @Success 200 {array} models.Lfguild
// @Failure 500 {string} string "Bad query request"
// @Router /lfguilds/count [get]
func (e *LfguildController) getLfguildsCount(c echo.Context) error {
	var count int64
	err := e.db.QueryContext(models.Lfguild{}, c).Count(&count).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err})
	}

	return c.JSON(http.StatusOK, echo.Map{"count": count})
}