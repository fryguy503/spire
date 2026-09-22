package updater

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInstallExecutableStagesBeforeReplacing(t *testing.T) {
	dir := t.TempDir()
	live := filepath.Join(dir, "Spire with spaces.exe")
	source := filepath.Join(t.TempDir(), "release.exe")
	if err := os.WriteFile(live, []byte("old executable"), 0755); err != nil {
		t.Fatal(err)
	}
	for _, content := range []string{"new executable", "second executable"} {
		if err := os.WriteFile(source, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
		if err := installExecutable(source, live); err != nil {
			t.Fatal(err)
		}
		if actual, err := os.ReadFile(live); err != nil || string(actual) != content {
			t.Fatalf("installed file = %q, %v", actual, err)
		}
	}
	backups, err := filepath.Glob(live + ".spire-update-old-*")
	if err != nil || len(backups) != 2 {
		t.Fatalf("backups = %v, %v", backups, err)
	}
	staging, _ := filepath.Glob(filepath.Join(dir, ".spire-update-new-*"))
	if len(staging) != 0 {
		t.Fatalf("staging files remain: %v", staging)
	}
}

func TestInstallExecutablePreservesOriginalOnInvalidSource(t *testing.T) {
	for _, kind := range []string{"missing", "empty", "directory"} {
		t.Run(kind, func(t *testing.T) {
			dir := t.TempDir()
			live, source := filepath.Join(dir, "spire"), filepath.Join(dir, "release")
			if err := os.WriteFile(live, []byte("original"), 0755); err != nil {
				t.Fatal(err)
			}
			if kind == "empty" {
				if err := os.WriteFile(source, nil, 0644); err != nil {
					t.Fatal(err)
				}
			}
			if kind == "directory" {
				if err := os.Mkdir(source, 0755); err != nil {
					t.Fatal(err)
				}
			}
			if err := installExecutable(source, live); err == nil {
				t.Fatal("expected installation error")
			}
			if actual, err := os.ReadFile(live); err != nil || string(actual) != "original" {
				t.Fatalf("original executable changed: %q, %v", actual, err)
			}
		})
	}
}
