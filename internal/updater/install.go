package updater

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Stage on the destination filesystem before renaming the running executable.
// A failed copy must never leave a partially written executable at the live path.
func installExecutable(source, executable string) error {
	input, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("open updated executable: %w", err)
	}
	defer input.Close()
	info, err := input.Stat()
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() || info.Size() == 0 {
		return fmt.Errorf("updated executable is empty or not a regular file")
	}
	staged, err := os.CreateTemp(filepath.Dir(executable), ".spire-update-new-*")
	if err != nil {
		return err
	}
	defer os.Remove(staged.Name())
	defer staged.Close()
	if _, err := io.Copy(staged, input); err != nil {
		return fmt.Errorf("stage updated executable: %w", err)
	}
	if err := staged.Chmod(0755); err != nil {
		return err
	}
	if err := staged.Sync(); err != nil {
		return err
	}
	if err := staged.Close(); err != nil {
		return err
	}
	backup, err := os.CreateTemp(filepath.Dir(executable), filepath.Base(executable)+".spire-update-old-*")
	if err != nil {
		return err
	}
	backupPath := backup.Name()
	if err := backup.Close(); err != nil {
		return err
	}
	if err := os.Remove(backupPath); err != nil {
		return err
	}
	if err := os.Rename(executable, backupPath); err != nil {
		return fmt.Errorf("back up running executable: %w", err)
	}
	if err := os.Rename(staged.Name(), executable); err != nil {
		if rollbackErr := os.Rename(backupPath, executable); rollbackErr != nil {
			return fmt.Errorf("install update: %v; restore executable from %s: %w", err, backupPath, rollbackErr)
		}
		return fmt.Errorf("install update: %w", err)
	}
	// Keep the backup until the replacement starts. Windows cannot delete it
	// while the old server or launcher still has the executable open.
	return nil
}

// CleanupOldExecutables removes only backups created by this updater. On
// Windows the launcher's own backup stays locked until the launcher exits.
func CleanupOldExecutables() {
	executable, err := os.Executable()
	if err != nil {
		return
	}
	entries, err := os.ReadDir(filepath.Dir(executable))
	if err != nil {
		return
	}
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasPrefix(entry.Name(), filepath.Base(executable)+".spire-update-old-") {
			_ = os.Remove(filepath.Join(filepath.Dir(executable), entry.Name()))
		}
	}
}
