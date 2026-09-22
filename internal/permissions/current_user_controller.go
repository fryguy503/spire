package permissions

import (
	"net/http"

	"github.com/EQEmuTools/spire/internal/database"
	"github.com/EQEmuTools/spire/internal/http/request"
	"github.com/EQEmuTools/spire/internal/http/routes"
	"github.com/labstack/echo/v4"
)

// CurrentUserController requires authentication, but no resource grants: users
// must be able to discover their own permissions even when they have none.
type CurrentUserController struct {
	db          *database.Resolver
	permissions *Service
}

func NewCurrentUserController(db *database.Resolver, permissions *Service) *CurrentUserController {
	return &CurrentUserController{db: db, permissions: permissions}
}

func (p *CurrentUserController) Routes() []*routes.Route {
	return []*routes.Route{
		routes.RegisterRoute(http.MethodGet, "permissions/me", p.getCurrentPermissions, nil),
	}
}

// Resolve the caller's active connection on the server. Never accept a user or
// connection ID from the browser when reporting effective permissions.
func (p *CurrentUserController) getCurrentPermissions(c echo.Context) error {
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
