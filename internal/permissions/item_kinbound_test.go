package permissions

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EQEmuTools/spire/internal/models"
	"github.com/labstack/echo/v4"
	gocache "github.com/patrickmn/go-cache"
)

func TestItemKinboundPermission(t *testing.T) {
	for _, tc := range []struct {
		name, method      string
		read, write, want bool
	}{
		{"reader", http.MethodGet, true, false, true},
		{"reader cannot save", http.MethodPatch, true, false, false},
		{"writer", http.MethodPatch, true, true, true},
		{"no grant", http.MethodGet, false, false, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cache := gocache.New(gocache.NoExpiration, 0)
			service := &Service{cache: cache}
			prefixes := service.RegisterManualResources()["Item Kinbound"]
			if len(prefixes) != 1 || prefixes[0] != "item-kinbound" {
				t.Fatalf("unexpected Kinbound resource: %v", prefixes)
			}
			cache.Set("user-permissions-2", userPermissions{connectionID: 1, permissions: []Resource{{RouteMatchPrefixes: prefixes, CanRead: tc.read, CanWrite: tc.write}}}, gocache.NoExpiration)
			c := echo.New().NewContext(httptest.NewRequest(tc.method, "/api/v1/item-kinbound/100", nil), httptest.NewRecorder())
			if got := service.CanAccessResource(c, models.User{ID: 2}, 1); got != tc.want {
				t.Fatalf("access=%v, want %v", got, tc.want)
			}
		})
	}
}
