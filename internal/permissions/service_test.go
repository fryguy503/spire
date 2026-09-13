package permissions

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EQEmuTools/spire/internal/models"
	"github.com/labstack/echo/v4"
	gocache "github.com/patrickmn/go-cache"
)

func TestConnectionPermissionsWithoutInstanceAdmin(t *testing.T) {
	for _, tc := range []struct {
		name        string
		method      string
		permissions userPermissions
		want        bool
	}{
		{name: "all read", method: http.MethodGet, permissions: userPermissions{canReadAll: true}, want: true},
		{name: "all write", method: http.MethodPut, permissions: userPermissions{canWriteAll: true}, want: true},
		{name: "read only cannot write", method: http.MethodPut, permissions: userPermissions{canReadAll: true}},
		{name: "write only cannot read", method: http.MethodGet, permissions: userPermissions{canWriteAll: true}},
		{name: "no grants cannot read", method: http.MethodGet},
		{name: "connection owner", method: http.MethodPut, permissions: userPermissions{isConnectionOwner: true}, want: true},
		{
			name: "scoped server read", method: http.MethodGet, want: true,
			permissions: userPermissions{permissions: []Resource{{
				RouteMatchPrefixes: []string{"admin/serverconfig"}, CanRead: true,
			}}},
		},
		{
			name: "unrelated grant cannot read", method: http.MethodGet,
			permissions: userPermissions{permissions: []Resource{{
				RouteMatchPrefixes: []string{"items"}, CanRead: true, CanWrite: true,
			}}},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cache := gocache.New(gocache.NoExpiration, 0)
			// Seed the same permission cache used after resolving a connection.
			cache.Set("user-permissions-2", tc.permissions, gocache.NoExpiration)
			service := &Service{cache: cache}
			e := echo.New()
			request := httptest.NewRequest(tc.method, "/api/v1/admin/serverconfig", nil)
			context := e.NewContext(request, httptest.NewRecorder())
			user := models.User{ID: 2, IsAdmin: false}
			if got := service.CanAccessResource(context, user, 1); got != tc.want {
				t.Errorf("CanAccessResource() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestIsWriteRequestIncludesDelete(t *testing.T) {
	service := &Service{}
	e := echo.New()

	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete} {
		request := httptest.NewRequest(method, "/api/v1/faction-editor/faction/10", nil)
		context := e.NewContext(request, httptest.NewRecorder())
		if !service.IsWriteRequest(context) {
			t.Errorf("IsWriteRequest(%s) = false, want true", method)
		}
	}
}

func TestIsWriteRequestTreatsReadsAndBulkAsNonWrites(t *testing.T) {
	service := &Service{}
	e := echo.New()

	for _, test := range []struct {
		method string
		path   string
	}{
		{method: http.MethodGet, path: "/api/v1/faction-editor/factions"},
		{method: http.MethodPost, path: "/api/v1/items/bulk"},
	} {
		request := httptest.NewRequest(test.method, test.path, nil)
		context := e.NewContext(request, httptest.NewRecorder())
		if service.IsWriteRequest(context) {
			t.Errorf("IsWriteRequest(%s %s) = true, want false", test.method, test.path)
		}
	}
}

func TestFactionEditorIsRegisteredAsManualResource(t *testing.T) {
	resources := (&Service{}).RegisterManualResources()
	prefixes, ok := resources["Faction Editor"]
	if !ok || len(prefixes) != 1 || prefixes[0] != "faction-editor" {
		t.Fatalf("Faction Editor resource = %#v, want [faction-editor]", prefixes)
	}
}

func TestContentFlagEditorIsRegisteredAsManualResource(t *testing.T) {
	resources := (&Service{}).RegisterManualResources()
	prefixes, ok := resources["Content Flag Editor"]
	if !ok || len(prefixes) != 1 || prefixes[0] != "content-flag-editor" {
		t.Fatalf("Content Flag Editor resource = %#v, want [content-flag-editor]", prefixes)
	}
}

func TestAlternateCurrencyEditorIsRegisteredAsManualResource(t *testing.T) {
	resources := (&Service{}).RegisterManualResources()
	prefixes, ok := resources["Alternate Currency Editor"]
	if !ok || len(prefixes) != 1 || prefixes[0] != "alternate-currency-editor" {
		t.Fatalf("Alternate Currency Editor resource = %#v, want [alternate-currency-editor]", prefixes)
	}
}

func TestPlayerOperationsIsRegisteredAsManualResource(t *testing.T) {
	resources := (&Service{}).RegisterManualResources()
	prefixes, ok := resources["Player Operations"]
	if !ok || len(prefixes) != 1 || prefixes[0] != "player-operations" {
		t.Fatalf("Player Operations resource = %#v, want [player-operations]", prefixes)
	}
}

func TestMailParcelsEditorIsRegisteredAsManualResource(t *testing.T) {
	resources := (&Service{}).RegisterManualResources()
	prefixes, ok := resources["Mail & Parcels Editor"]
	if !ok || len(prefixes) != 1 || prefixes[0] != "mail-parcels-editor" {
		t.Fatalf("Mail & Parcels Editor resource = %#v, want [mail-parcels-editor]", prefixes)
	}
}

func TestInventoryKeyringIsRegisteredAsManualResource(t *testing.T) {
	resources := (&Service{}).RegisterManualResources()
	prefixes, ok := resources["Inventory & Keyring"]
	if !ok || len(prefixes) != 1 || prefixes[0] != "inventory-keyring" {
		t.Fatalf("Inventory & Keyring resource = %#v, want [inventory-keyring]", prefixes)
	}
}

func TestOperationalDataEditorsAreRegisteredAsManualResources(t *testing.T) {
	resources := (&Service{}).RegisterManualResources()
	tests := map[string]string{
		"Data Buckets Editor": "data-bucket-editor",
		"QGlobals Editor":     "qglobal-editor",
		"Chat Administration": "chat-administration",
	}
	for name, prefix := range tests {
		prefixes, ok := resources[name]
		if !ok || len(prefixes) != 1 || prefixes[0] != prefix {
			t.Fatalf("%s resource = %#v, want [%s]", name, prefixes, prefix)
		}
	}
}

func TestSpireApplicationUpdateIsRegisteredAsManualResource(t *testing.T) {
	resources := (&Service{}).RegisterManualResources()
	prefixes, ok := resources["Spire Application Update"]
	if !ok || len(prefixes) != 1 || prefixes[0] != "app/update" {
		t.Fatalf("Spire Application Update resource = %#v, want [app/update]", prefixes)
	}
}
