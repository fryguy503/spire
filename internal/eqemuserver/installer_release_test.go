package eqemuserver

import (
	"testing"

	spirerelease "github.com/EQEmuTools/spire/internal/release"
)

func TestSpireInstallerReleaseRepositoryUsesCombinedFork(t *testing.T) {
	owner, repo, err := spireInstallerReleaseRepository()
	if err != nil {
		t.Fatalf("spireInstallerReleaseRepository() returned error: %v", err)
	}

	got := owner + "/" + repo
	if got != spirerelease.DefaultRepository {
		t.Fatalf("spireInstallerReleaseRepository() = %q, want configured default %q", got, spirerelease.DefaultRepository)
	}
	if got != "fryguy503/spire" {
		t.Fatalf("spireInstallerReleaseRepository() = %q, want fryguy503/spire", got)
	}
}
