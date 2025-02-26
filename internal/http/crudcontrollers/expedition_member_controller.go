package crudcontrollers

import (
	"fmt"
	"github.com/Akkadius/spire/internal/database"
	"github.com/Akkadius/spire/internal/http/routes"
	"github.com/Akkadius/spire/internal/models"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type ExpeditionMemberController struct {
	db	 *database.DatabaseResolver
	logger *logrus.Logger
}

func NewExpeditionMemberController(
	db *database.DatabaseResolver,
	logger *logrus.Logger,
) *ExpeditionMemberController {
	return &ExpeditionMemberController{
		db:	 db,
		logger: logger,
	}
}

func (e *ExpeditionMemberController) Routes() []*routes.Route {
	return []*routes.Route{
		routes.RegisterRoute(http.MethodGet, "expedition_member/:id", e.getExpeditionMember, nil),
		routes.RegisterRoute(http.MethodGet, "expedition_members", e.listExpeditionMembers, nil),
		routes.RegisterRoute(http.MethodPut, "expedition_member", e.createExpeditionMember, nil),
		routes.RegisterRoute(http.MethodDelete, "expedition_member/:id", e.deleteExpeditionMember, nil),
		routes.RegisterRoute(http.MethodPatch, "expedition_member/:id", e.updateExpeditionMember, nil),
		routes.RegisterRoute(http.MethodPost, "expedition_members/bulk", e.getExpeditionMembersBulk, nil),
	}
}

// listExpeditionMembers godoc
// @Id listExpeditionMembers
// @Summary Lists ExpeditionMembers
// @Accept json
// @Produce json
// @Tags ExpeditionMember
// @Param includes query string false "Relationships [all] for all [number] for depth of relationships to load or [.] separated relationship names "
// @Param where query string false "Filter on specific fields. Multiple conditions [.] separated Example: col_like_value.col2__val2"
// @Param whereOr query string false "Filter on specific fields (Chained ors). Multiple conditions [.] separated Example: col_like_value.col2__val2"
// @Param groupBy query string false "Group by field. Multiple conditions [.] separated Example: field1.field2"
// @Param limit query string false "Rows to limit in response (Default: 10,000)"
// @Param orderBy query string false "Order by [field]"
// @Param orderDirection query string false "Order by field direction"
// @Param select query string false "Column names [.] separated to fetch specific fields in response"
// @Success 200 {array} models.ExpeditionMember
// @Failure 500 {string} string "Bad query request"
// @Router /expedition_members [get]
func (e *ExpeditionMemberController) listExpeditionMembers(c echo.Context) error {
	var results []models.ExpeditionMember
	err := e.db.QueryContext(models.ExpeditionMember{}, c).Find(&results).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err})
	}

	return c.JSON(http.StatusOK, results)
}

// getExpeditionMember godoc
// @Id getExpeditionMember
// @Summary Gets ExpeditionMember
// @Accept json
// @Produce json
// @Tags ExpeditionMember
// @Param id path int true "Id"
// @Param includes query string false "Relationships [all] for all [number] for depth of relationships to load or [.] separated relationship names "
// @Param select query string false "Column names [.] separated to fetch specific fields in response"
// @Success 200 {array} models.ExpeditionMember
// @Failure 404 {string} string "Entity not found"
// @Failure 500 {string} string "Cannot find param"
// @Failure 500 {string} string "Bad query request"
// @Router /expedition_member/{id} [get]
func (e *ExpeditionMemberController) getExpeditionMember(c echo.Context) error {
	var params []interface{}
	var keys []string

	// primary key param
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Cannot find param [ID]"})
	}
	params = append(params, id)
	keys = append(keys, "id = ?")

	// query builder
	var result models.ExpeditionMember
	query := e.db.QueryContext(models.ExpeditionMember{}, c)
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

