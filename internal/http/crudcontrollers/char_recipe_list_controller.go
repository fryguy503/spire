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

type CharRecipeListController struct {
	db       *database.DatabaseResolver
	logger   *logrus.Logger
	auditLog *auditlog.UserEvent
}

func NewCharRecipeListController(
	db *database.DatabaseResolver,
	logger *logrus.Logger,
	auditLog *auditlog.UserEvent,
) *CharRecipeListController {
	return &CharRecipeListController{
		db:       db,
		logger:   logger,
		auditLog: auditLog,
	}
}

func (e *CharRecipeListController) Routes() []*routes.Route {
	return []*routes.Route{
		routes.RegisterRoute(http.MethodGet, "char_recipe_list/:charId", e.getCharRecipeList, nil),
		routes.RegisterRoute(http.MethodGet, "char_recipe_lists", e.listCharRecipeLists, nil),
		routes.RegisterRoute(http.MethodGet, "char_recipe_lists/count", e.getCharRecipeListsCount, nil),
		routes.RegisterRoute(http.MethodPut, "char_recipe_list", e.createCharRecipeList, nil),
		routes.RegisterRoute(http.MethodDelete, "char_recipe_list/:charId", e.deleteCharRecipeList, nil),
		routes.RegisterRoute(http.MethodPatch, "char_recipe_list/:charId", e.updateCharRecipeList, nil),
		routes.RegisterRoute(http.MethodPost, "char_recipe_lists/bulk", e.getCharRecipeListsBulk, nil),
	}
}

// listCharRecipeLists godoc
// @Id listCharRecipeLists
// @Summary Lists CharRecipeLists
// @Accept json
// @Produce json
// @Tags CharRecipeList
// @Param includes query string false "Relationships [all] for all [number] for depth of relationships to load or [.] separated relationship names"
// @Param where query string false "Filter on specific fields. Multiple conditions [.] separated Example: col_like_value.col2__val2"
// @Param whereOr query string false "Filter on specific fields (Chained ors). Multiple conditions [.] separated Example: col_like_value.col2__val2"
// @Param groupBy query string false "Group by field. Multiple conditions [.] separated Example: field1.field2"
// @Param limit query string false "Rows to limit in response (Default: 10,000)"
// @Param page query int 0 "Pagination page"
// @Param orderBy query string false "Order by [field]"
// @Param orderDirection query string false "Order by field direction"
// @Param select query string false "Column names [.] separated to fetch specific fields in response"
// @Success 200 {array} models.CharRecipeList
// @Failure 500 {string} string "Bad query request"
// @Router /char_recipe_lists [get]
func (e *CharRecipeListController) listCharRecipeLists(c echo.Context) error {
	var results []models.CharRecipeList
	err := e.db.QueryContext(models.CharRecipeList{}, c).Find(&results).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err})
	}

	return c.JSON(http.StatusOK, results)
}

// getCharRecipeList godoc
// @Id getCharRecipeList
// @Summary Gets CharRecipeList
// @Accept json
// @Produce json
// @Tags CharRecipeList
// @Param id path int true "Id"
// @Param includes query string false "Relationships [all] for all [number] for depth of relationships to load or [.] separated relationship names"
// @Param select query string false "Column names [.] separated to fetch specific fields in response"
// @Success 200 {array} models.CharRecipeList
// @Failure 404 {string} string "Entity not found"
// @Failure 500 {string} string "Cannot find param"
// @Failure 500 {string} string "Bad query request"
// @Router /char_recipe_list/{id} [get]
func (e *CharRecipeListController) getCharRecipeList(c echo.Context) error {
	var params []interface{}
	var keys []string

	// primary key param
	charId, err := strconv.Atoi(c.Param("charId"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Cannot find param [CharId]"})
	}
	params = append(params, charId)
	keys = append(keys, "char_id = ?")

	// key param [recipe_id] position [2] type [int]
	if len(c.QueryParam("recipe_id")) > 0 {
		recipeIdParam, err := strconv.Atoi(c.QueryParam("recipe_id"))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("Error parsing query param [recipe_id] err [%s]", err.Error())})
		}

		params = append(params, recipeIdParam)
		keys = append(keys, "recipe_id = ?")
	}

	// query builder
	var result models.CharRecipeList
	query := e.db.QueryContext(models.CharRecipeList{}, c)
	for i, _ := range keys {
		query = query.Where(keys[i], params[i])
	}

	// grab first entry
	err = query.First(&result).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	// couldn't find entity
	if result.CharId == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Cannot find entity"})
	}

	return c.JSON(http.StatusOK, result)
}

