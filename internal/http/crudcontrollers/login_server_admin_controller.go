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

type LoginServerAdminController struct {
	db       *database.DatabaseResolver
	logger   *logrus.Logger
	auditLog *auditlog.UserEvent
}

func NewLoginServerAdminController(
	db *database.DatabaseResolver,
	logger *logrus.Logger,
	auditLog *auditlog.UserEvent,
) *LoginServerAdminController {
	return &LoginServerAdminController{
		db:       db,
		logger:   logger,
		auditLog: auditLog,
	}
}

func (e *LoginServerAdminController) Routes() []*routes.Route {
	return []*routes.Route{
		routes.RegisterRoute(http.MethodGet, "login_server_admin/:id", e.getLoginServerAdmin, nil),
		routes.RegisterRoute(http.MethodGet, "login_server_admins", e.listLoginServerAdmins, nil),
		routes.RegisterRoute(http.MethodGet, "login_server_admins/count", e.getLoginServerAdminsCount, nil),
		routes.RegisterRoute(http.MethodPut, "login_server_admin", e.createLoginServerAdmin, nil),
		routes.RegisterRoute(http.MethodDelete, "login_server_admin/:id", e.deleteLoginServerAdmin, nil),
		routes.RegisterRoute(http.MethodPatch, "login_server_admin/:id", e.updateLoginServerAdmin, nil),
		routes.RegisterRoute(http.MethodPost, "login_server_admins/bulk", e.getLoginServerAdminsBulk, nil),
	}
}

// listLoginServerAdmins godoc
// @Id listLoginServerAdmins
// @Summary Lists LoginServerAdmins
// @Accept json
// @Produce json
// @Tags LoginServerAdmin
// @Param includes query string false "Relationships [all] for all [number] for depth of relationships to load or [.] separated relationship names"
// @Param where query string false "Filter on specific fields. Multiple conditions [.] separated Example: col_like_value.col2__val2"
// @Param whereOr query string false "Filter on specific fields (Chained ors). Multiple conditions [.] separated Example: col_like_value.col2__val2"
// @Param groupBy query string false "Group by field. Multiple conditions [.] separated Example: field1.field2"
// @Param limit query string false "Rows to limit in response (Default: 10,000)"
// @Param page query int 0 "Pagination page"
// @Param orderBy query string false "Order by [field]"
// @Param orderDirection query string false "Order by field direction"
// @Param select query string false "Column names [.] separated to fetch specific fields in response"
// @Success 200 {array} models.LoginServerAdmin
// @Failure 500 {string} string "Bad query request"
// @Router /login_server_admins [get]
func (e *LoginServerAdminController) listLoginServerAdmins(c echo.Context) error {
	var results []models.LoginServerAdmin
	err := e.db.QueryContext(models.LoginServerAdmin{}, c).Find(&results).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err})
	}

	return c.JSON(http.StatusOK, results)
}

// getLoginServerAdmin godoc
// @Id getLoginServerAdmin
// @Summary Gets LoginServerAdmin
// @Accept json
// @Produce json
// @Tags LoginServerAdmin
// @Param id path int true "Id"
// @Param includes query string false "Relationships [all] for all [number] for depth of relationships to load or [.] separated relationship names"
// @Param select query string false "Column names [.] separated to fetch specific fields in response"
// @Success 200 {array} models.LoginServerAdmin
// @Failure 404 {string} string "Entity not found"
// @Failure 500 {string} string "Cannot find param"
// @Failure 500 {string} string "Bad query request"
// @Router /login_server_admin/{id} [get]
func (e *LoginServerAdminController) getLoginServerAdmin(c echo.Context) error {
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
	var result models.LoginServerAdmin
	query := e.db.QueryContext(models.LoginServerAdmin{}, c)
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

// updateLoginServerAdmin godoc
// @Id updateLoginServerAdmin
// @Summary Updates LoginServerAdmin
// @Accept json
// @Produce json
// @Tags LoginServerAdmin
// @Param id path int true "Id"
// @Param login_server_admin body models.LoginServerAdmin true "LoginServerAdmin"
// @Success 200 {array} models.LoginServerAdmin
// @Failure 404 {string} string "Cannot find entity"
// @Failure 500 {string} string "Error binding to entity"
// @Failure 500 {string} string "Error updating entity"
// @Router /login_server_admin/{id} [patch]
func (e *LoginServerAdminController) updateLoginServerAdmin(c echo.Context) error {
	request := new(models.LoginServerAdmin)
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
	var result models.LoginServerAdmin
	query := e.db.QueryContext(models.LoginServerAdmin{}, c)
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
		event := fmt.Sprintf("Updated [LoginServerAdmin] [%v] fields [%v]", strings.Join(ids, ", "), strings.Join(fieldsUpdated, ", "))
		e.auditLog.LogUserEvent(c, "UPDATE", event)
	}

	return c.JSON(http.StatusOK, request)
}

