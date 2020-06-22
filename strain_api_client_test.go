package strainapiclient

import (
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func createTestDefaultClient(t *testing.T) *DefaultClient {
	apiKey, found := os.LookupEnv("STRAIN_API_KEY")

	if !found && t != nil {
		t.Errorf("Problem getting STRAIN_API_KEY from environemnt")
	}

	return NewDefaultClient(apiKey)
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
	var client Client = createTestDefaultClient(t)
	expectedCount := 33
	expectedFirstName := "Relaxed"
	expectedFirstType := EffectTypePositive

	allEffects, err := client.ListAllEffects()

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
	var client Client = createTestDefaultClient(t)
	expectedCount := 48
	expectedFirstFlavor := Flavor("Earthy")

	allFlavors, err := client.ListAllFlavors()
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
	var client Client = createTestDefaultClient(t)
	expectedCount := 1970
	expectedFirstStrain := commonFirstStrain()

	allStrains, err := client.ListAllStrains()
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
	var client Client = createTestDefaultClient(t)
	expectedFirstStrainResult := commonFirstSearchStrainByNameResult()

	allStrains, err := client.SearchStrainsByName("Af")
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
	var client Client = createTestDefaultClient(t)
	expectedFirstStrainResult := commonFirstSearchStrainByRaceResult()

	allStrains, err := client.SearchStrainsByRace(RaceHybrid)
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
	var client Client = createTestDefaultClient(t)
	expectedFirstStrainResult := commonFirstSearchStrainByEffectNameResult()

	allStrains, err := client.SearchStrainsByEffectName("Happy")
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
	var client Client = createTestDefaultClient(t)
	expectedFirstStrainResult := commonFirstSearchStrainByFlavorResult()

	allStrains, err := client.SearchStrainsByFlavor("Earthy")
	if err != nil {
		t.Error("Failed trying to search strains by name:", expectedFirstStrainResult.Name, err)
	}

	actualCount := len(allStrains)
	actualStrain := allStrains[0]

	if actualCount == 0 || !cmp.Equal(expectedFirstStrainResult, actualStrain) {
		t.Errorf("Expected at least one strain with %v as the first, got %d results of %v instead.", expectedFirstStrainResult, actualCount, actualStrain)
	}
}

func TestGetStrainDescriptionByStrainID(t *testing.T) {
	// These are purposefully dumb tests that should fail as new data gets added
	// but using for now for SOMETHING.
	// Bad things about it: 1) live HTTP calls, 2) hard-coded count results, 3) hard-coded first result
	var client Client = createTestDefaultClient(t)
	expectedDescription := commonFirstSearchStrainByNameResult().Description

	strainID := 1
	actualDescription, err := client.GetStrainDescriptionByStrainID(strainID)

	if err != nil {
		t.Errorf("Failed trying to get description for strain with ID of %d: %s", strainID, err)
	}

	// Sample has extra non-printable UTF-8 character c2a0
	actualDescription = strings.TrimSpace(strings.ReplaceAll(actualDescription, string([]byte{0xc2, 0xa0}), ""))

	if actualDescription != expectedDescription {
		t.Errorf("For strain with ID %d, expected description '%s' but got '%s'", strainID, expectedDescription, actualDescription)
	}
}

func TestGetStrainFavorsByStrainID(t *testing.T) {
	// These are purposefully dumb tests that should fail as new data gets added
	// but using for now for SOMETHING.
	// Bad things about it: 1) live HTTP calls, 2) hard-coded count results, 3) hard-coded first result
	var client Client = createTestDefaultClient(t)
	expectedFlavors := []Flavor{"Earthy", "Chemical", "Pine"}

	strainID := 1
	actualFlavors, err := client.GetStrainFavorsByStrainID(strainID)

	if err != nil {
		t.Errorf("Failed trying to get flavors for strain with ID of %d: %s", strainID, err)
	}

	if !reflect.DeepEqual(actualFlavors, expectedFlavors) {
		t.Errorf("For strain with ID %d, expected '%v' but got '%v'", strainID, expectedFlavors, actualFlavors)
	}
}

func TestGetStrainEffectsByStrainID(t *testing.T) {
	// These are purposefully dumb tests that should fail as new data gets added
	// but using for now for SOMETHING.
	// Bad things about it: 1) live HTTP calls, 2) hard-coded count results, 3) hard-coded first result
	var client Client = createTestDefaultClient(t)
	expectedPositiveEffects := []Effect{
		{Name: "Relaxed", Type: EffectTypePositive},
		{Name: "Happy", Type: EffectTypePositive},
		{Name: "Hungry", Type: EffectTypePositive},
		{Name: "Sleepy", Type: EffectTypePositive},
	}
	expectedNegativeEffects := []Effect{
		{Name: "Dizzy", Type: EffectTypeNegative},
	}

	expectedMedicalEffects := []Effect{
		{Name: "Depression", Type: EffectTypeMedical},
		{Name: "Stress", Type: EffectTypeMedical},
		{Name: "Lack of Appetite", Type: EffectTypeMedical},
		{Name: "Insomnia", Type: EffectTypeMedical},
		{Name: "Pain", Type: EffectTypeMedical},
	}

	expectedEffects := make(EffectsByEffectType)
	expectedEffects[EffectTypePositive] = expectedPositiveEffects
	expectedEffects[EffectTypeNegative] = expectedNegativeEffects
	expectedEffects[EffectTypeMedical] = expectedMedicalEffects

	strainID := 1
	actualEffects, err := client.GetStrainEffectsByStrainID(strainID)

	if err != nil {
		t.Errorf("Failed trying to get effects for strain with ID of %d: %s", strainID, err)
	}

	if !reflect.DeepEqual(actualEffects, expectedEffects) {
		t.Errorf("For strain with ID %d, expected effects '%v' but got '%v'", strainID, expectedEffects, actualEffects)
	}
}