// updateCharRecipeList godoc
// @Id updateCharRecipeList
// @Summary Updates CharRecipeList
// @Accept json
// @Produce json
// @Tags CharRecipeList
// @Param id path int true "Id"
// @Param char_recipe_list body models.CharRecipeList true "CharRecipeList"
// @Success 200 {array} models.CharRecipeList
// @Failure 404 {string} string "Cannot find entity"
// @Failure 500 {string} string "Error binding to entity"
// @Failure 500 {string} string "Error updating entity"
// @Router /char_recipe_list/{id} [patch]
func (e *CharRecipeListController) updateCharRecipeList(c echo.Context) error {
	request := new(models.CharRecipeList)
	if err := c.Bind(request); err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			echo.Map{"error": fmt.Sprintf("Error binding to entity [%v]", err.Error())},
		)
	}

	var params []interface{}
	var keys []string

	// primary key param
	charId, err := strconv.Atoi(c.Param("charId"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Cannot find param [CharId]"})
	}
	params = append(params, charId)
	keys = append(keys, "char_id = ?")

	// key param [recipe_id] position [2] type [int]
	if len(c.QueryParam("recipe_id")) > 0 {
		recipeIdParam, err := strconv.Atoi(c.QueryParam("recipe_id"))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("Error parsing query param [recipe_id] err [%s]", err.Error())})
		}

		params = append(params, recipeIdParam)
		keys = append(keys, "recipe_id = ?")
	}

	// query builder
	var result models.CharRecipeList
	query := e.db.QueryContext(models.CharRecipeList{}, c)
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
		event := fmt.Sprintf("Updated [CharRecipeList] [%v] fields [%v]", strings.Join(ids, ", "), strings.Join(fieldsUpdated, ", "))
		e.auditLog.LogUserEvent(c, "UPDATE", event)
	}

	return c.JSON(http.StatusOK, request)
}

// createCharRecipeList godoc
// @Id createCharRecipeList
// @Summary Creates CharRecipeList
// @Accept json
// @Produce json
// @Param char_recipe_list body models.CharRecipeList true "CharRecipeList"
// @Tags CharRecipeList
// @Success 200 {array} models.CharRecipeList
// @Failure 500 {string} string "Error binding to entity"
// @Failure 500 {string} string "Error inserting entity"
// @Router /char_recipe_list [put]
func (e *CharRecipeListController) createCharRecipeList(c echo.Context) error {
	charRecipeList := new(models.CharRecipeList)
	if err := c.Bind(charRecipeList); err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			echo.Map{"error": fmt.Sprintf("Error binding to entity [%v]", err.Error())},
		)
	}

	db := e.db.Get(models.CharRecipeList{}, c).Model(&models.CharRecipeList{})

	// save associations
	if c.QueryParam("save_associations") != "true" {
		db = db.Omit(clause.Associations)
	}

	err := db.Create(&charRecipeList).Error
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			echo.Map{"error": fmt.Sprintf("Error inserting entity [%v]", err.Error())},
		)
	}

	// log create event
	if e.db.GetSpireDb() != nil {
		// diff between an empty model and the created
		diff := database.ResultDifference(models.CharRecipeList{}, charRecipeList)
		// build fields created
		var fields []string
		for k, v := range diff {
			fields = append(fields, fmt.Sprintf("%v = %v", k, v))
		}
		// record event
		event := fmt.Sprintf("Created [CharRecipeList] [%v] data [%v]", charRecipeList.CharId, strings.Join(fields, ", "))
		e.auditLog.LogUserEvent(c, "CREATE", event)
	}

	return c.JSON(http.StatusOK, charRecipeList)
}

