package main

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/EQEmuTools/spire/boot"
	"github.com/EQEmuTools/spire/internal/console"
	"github.com/EQEmuTools/spire/internal/env"
	"github.com/EQEmuTools/spire/internal/selfrestart"
	spiresentry "github.com/EQEmuTools/spire/internal/sentry"
	"github.com/EQEmuTools/spire/internal/updater"
	"github.com/henvic/httpretty"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	if handled, code, err := selfrestart.Run(); handled {
		if err != nil {
			log.Print(err)
		}
		os.Exit(code)
	}
	updater.CleanupOldExecutables()
	// uncomment for profiling
	//f, err := os.Create("main.prof")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//pprof.StartCPUProfile(f)
	//defer pprof.StopCPUProfile()

	registerHttpLogging()

	// self update service
	if len(os.Args) == 1 && !selfrestart.ResumedDesktop() {
		if updated, err := updater.NewUpdater(packageJson).CheckForUpdates(context.Background(), true); err != nil {
			log.Printf("Could not update Spire: %v", err)
		} else if updated != "" {
			if err := selfrestart.Prepare(); err != nil {
				log.Printf("Could not restart Spire: %v", err)
			} else {
				selfrestart.Exit()
			}
		}
	}

	// default (only if not provided by environment)
	if len(os.Getenv("APP_ENV")) == 0 {
		_ = os.Setenv("APP_ENV", "local")
	}

	// load env
	if err := env.LoadEnvFileIfExists(); err != nil {
		Fatal(err)
	}

	// initialize sentry (optional - only if SENTRY_DSN is set)
	if err := spiresentry.Init(); err != nil {
		log.Printf("Warning: Sentry initialization failed: %v", err)
	}
	defer spiresentry.Flush()

	// boot app
	app, err := boot.InitializeApplication()
	if err != nil {
		Fatal(err)
	}

	// load embedded resources
	loadEmbedded(&app)

	// ran via executable on desktop
	if len(os.Args) == 1 {
		_ = os.Setenv("APP_ENV", "desktop")
		app.Desktop().Boot()
	}

	// run cmd
	if err := console.Run(app.Commands()); err != nil {
		Fatal(err)
	}
}

func Fatal(err error) {
	log.Println(err)

	// only hang if we're not running a command
	if len(os.Args) == 1 {
		fmt.Print("Automatically shutting down in 10 seconds...")
		time.Sleep(10 * time.Second)
		os.Exit(1)
	}
}

var (
	//go:embed CHANGELOG.md
	changelog string
	//go:embed package.json
	packageJson []byte
)

func loadEmbedded(app *boot.App) {
	// load embedded files into cache
	app.Cache().Set("changelog", changelog, -1)
	app.Cache().Set("packageJson", packageJson, -1)
}

func registerHttpLogging() {
	if len(os.Getenv("HTTP_DEBUG")) > 0 {
		logger := &httpretty.Logger{
			Time:           true,
			TLS:            true,
			RequestHeader:  true,
			RequestBody:    true,
			ResponseHeader: true,
			ResponseBody:   true,
			Colors:         true,
			Formatters:     []httpretty.Formatter{&httpretty.JSONFormatter{}},
		}

		http.DefaultClient.Transport = logger.RoundTripper(http.DefaultClient.Transport)
	}
}
