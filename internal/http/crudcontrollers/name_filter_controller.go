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

type NameFilterController struct {
	db       *database.DatabaseResolver
	logger   *logrus.Logger
	auditLog *auditlog.UserEvent
}

func NewNameFilterController(
	db *database.DatabaseResolver,
	logger *logrus.Logger,
	auditLog *auditlog.UserEvent,
) *NameFilterController {
	return &NameFilterController{
		db:       db,
		logger:   logger,
		auditLog: auditLog,
	}
}

func (e *NameFilterController) Routes() []*routes.Route {
	return []*routes.Route{
		routes.RegisterRoute(http.MethodGet, "name_filter/:id", e.getNameFilter, nil),
		routes.RegisterRoute(http.MethodGet, "name_filters", e.listNameFilters, nil),
		routes.RegisterRoute(http.MethodGet, "name_filters/count", e.getNameFiltersCount, nil),
		routes.RegisterRoute(http.MethodPut, "name_filter", e.createNameFilter, nil),
		routes.RegisterRoute(http.MethodDelete, "name_filter/:id", e.deleteNameFilter, nil),
		routes.RegisterRoute(http.MethodPatch, "name_filter/:id", e.updateNameFilter, nil),
		routes.RegisterRoute(http.MethodPost, "name_filters/bulk", e.getNameFiltersBulk, nil),
	}
}

// listNameFilters godoc
// @Id listNameFilters
// @Summary Lists NameFilters
// @Accept json
// @Produce json
// @Tags NameFilter
// @Param includes query string false "Relationships [all] for all [number] for depth of relationships to load or [.] separated relationship names"
// @Param where query string false "Filter on specific fields. Multiple conditions [.] separated Example: col_like_value.col2__val2"
// @Param whereOr query string false "Filter on specific fields (Chained ors). Multiple conditions [.] separated Example: col_like_value.col2__val2"
// @Param groupBy query string false "Group by field. Multiple conditions [.] separated Example: field1.field2"
// @Param limit query string false "Rows to limit in response (Default: 10,000)"
// @Param page query int 0 "Pagination page"
// @Param orderBy query string false "Order by [field]"
// @Param orderDirection query string false "Order by field direction"
// @Param select query string false "Column names [.] separated to fetch specific fields in response"
// @Success 200 {array} models.NameFilter
// @Failure 500 {string} string "Bad query request"
// @Router /name_filters [get]
func (e *NameFilterController) listNameFilters(c echo.Context) error {
	var results []models.NameFilter
	err := e.db.QueryContext(models.NameFilter{}, c).Find(&results).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err})
	}

	return c.JSON(http.StatusOK, results)
}

// getNameFilter godoc
// @Id getNameFilter
// @Summary Gets NameFilter
// @Accept json
// @Produce json
// @Tags NameFilter
// @Param id path int true "Id"
// @Param includes query string false "Relationships [all] for all [number] for depth of relationships to load or [.] separated relationship names"
// @Param select query string false "Column names [.] separated to fetch specific fields in response"
// @Success 200 {array} models.NameFilter
// @Failure 404 {string} string "Entity not found"
// @Failure 500 {string} string "Cannot find param"
// @Failure 500 {string} string "Bad query request"
// @Router /name_filter/{id} [get]
func (e *NameFilterController) getNameFilter(c echo.Context) error {
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
	var result models.NameFilter
	query := e.db.QueryContext(models.NameFilter{}, c)
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

// updateNameFilter godoc
// @Id updateNameFilter
// @Summary Updates NameFilter
// @Accept json
// @Produce json
// @Tags NameFilter
// @Param id path int true "Id"
// @Param name_filter body models.NameFilter true "NameFilter"
// @Success 200 {array} models.NameFilter
// @Failure 404 {string} string "Cannot find entity"
// @Failure 500 {string} string "Error binding to entity"
// @Failure 500 {string} string "Error updating entity"
// @Router /name_filter/{id} [patch]
func (e *NameFilterController) updateNameFilter(c echo.Context) error {
	request := new(models.NameFilter)
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
	var result models.NameFilter
	query := e.db.QueryContext(models.NameFilter{}, c)
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
		event := fmt.Sprintf("Updated [NameFilter] [%v] fields [%v]", strings.Join(ids, ", "), strings.Join(fieldsUpdated, ", "))
		e.auditLog.LogUserEvent(c, "UPDATE", event)
	}

	return c.JSON(http.StatusOK, request)
}

