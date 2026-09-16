package generators

import (
	"os"
	"testing"
)

func TestRaceInventoryReference(t *testing.T) {
	contents, err := os.ReadFile("testdata/race-inventory.html")
	if err != nil {
		t.Fatal(err)
	}
	data := (&RaceModelMapsCommand{}).Parse(string(contents))
	races := map[int]RaceEntry{}
	for _, race := range data.Race {
		races[race.RaceId] = race
	}
	if len(races) != 7 {
		t.Fatalf("expected all 7 fixture races including IDs above 1000, got %d", len(races))
	}
	if races[119].Description != "Saber-toothed Cat" {
		t.Fatalf("hyphenated race name truncated: %q", races[119].Description)
	}
	if races[2253].MaleGenderModelCode != "PPOINT" || len(races[2253].Sources) != 0 {
		t.Fatal("source-free reference entry was lost or fabricated")
	}
	human := races[1]
	if !human.IsPlayable || len(human.Sources) != 3 || !human.Sources[0].IsGlobal || human.Sources[0].LoadedViaDescription != "Luclin models" {
		t.Fatal("global source and conditional loading information must be retained")
	}
	aviak := races[13]
	foundLocal, foundImported := false, false
	for _, source := range aviak.Sources {
		for _, model := range source.Models {
			if model.Code != "AVI" {
				t.Fatalf("unexpected Aviak code %q", model.Code)
			}
		}
		if source.SourceFile == "southkarana_chr.s3d" {
			for _, zone := range source.Zones {
				if zone.ID == 14 && zone.ShortName == "southkarana" && zone.ImportType == "Local" {
					foundLocal = true
				}
				if zone.ID == 383 && zone.ShortName == "freeportwest" && zone.ImportType == "Imported" {
					foundImported = true
				}
			}
		}
	}
	if !foundLocal || !foundImported {
		t.Fatal("archive-specific local/imported zone relationships were lost")
	}
	for id, code := range map[int]string{487: "BSE", 724: "LUC"} {
		models := races[id].Sources[0].Models
		if len(models) != 1 || models[0].Code != code || models[0].Gender != 1 {
			t.Fatalf("race %d: single female code must be associated with its archive: %+v", id, models)
		}
	}
	if len(races[732].Sources) != 1 || len(races[732].Sources[0].Zones) != 0 || races[732].Sources[0].IsGlobal {
		t.Fatal("unlocated archives must not be treated as global")
	}
}