// createLoginServerAdmin godoc
// @Id createLoginServerAdmin
// @Summary Creates LoginServerAdmin
// @Accept json
// @Produce json
// @Param login_server_admin body models.LoginServerAdmin true "LoginServerAdmin"
// @Tags LoginServerAdmin
// @Success 200 {array} models.LoginServerAdmin
// @Failure 500 {string} string "Error binding to entity"
// @Failure 500 {string} string "Error inserting entity"
// @Router /login_server_admin [put]
func (e *LoginServerAdminController) createLoginServerAdmin(c echo.Context) error {
	loginServerAdmin := new(models.LoginServerAdmin)
	if err := c.Bind(loginServerAdmin); err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			echo.Map{"error": fmt.Sprintf("Error binding to entity [%v]", err.Error())},
		)
	}

	db := e.db.Get(models.LoginServerAdmin{}, c).Model(&models.LoginServerAdmin{})

	// save associations
	if c.QueryParam("save_associations") != "true" {
		db = db.Omit(clause.Associations)
	}

	err := db.Create(&loginServerAdmin).Error
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			echo.Map{"error": fmt.Sprintf("Error inserting entity [%v]", err.Error())},
		)
	}

	// log create event
	if e.db.GetSpireDb() != nil {
		// diff between an empty model and the created
		diff := database.ResultDifference(models.LoginServerAdmin{}, loginServerAdmin)
		// build fields created
		var fields []string
		for k, v := range diff {
			fields = append(fields, fmt.Sprintf("%v = %v", k, v))
		}
		// record event
		event := fmt.Sprintf("Created [LoginServerAdmin] [%v] data [%v]", loginServerAdmin.ID, strings.Join(fields, ", "))
		e.auditLog.LogUserEvent(c, "CREATE", event)
	}

	return c.JSON(http.StatusOK, loginServerAdmin)
}

// deleteLoginServerAdmin godoc
// @Id deleteLoginServerAdmin
// @Summary Deletes LoginServerAdmin
// @Accept json
// @Produce json
// @Tags LoginServerAdmin
// @Param id path int true "id"
// @Success 200 {string} string "Entity deleted successfully"
// @Failure 404 {string} string "Cannot find entity"
// @Failure 500 {string} string "Error binding to entity"
// @Failure 500 {string} string "Error deleting entity"
// @Router /login_server_admin/{id} [delete]
func (e *LoginServerAdminController) deleteLoginServerAdmin(c echo.Context) error {
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
	var result models.LoginServerAdmin
	query := e.db.QueryContext(models.LoginServerAdmin{}, c)
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
		event := fmt.Sprintf("Deleted [LoginServerAdmin] [%v] keys [%v]", result.ID, strings.Join(ids, ", "))
		e.auditLog.LogUserEvent(c, "DELETE", event)
	}

	return c.JSON(http.StatusOK, echo.Map{"success": "Entity deleted successfully"})
}

// getLoginServerAdminsBulk godoc
// @Id getLoginServerAdminsBulk
// @Summary Gets LoginServerAdmins in bulk
// @Accept json
// @Produce json
// @Param Body body BulkFetchByIdsGetRequest true "body"
// @Tags LoginServerAdmin
// @Success 200 {array} models.LoginServerAdmin
// @Failure 500 {string} string "Bad query request"
// @Router /login_server_admins/bulk [post]
func (e *LoginServerAdminController) getLoginServerAdminsBulk(c echo.Context) error {
	var results []models.LoginServerAdmin

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

	err := e.db.QueryContext(models.LoginServerAdmin{}, c).Find(&results, r.IDs).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, results)
}

// getLoginServerAdminsCount godoc
// @Id getLoginServerAdminsCount
// @Summary Counts LoginServerAdmins
// @Accept json
// @Produce json
// @Tags LoginServerAdmin
// @Param includes query string false "Relationships [all] for all [number] for depth of relationships to load or [.] separated relationship names"
// @Param where query string false "Filter on specific fields. Multiple conditions [.] separated Example: col_like_value.col2__val2"
// @Param whereOr query string false "Filter on specific fields (Chained ors). Multiple conditions [.] separated Example: col_like_value.col2__val2"
// @Param groupBy query string false "Group by field. Multiple conditions [.] separated Example: field1.field2"
// @Param limit query string false "Rows to limit in response (Default: 10,000)"
// @Param page query int 0 "Pagination page"
// @Param orderBy query string false "Order by [field]"
// @Param orderDirection query string false "Order by field direction"
// @Param select query string false "Column names [.] separated to fetch specific fields in response"
// @Success 200 {array} models.LoginServerAdmin
// @Failure 500 {string} string "Bad query request"
// @Router /login_server_admins/count [get]
func (e *LoginServerAdminController) getLoginServerAdminsCount(c echo.Context) error {
	var count int64
	err := e.db.QueryContext(models.LoginServerAdmin{}, c).Count(&count).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err})
	}

	return c.JSON(http.StatusOK, echo.Map{"count": count})
}