// updateExpeditionMember godoc
// @Id updateExpeditionMember
// @Summary Updates ExpeditionMember
// @Accept json
// @Produce json
// @Tags ExpeditionMember
// @Param id path int true "Id"
// @Param expedition_member body models.ExpeditionMember true "ExpeditionMember"
// @Success 200 {array} models.ExpeditionMember
// @Failure 404 {string} string "Cannot find entity"
// @Failure 500 {string} string "Error binding to entity"
// @Failure 500 {string} string "Error updating entity"
// @Router /expedition_member/{id} [patch]
func (e *ExpeditionMemberController) updateExpeditionMember(c echo.Context) error {
	request := new(models.ExpeditionMember)
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
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Cannot find param [ID]"})
	}
	params = append(params, id)
	keys = append(keys, "id = ?")

	// query builder
	var result models.ExpeditionMember
	query := e.db.QueryContext(models.ExpeditionMember{}, c)
	for i, _ := range keys {
		query = query.Where(keys[i], params[i])
	}

	// grab first entry
	err = query.First(&result).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("Cannot find entity [%s]", err.Error())})
	}

	err = query.Select("*").Updates(&request).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("Error updating entity [%v]", err.Error())})
	}

	return c.JSON(http.StatusOK, request)
}

// createExpeditionMember godoc
// @Id createExpeditionMember
// @Summary Creates ExpeditionMember
// @Accept json
// @Produce json
// @Param expedition_member body models.ExpeditionMember true "ExpeditionMember"
// @Tags ExpeditionMember
// @Success 200 {array} models.ExpeditionMember
// @Failure 500 {string} string "Error binding to entity"
// @Failure 500 {string} string "Error inserting entity"
// @Router /expedition_member [put]
func (e *ExpeditionMemberController) createExpeditionMember(c echo.Context) error {
	expeditionMember := new(models.ExpeditionMember)
	if err := c.Bind(expeditionMember); err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			echo.Map{"error": fmt.Sprintf("Error binding to entity [%v]", err.Error())},
		)
	}

	err := e.db.Get(models.ExpeditionMember{}, c).Model(&models.ExpeditionMember{}).Create(&expeditionMember).Error
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			echo.Map{"error": fmt.Sprintf("Error inserting entity [%v]", err.Error())},
		)
	}

	return c.JSON(http.StatusOK, expeditionMember)
}

// deleteExpeditionMember godoc
// @Id deleteExpeditionMember
// @Summary Deletes ExpeditionMember
// @Accept json
// @Produce json
// @Tags ExpeditionMember
// @Param id path int true "id"
// @Success 200 {string} string "Entity deleted successfully"
// @Failure 404 {string} string "Cannot find entity"
// @Failure 500 {string} string "Error binding to entity"
// @Failure 500 {string} string "Error deleting entity"
// @Router /expedition_member/{id} [delete]
func (e *ExpeditionMemberController) deleteExpeditionMember(c echo.Context) error {
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
	var result models.ExpeditionMember
	query := e.db.QueryContext(models.ExpeditionMember{}, c)
	for i, _ := range keys {
		query = query.Where(keys[i], params[i])
	}

	// grab first entry
	err = query.First(&result).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	err = e.db.Get(models.ExpeditionMember{}, c).Model(&models.ExpeditionMember{}).Delete(&result).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Error deleting entity"})
	}

	return c.JSON(http.StatusOK, echo.Map{"success": "Entity deleted successfully"})
}

// getExpeditionMembersBulk godoc
// @Id getExpeditionMembersBulk
// @Summary Gets ExpeditionMembers in bulk
// @Accept json
// @Produce json
// @Param Body body BulkFetchByIdsGetRequest true "body"
// @Tags ExpeditionMember
// @Success 200 {array} models.ExpeditionMember
// @Failure 500 {string} string "Bad query request"
// @Router /expedition_members/bulk [post]
func (e *ExpeditionMemberController) getExpeditionMembersBulk(c echo.Context) error {
	var results []models.ExpeditionMember

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

	err := e.db.QueryContext(models.ExpeditionMember{}, c).Find(&results, r.IDs).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, results)
}
