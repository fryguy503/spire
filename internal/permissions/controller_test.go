package permissions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/EQEmuTools/spire/internal/database"
	"github.com/EQEmuTools/spire/internal/models"
	"github.com/labstack/echo/v4"
	gocache "github.com/patrickmn/go-cache"
)

func TestCurrentPermissionsRequiresAuthentication(t *testing.T) {
	e := echo.New()
	context := e.NewContext(httptest.NewRequest(http.MethodGet, "/api/v1/permissions/me", nil), httptest.NewRecorder())
	err := (&Controller{}).getCurrentPermissions(context)
	if httpError, ok := err.(*echo.HTTPError); !ok || httpError.Code != http.StatusUnauthorized {
		t.Fatalf("getCurrentPermissions() error = %v, want 401", err)
	}
}

func TestCurrentPermissionsUsesActiveConnectionACL(t *testing.T) {
	for _, tc := range []struct {
		name       string
		connection uint
		grants     userPermissions
		want       permissionSnapshot
	}{
		{name: "no active connection"},
		{name: "no grants despite instance admin", connection: 7, want: permissionSnapshot{ConnectionID: 7}},
		{name: "owner", connection: 7, grants: userPermissions{isConnectionOwner: true}, want: permissionSnapshot{ConnectionID: 7, ReadAll: true, WriteAll: true}},
		{name: "all read", connection: 7, grants: userPermissions{canReadAll: true}, want: permissionSnapshot{ConnectionID: 7, ReadAll: true}},
		{name: "all write", connection: 7, grants: userPermissions{canWriteAll: true}, want: permissionSnapshot{ConnectionID: 7, WriteAll: true}},
		{
			name: "scoped read and write", connection: 7,
			grants: userPermissions{permissions: []Resource{
				{RouteMatchPrefixes: []string{"variable", "variables"}, CanRead: true},
				{RouteMatchPrefixes: []string{"admin/serverconfig"}, CanWrite: true},
			}},
			want: permissionSnapshot{ConnectionID: 7, Read: []string{"variable", "variables"}, Write: []string{"admin/serverconfig"}},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cache := gocache.New(gocache.NoExpiration, 0)
			// Membership ID deliberately differs from the underlying connection ID.
			cache.Set("active-user-db-connection-2", models.UserServerDatabaseConnection{
				ID: 99, UserId: 2, ServerDatabaseConnectionId: tc.connection, Active: 1,
			}, gocache.NoExpiration)
			tc.grants.connectionID = tc.connection
			cache.Set("user-permissions-2", tc.grants, gocache.NoExpiration)
			resolver := database.NewResolver(nil, nil, nil, cache)
			controller := NewController(resolver, &Service{cache: cache})
			e := echo.New()
			recorder := httptest.NewRecorder()
			// Query parameters cannot select another user's permissions.
			context := e.NewContext(httptest.NewRequest(http.MethodGet, "/api/v1/permissions/me?user_id=1&connection_id=99", nil), recorder)
			context.Set("user", models.User{ID: 2, IsAdmin: true})
			if err := controller.getCurrentPermissions(context); err != nil {
				t.Fatal(err)
			}
			if recorder.Code != http.StatusOK {
				t.Fatalf("status = %v", recorder.Code)
			}
			var got permissionSnapshot
			if err := json.Unmarshal(recorder.Body.Bytes(), &got); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("snapshot = %#v, want %#v", got, tc.want)
			}
			controller.permissions.ClearUserPermissionsCache(2)
			if _, found := cache.Get(fmt.Sprintf("user-permissions-%v", 2)); found {
				t.Fatal("permission edits must invalidate the snapshot cache")
			}
		})
	}
}
