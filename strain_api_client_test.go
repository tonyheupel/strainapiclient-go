package strainapiclient

import (
	"io/ioutil"
	"os"
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
	expectedFirstStrain := Strain{
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
