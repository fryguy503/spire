package updater

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/EQEmuTools/spire/internal/download"
	"github.com/EQEmuTools/spire/internal/env"
	"github.com/EQEmuTools/spire/internal/eqemuserverconfig"
	"github.com/EQEmuTools/spire/internal/logger"
	"github.com/EQEmuTools/spire/internal/pathmgmt"
	"github.com/EQEmuTools/spire/internal/unzip"
	"github.com/google/go-github/v41/github"
)

// Serialize replacement across updater instances, including the time between
// installation and the running process exiting.
var installMu sync.Mutex
var installedVersion string

// Updater is a service that checks for updates to the app.
type Updater struct {
	packageJson               []byte
	logger                    *logger.AppLogger
	serverconfig              *eqemuserverconfig.Config
	unzipper                  *unzip.Unzipper
	githubClient              *github.Client
	goos                      string
	goarch                    string
	releaseRepositoryOverride string
}

func NewUpdater(packageJson []byte) *Updater {
	appLogger := logger.ProvideAppLogger()
	pathmgr := pathmgmt.NewPathManagement(appLogger)
	return &Updater{
		packageJson:               packageJson,
		logger:                    appLogger,
		serverconfig:              eqemuserverconfig.NewConfig(appLogger, pathmgr),
		unzipper:                  unzip.NewUnzipper(appLogger),
		githubClient:              newDefaultGitHubClient(),
		goos:                      runtime.GOOS,
		goarch:                    runtime.GOARCH,
		releaseRepositoryOverride: os.Getenv("SPIRE_RELEASE_REPO"),
	}
}

type EnvResponse struct {
	Env     string `json:"env"`
	Version string `json:"version"`
}

type PackageJson struct {
	Name       string `json:"name"`
	Version    string `json:"version"`
	Repository struct {
		Type string `json:"type"`
		URL  string `json:"url"`
	} `json:"repository"`
}

func (s *Updater) getAppVersion() (error, EnvResponse) {
	var pkg PackageJson
	if err := json.Unmarshal(s.packageJson, &pkg); err != nil {
		return err, EnvResponse{}
	}
	return nil, EnvResponse{Env: env.Get("APP_ENV", "local"), Version: pkg.Version}
}

// CheckForUpdates installs an eligible update and returns its version, or an
// empty string when no update is needed. Errors leave the running app alive.
func (s *Updater) CheckForUpdates(interactive bool) (string, error) {
	installMu.Lock()
	defer installMu.Unlock()
	if installedVersion != "" {
		return installedVersion, nil
	}
	config, configErr := s.serverconfig.Get()
	if configErr == nil && config.Spire.DisableAutoUpdates && interactive {
		s.logger.Info().Msg("Auto updates are disabled via config")
		return "", nil
	}
	executable, err := os.Executable()
	if err != nil {
		return "", err
	}
	executable, err = filepath.EvalSymlinks(executable)
	if err != nil {
		return "", err
	}
	if name := filepath.Base(executable); name == "main.exe" || name == "main" {
		s.logger.Info().Msg("Running as go run main.go, ignoring updates")
		return "", nil
	}
	status, selected, err := s.resolveUpdate(context.Background())
	if err != nil {
		return "", fmt.Errorf("resolve Spire release: %w", err)
	}
	if !status.Available || selected == nil {
		s.logger.Info().Any("version", status.CurrentVersion).Msg("Spire is already up to date")
		return "", nil
	}
	tmpdir, err := os.MkdirTemp("", "spire-update-")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tmpdir)
	archiveName := targetReleaseAssetName(s.goos, s.goarch)
	downloadPath := filepath.Join(tmpdir, archiveName)
	if err := download.WithProgress(downloadPath, selected.asset.GetBrowserDownloadURL()); err != nil {
		return "", fmt.Errorf("download Spire update: %w", err)
	}
	if err := s.unzipper.Extract(downloadPath, tmpdir); err != nil {
		return "", fmt.Errorf("extract Spire update: %w", err)
	}
	// The Windows archive contains an .exe; the Linux archive does not.
	source := filepath.Join(tmpdir, strings.TrimSuffix(archiveName, ".zip"))
	if err := installExecutable(source, executable); err != nil {
		return "", err
	}
	installedVersion = strings.TrimPrefix(selected.release.GetTagName(), "v")
	_ = os.Remove(filepath.Join(os.TempDir(), "spire_asset_last_check"))
	s.logger.Info().Any("version", installedVersion).Msg("Spire successfully updated; restarting")
	return installedVersion, nil
}
