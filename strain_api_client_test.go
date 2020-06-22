package strainapiclient

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func createTestDefaultClient(t *testing.T) *DefaultClient {
	apiKey, found := os.LookupEnv("STRAIN_API_KEY")

	if !found {
		t.Errorf("Problem getting STRAIN_API_KEY from environemnt")
	}

	return NewDefaultClient(apiKey)
}

func readAPIKeyFile(environment string) (string, error) {
	apiKeyBytes, err := ioutil.ReadFile("./api_key_" + environment)

	if err != nil {
		return "", err
	}

	return string(apiKeyBytes), nil
}

func commonFirstStrain() Strain {
	return Strain{
		Name:    "Afpak",
		ID:      1,
		Race:    RaceHybrid,
		Flavors: []Flavor{"Earthy", "Chemical", "Pine"},
		Effects: map[EffectType][]string{
			"positive": {"Relaxed", "Hungry", "Happy", "Sleepy"},
			"negative": {"Dizzy"},
			"medical":  {"Depression", "Insomnia", "Pain", "Stress", "Lack of Appetite"},
		},
	}
}

func commonFirstSearchStrainByNameResult() SearchStrainsByNameResult {
	return SearchStrainsByNameResult{
		Name:        "Afpak",
		ID:          1,
		Race:        RaceHybrid,
		Description: "Afpak, named for its direct Afghani and Pakistani landrace heritage, is a beautiful indica-dominant hybrid with light green and deep bluish purple leaves. The taste and aroma are floral with a touch of lemon, making the inhale light and smooth. Its effects start in the stomach by activating the appetite. There is also a potent relaxation that starts in the head and face, and gradually sinks down into the body. Enjoy this strain if you’re suffering from stress, mild physical discomfort, or having difficulty eating.",
	}
}

func commonFirstSearchStrainByRaceResult() SearchStrainsByRaceResult {
	return SearchStrainsByRaceResult{
		Name: "Afpak",
		ID:   1,
		Race: RaceHybrid,
	}
}

func commonFirstSearchStrainByEffectNameResult() SearchStrainsByEffectNameResult {
	return SearchStrainsByEffectNameResult{
		Name:       "Afpak",
		ID:         1,
		Race:       RaceHybrid,
		EffectName: "Happy",
	}
}

func commonFirstSearchStrainByFlavorResult() SearchStrainsByFlavorResult {
	return SearchStrainsByFlavorResult{
		Name:   "Afpak",
		ID:     1,
		Race:   RaceHybrid,
		Flavor: "Earthy",
	}
}
func TestConnect(t *testing.T) {
	expected := true
	if returnValue := createTestDefaultClient(t).CanConnect(); !returnValue {
		t.Error("Could not connect!", returnValue, expected)
	}
}

func TestListAllEffects(t *testing.T) {
	// These are purposefully dumb tests that should fail as new data gets added
	// but using for now for SOMETHING.
	// Bad things about it: 1) live HTTP calls, 2) hard-coded count results, 3) hard-coded first result
	expectedCount := 33
	expectedFirstName := "Relaxed"
	expectedFirstType := EffectTypePositive

	allEffects, err := createTestDefaultClient(t).ListAllEffects()

	if err != nil {
		t.Error("Failed trying to list all effects", err)
		return
	}

	actualCount := len(allEffects)

	if actualCount != expectedCount ||
		allEffects[0].Name != expectedFirstName || allEffects[0].Type != expectedFirstType {
		t.Errorf("Expected %d effects, where the first one is %s: %s; Actual: %d items", expectedCount, expectedFirstName, expectedFirstType, actualCount)
		return
	}
}

func TestListALlFlavors(t *testing.T) {
	// These are purposefully dumb tests that should fail as new data gets added
	// but using for now for SOMETHING.
	// Bad things about it: 1) live HTTP calls, 2) hard-coded count results, 3) hard-coded first result
	expectedCount := 48
	expectedFirstFlavor := Flavor("Earthy")

	allFlavors, err := createTestDefaultClient(t).ListAllFlavors()
	if err != nil {
		t.Error("Failed trying to list all flavors", err)
	}

	actualCount := len(allFlavors)

	if actualCount != expectedCount || expectedFirstFlavor != allFlavors[0] {
		t.Errorf("Expected %d flavors with %s as the first, got %v instead.", expectedCount, expectedFirstFlavor, allFlavors)
	}
}

