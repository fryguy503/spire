package permissions

import (
	"github.com/EQEmuTools/spire/internal/database"
	"github.com/EQEmuTools/spire/internal/http/request"
	"github.com/EQEmuTools/spire/internal/http/routes"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Controller struct {
	db          *database.Resolver
	permissions *Service
}

func NewController(
	db *database.Resolver,
	permissions *Service,
) *Controller {
	return &Controller{
		db:          db,
		permissions: permissions,
	}
}

func (p *Controller) Routes() []*routes.Route {
	return []*routes.Route{
		routes.RegisterRoute(http.MethodGet, "permissions/resources", p.getPermissionResources, nil),
		routes.RegisterRoute(http.MethodGet, "permissions/me", p.getCurrentPermissions, nil),
	}
}

// Resolve the caller's active connection on the server. Never accept a user or
// connection ID from the browser when reporting effective permissions.
func (p *Controller) getCurrentPermissions(c echo.Context) error {
	user := request.GetUser(c)
	if user.ID == 0 {
		return echo.NewHTTPError(http.StatusUnauthorized, "Authentication required")
	}
	connection := p.db.GetUserConnection(user)
	if connection.ServerDatabaseConnectionId == 0 {
		return c.JSON(http.StatusOK, permissionSnapshot{})
	}
	grants := p.permissions.getUserPermissions(c, user, connection.ServerDatabaseConnectionId)
	return c.JSON(http.StatusOK, grants.snapshot(connection.ServerDatabaseConnectionId))
}

func (p *Controller) getPermissionResources(c echo.Context) error {
	return c.JSON(http.StatusOK, p.permissions.GetResources(c.Echo().Routes()))
}