// deleteCharRecipeList godoc
// @Id deleteCharRecipeList
// @Summary Deletes CharRecipeList
// @Accept json
// @Produce json
// @Tags CharRecipeList
// @Param id path int true "charId"
// @Success 200 {string} string "Entity deleted successfully"
// @Failure 404 {string} string "Cannot find entity"
// @Failure 500 {string} string "Error binding to entity"
// @Failure 500 {string} string "Error deleting entity"
// @Router /char_recipe_list/{id} [delete]
func (e *CharRecipeListController) deleteCharRecipeList(c echo.Context) error {
	var params []interface{}
	var keys []string

	// primary key param
	charId, err := strconv.Atoi(c.Param("charId"))
	if err != nil {
		e.logger.Error(err)
	}
	params = append(params, charId)
	keys = append(keys, "char_id = ?")

	// key param [recipe_id] position [2] type [int]
	if len(c.QueryParam("recipe_id")) > 0 {
		recipeIdParam, err := strconv.Atoi(c.QueryParam("recipe_id"))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("Error parsing query param [recipe_id] err [%s]", err.Error())})
		}

		params = append(params, recipeIdParam)
		keys = append(keys, "recipe_id = ?")
	}

	// query builder
	var result models.CharRecipeList
	query := e.db.QueryContext(models.CharRecipeList{}, c)
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
		event := fmt.Sprintf("Deleted [CharRecipeList] [%v] keys [%v]", result.CharId, strings.Join(ids, ", "))
		e.auditLog.LogUserEvent(c, "DELETE", event)
	}

	return c.JSON(http.StatusOK, echo.Map{"success": "Entity deleted successfully"})
}

// getCharRecipeListsBulk godoc
// @Id getCharRecipeListsBulk
// @Summary Gets CharRecipeLists in bulk
// @Accept json
// @Produce json
// @Param Body body BulkFetchByIdsGetRequest true "body"
// @Tags CharRecipeList
// @Success 200 {array} models.CharRecipeList
// @Failure 500 {string} string "Bad query request"
// @Router /char_recipe_lists/bulk [post]
func (e *CharRecipeListController) getCharRecipeListsBulk(c echo.Context) error {
	var results []models.CharRecipeList

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

	err := e.db.QueryContext(models.CharRecipeList{}, c).Find(&results, r.IDs).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, results)
}

// getCharRecipeListsCount godoc
// @Id getCharRecipeListsCount
// @Summary Counts CharRecipeLists
// @Accept json
// @Produce json
// @Tags CharRecipeList
// @Param includes query string false "Relationships [all] for all [number] for depth of relationships to load or [.] separated relationship names"
// @Param where query string false "Filter on specific fields. Multiple conditions [.] separated Example: col_like_value.col2__val2"
// @Param whereOr query string false "Filter on specific fields (Chained ors). Multiple conditions [.] separated Example: col_like_value.col2__val2"
// @Param groupBy query string false "Group by field. Multiple conditions [.] separated Example: field1.field2"
// @Param limit query string false "Rows to limit in response (Default: 10,000)"
// @Param page query int 0 "Pagination page"
// @Param orderBy query string false "Order by [field]"
// @Param orderDirection query string false "Order by field direction"
// @Param select query string false "Column names [.] separated to fetch specific fields in response"
// @Success 200 {array} models.CharRecipeList
// @Failure 500 {string} string "Bad query request"
// @Router /char_recipe_lists/count [get]
func (e *CharRecipeListController) getCharRecipeListsCount(c echo.Context) error {
	var count int64
	err := e.db.QueryContext(models.CharRecipeList{}, c).Count(&count).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err})
	}

	return c.JSON(http.StatusOK, echo.Map{"count": count})
}