// createNameFilter godoc
// @Id createNameFilter
// @Summary Creates NameFilter
// @Accept json
// @Produce json
// @Param name_filter body models.NameFilter true "NameFilter"
// @Tags NameFilter
// @Success 200 {array} models.NameFilter
// @Failure 500 {string} string "Error binding to entity"
// @Failure 500 {string} string "Error inserting entity"
// @Router /name_filter [put]
func (e *NameFilterController) createNameFilter(c echo.Context) error {
	nameFilter := new(models.NameFilter)
	if err := c.Bind(nameFilter); err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			echo.Map{"error": fmt.Sprintf("Error binding to entity [%v]", err.Error())},
		)
	}

	db := e.db.Get(models.NameFilter{}, c).Model(&models.NameFilter{})

	// save associations
	if c.QueryParam("save_associations") != "true" {
		db = db.Omit(clause.Associations)
	}

	err := db.Create(&nameFilter).Error
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			echo.Map{"error": fmt.Sprintf("Error inserting entity [%v]", err.Error())},
		)
	}

	// log create event
	if e.db.GetSpireDb() != nil {
		// diff between an empty model and the created
		diff := database.ResultDifference(models.NameFilter{}, nameFilter)
		// build fields created
		var fields []string
		for k, v := range diff {
			fields = append(fields, fmt.Sprintf("%v = %v", k, v))
		}
		// record event
		event := fmt.Sprintf("Created [NameFilter] [%v] data [%v]", nameFilter.ID, strings.Join(fields, ", "))
		e.auditLog.LogUserEvent(c, "CREATE", event)
	}

	return c.JSON(http.StatusOK, nameFilter)
}

// deleteNameFilter godoc
// @Id deleteNameFilter
// @Summary Deletes NameFilter
// @Accept json
// @Produce json
// @Tags NameFilter
// @Param id path int true "id"
// @Success 200 {string} string "Entity deleted successfully"
// @Failure 404 {string} string "Cannot find entity"
// @Failure 500 {string} string "Error binding to entity"
// @Failure 500 {string} string "Error deleting entity"
// @Router /name_filter/{id} [delete]
func (e *NameFilterController) deleteNameFilter(c echo.Context) error {
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
	var result models.NameFilter
	query := e.db.QueryContext(models.NameFilter{}, c)
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
		event := fmt.Sprintf("Deleted [NameFilter] [%v] keys [%v]", result.ID, strings.Join(ids, ", "))
		e.auditLog.LogUserEvent(c, "DELETE", event)
	}

	return c.JSON(http.StatusOK, echo.Map{"success": "Entity deleted successfully"})
}

// getNameFiltersBulk godoc
// @Id getNameFiltersBulk
// @Summary Gets NameFilters in bulk
// @Accept json
// @Produce json
// @Param Body body BulkFetchByIdsGetRequest true "body"
// @Tags NameFilter
// @Success 200 {array} models.NameFilter
// @Failure 500 {string} string "Bad query request"
// @Router /name_filters/bulk [post]
func (e *NameFilterController) getNameFiltersBulk(c echo.Context) error {
	var results []models.NameFilter

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

	err := e.db.QueryContext(models.NameFilter{}, c).Find(&results, r.IDs).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, results)
}

// getNameFiltersCount godoc
// @Id getNameFiltersCount
// @Summary Counts NameFilters
// @Accept json
// @Produce json
// @Tags NameFilter
// @Param includes query string false "Relationships [all] for all [number] for depth of relationships to load or [.] separated relationship names"
// @Param where query string false "Filter on specific fields. Multiple conditions [.] separated Example: col_like_value.col2__val2"
// @Param whereOr query string false "Filter on specific fields (Chained ors). Multiple conditions [.] separated Example: col_like_value.col2__val2"
// @Param groupBy query string false "Group by field. Multiple conditions [.] separated Example: field1.field2"
// @Param limit query string false "Rows to limit in response (Default: 10,000)"
// @Param page query int 0 "Pagination page"
// @Param orderBy query string false "Order by [field]"
// @Param orderDirection query string false "Order by field direction"
// @Param select query string false "Column names [.] separated to fetch specific fields in response"
// @Success 200 {array} models.NameFilter
// @Failure 500 {string} string "Bad query request"
// @Router /name_filters/count [get]
func (e *NameFilterController) getNameFiltersCount(c echo.Context) error {
	var count int64
	err := e.db.QueryContext(models.NameFilter{}, c).Count(&count).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err})
	}

	return c.JSON(http.StatusOK, echo.Map{"count": count})
}