func TestListAllStrains(t *testing.T) {
	// These are purposefully dumb tests that should fail as new data gets added
	// but using for now for SOMETHING.
	// Bad things about it: 1) live HTTP calls, 2) hard-coded count results, 3) hard-coded first result
	expectedCount := 1970
	expectedFirstStrain := commonFirstStrain()

	allStrains, err := createTestDefaultClient(t).ListAllStrains()
	if err != nil {
		t.Error("Failed trying to list all strains", err)
	}

	actualCount := len(allStrains)
	actualStrain := allStrains[expectedFirstStrain.Name]

	if actualCount != expectedCount || !cmp.Equal(expectedFirstStrain, actualStrain) {
		t.Errorf("Expected %d strains with %s as the first, got %d results of %v instead.", expectedCount, expectedFirstStrain.Name, actualCount, actualStrain)
	}
}

func TestSearchStrainsByName(t *testing.T) {
	// These are purposefully dumb tests that should fail as new data gets added
	// but using for now for SOMETHING.
	// Bad things about it: 1) live HTTP calls, 2) hard-coded first result
	expectedFirstStrainResult := commonFirstSearchStrainByNameResult()

	allStrains, err := createTestDefaultClient(t).SearchStrainsByName("Af")
	if err != nil {
		t.Error("Failed trying to search strains by name:", expectedFirstStrainResult.Name, err)
	}

	actualCount := len(allStrains)
	actualStrain := allStrains[0]

	// Sample data has a trailing 32 and 0xc2a0.  Trim before continuing
	actualStrain.Description = strings.TrimSpace(
		strings.Replace(
			actualStrain.Description, string([]byte{0xc2, 0xa0}), "", 0))

	if actualCount == 0 || !cmp.Equal(expectedFirstStrainResult, actualStrain) {
		t.Errorf("Expected at least one strain with %v as the first, got %d results of %v instead.", expectedFirstStrainResult, actualCount, actualStrain)
	}
}

func TestSearchStrainsByRace(t *testing.T) {
	// These are purposefully dumb tests that should fail as new data gets added
	// but using for now for SOMETHING.
	// Bad things about it: 1) live HTTP calls, 2) hard-coded first result
	expectedFirstStrainResult := commonFirstSearchStrainByRaceResult()

	allStrains, err := createTestDefaultClient(t).SearchStrainsByRace(RaceHybrid)
	if err != nil {
		t.Error("Failed trying to search strains by race:", expectedFirstStrainResult.Name, err)
	}

	actualCount := len(allStrains)
	actualStrain := allStrains[0]

	if actualCount == 0 || !cmp.Equal(expectedFirstStrainResult, actualStrain) {
		t.Errorf("Expected at least one strain with %v as the first, got %d results of %v instead.", expectedFirstStrainResult, actualCount, actualStrain)
	}
}

func TestSearchStrainsByEffectName(t *testing.T) {
	// These are purposefully dumb tests that should fail as new data gets added
	// but using for now for SOMETHING.
	// Bad things about it: 1) live HTTP calls, 2) hard-coded first result
	expectedFirstStrainResult := commonFirstSearchStrainByEffectNameResult()

	allStrains, err := createTestDefaultClient(t).SearchStrainsByEffectName("Happy")
	if err != nil {
		t.Error("Failed trying to search strains by effect name:", expectedFirstStrainResult.Name, err)
	}

	actualCount := len(allStrains)
	actualStrain := allStrains[0]

	if actualCount == 0 || !cmp.Equal(expectedFirstStrainResult, actualStrain) {
		t.Errorf("Expected at least one strain with %v as the first, got %d results of %v instead.", expectedFirstStrainResult, actualCount, actualStrain)
	}
}

func TestSearchStrainsByFlavor(t *testing.T) {
	// These are purposefully dumb tests that should fail as new data gets added
	// but using for now for SOMETHING.
	// Bad things about it: 1) live HTTP calls, 2) hard-coded first result
	expectedFirstStrainResult := commonFirstSearchStrainByFlavorResult()

	allStrains, err := createTestDefaultClient(t).SearchStrainsByFlavor("Earthy")
	if err != nil {
		t.Error("Failed trying to search strains by name:", expectedFirstStrainResult.Name, err)
	}

	actualCount := len(allStrains)
	actualStrain := allStrains[0]

	if actualCount == 0 || !cmp.Equal(expectedFirstStrainResult, actualStrain) {
		t.Errorf("Expected at least one strain with %v as the first, got %d results of %v instead.", expectedFirstStrainResult, actualCount, actualStrain)
	}
}
