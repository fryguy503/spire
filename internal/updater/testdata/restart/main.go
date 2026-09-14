// This fixture runs the production update handler, installer, and launcher
// against a local release server. It never contacts GitHub or a database.
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/EQEmuTools/spire/internal/app"
	"github.com/EQEmuTools/spire/internal/selfrestart"
	"github.com/EQEmuTools/spire/internal/updater"
	"github.com/labstack/echo/v4"
	"github.com/patrickmn/go-cache"
)

var version = "1.0.0"

func main() {
	if os.Getenv("SPIRE_INTERNAL_RESTART_STATE") == "" {
		_ = os.Setenv("SPIRE_TEST_LAUNCHER_VERSION", version)
	}
	if handled, code, err := runLauncher(); handled {
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(code)
	}
	updater.CleanupOldExecutables()
	base, err := url.Parse(os.Getenv("SPIRE_TEST_RELEASE_URL"))
	if err != nil {
		panic(err)
	}
	// Updater's GitHub client uses the default transport when no transport is
	// configured. Substitute only in this fixture, never in production.
	originalTransport := http.DefaultTransport
	http.DefaultTransport = transportFunc(func(request *http.Request) (*http.Response, error) {
		if request.URL.Host == "api.github.com" {
			request = request.Clone(request.Context())
			request.URL.Scheme, request.URL.Host = base.Scheme, base.Host
		}
		return originalTransport.RoundTrip(request)
	})
	port, _ := strconv.Atoi(os.Args[3])
	if selfrestart.DesktopPort() > 0 {
		port = selfrestart.DesktopPort()
	}
	selfrestart.SetDesktopPort(port)
	data, _ := json.Marshal(map[string]string{"version": version})
	metadata := cache.New(0, 0)
	metadata.Set("packageJson", data, cache.NoExpiration)
	controller := app.NewController(metadata, nil, nil, nil, nil, nil, nil)
	e := echo.New()
	e.HideBanner = true
	e.POST("/block-log", func(c echo.Context) error {
		go func() { _, _ = os.Stderr.WriteString(strings.Repeat("x", 1024*1024)) }()
		return c.NoContent(200)
	})
	e.POST("/restart", func(c echo.Context) error {
		if err := selfrestart.Prepare(); err != nil {
			return err
		}
		err := c.JSON(200, echo.Map{"restarting": true})
		c.Response().Flush()
		time.AfterFunc(10*time.Millisecond, selfrestart.Exit)
		return err
	})
	e.POST("/stop", func(c echo.Context) error {
		code, _ := strconv.Atoi(c.QueryParam("code"))
		os.Exit(code)
		return nil
	})
	for _, route := range controller.Routes() {
		if route.Route() == "app/update" {
			e.POST("/api/v1/app/update", route.Handler())
		}
	}
	e.GET("/api/v1/app/env", func(c echo.Context) error {
		cwd, _ := os.Getwd()
		return c.JSON(200, echo.Map{"data": echo.Map{
			"version": version, "pid": os.Getpid(), "parent": os.Getppid(),
			"cwd": cwd, "args": os.Args[1:], "port": selfrestart.DesktopPort(),
			"resumed": selfrestart.ResumedDesktop(), "marker": os.Getenv("SPIRE_TEST_MARKER"),
			"worker_env":       os.Getenv("SPIRE_INTERNAL_RESTART_STATE"),
			"launcher_version": os.Getenv("SPIRE_TEST_LAUNCHER_VERSION"),
		}})
	})
	if err := e.Start(fmt.Sprintf("127.0.0.1:%d", port)); err != nil {
		panic(err)
	}
}

func runLauncher() (bool, int, error) {
	if os.Getenv("SPIRE_TEST_DIRECT") == "1" {
		return false, 0, nil
	}
	return selfrestart.Run()
}

type transportFunc func(*http.Request) (*http.Response, error)

func (f transportFunc) RoundTrip(request *http.Request) (*http.Response, error) { return f(request) }
