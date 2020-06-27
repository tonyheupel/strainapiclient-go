package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/tchype/strainapiclient-go"
)

func main() {
	const apiEnvironmentVariableName string = "STRAIN_API_KEY"
	apiKey, found := os.LookupEnv(apiEnvironmentVariableName)
	if !found {
		log.Fatalf("Did not find envrionment variable '%s'", apiEnvironmentVariableName)
	}

	client := strainapiclient.NewDefaultClient(apiKey)

	fmt.Println("Connecting to API...")
	canConnect := client.CanConnect()
	if !canConnect {
		log.Fatalf("Unable to connect to the API with apiKey: '%s'", apiKey)
	}

	const strainID int = 1
	fmt.Printf("Calling for all effects for Strain with ID %d\n", strainID)
	effectsByEffectType, err := client.GetStrainEffectsByStrainID(strainID)
	if err != nil {
		log.Fatalf("Problem calling for effects for Strain with ID %d: %s", strainID, err)
	}

	fmt.Println("Received the Effecs grouped by EffectType...")
	for effectType, effects := range effectsByEffectType {
		fmt.Printf("Type: %s\n", effectType)

		for _, effect := range effects {
			fmt.Printf("\t%s\n", effect.Name)
		}
	}

	fmt.Println("\nOverriding the HandleResourceRequestFunc with a function that always returns an error...")
	// Keep the original handler so we can set it back right after we use it
	originalResourceHandler := client.SetHandleResourceRequestFunc(alwaysReturnErrorRegardlessOfResourcePath)
	_, mockErr := client.ListAllStrains()
	// Don't forget to reset the handler back to the original
	_ = client.SetHandleResourceRequestFunc(originalResourceHandler)
	fmt.Println("Got error", mockErr)

	const mockStrainID int = 1234567890
	fmt.Println("\nOverriding the HandleResourceRequestFunc, but only for a specific GetFlavorsByStrainID call with a specific Strain ID:", mockStrainID, "...")

	originalResourceHandler = wrapDefaultResourceHandlerForFlavorsByStrainIDRequest(client, mockStrainID)
	defer client.SetHandleResourceRequestFunc(originalResourceHandler)

	// Since we are mocking the GetStrainFlavorsByStrainID call, we can
	// predict what it's going to be.
	expectedFlavorsForStrain := []strainapiclient.Flavor{"Snot", "Tar", "Unrealized Dreams"}
	flavorsForStrain, err := client.GetStrainFlavorsByStrainID(mockStrainID)
	if err != nil {
		log.Fatalf("Problem using GetStrainFlavorsByStrainID with the mock flavors hander: %s", err)
	}

	if !reflect.DeepEqual(flavorsForStrain, expectedFlavorsForStrain) {
		log.Fatalf("Expected to get back mock Flavors: %v\nbut got back: %v\n", expectedFlavorsForStrain, flavorsForStrain)
	}

	fmt.Printf("Got the expected mock flavors result back: %v\n", flavorsForStrain)

	const realStrainID int = 3
	fmt.Printf("\nRetrieving what is hopefully real and not mock flavors for Strain with ID %d ...\n", realStrainID)
	realStrainFlavors, err := client.GetStrainFlavorsByStrainID(realStrainID)

	fmt.Printf("See, it's only mocked for the one ID.  Here's Flavors for Strain ID of %d from real api calls: %v\n", realStrainID, realStrainFlavors)

	if len(realStrainFlavors) == 0 {
		return
	}

	firstFlavor := realStrainFlavors[0]
	fmt.Printf("\nAnd here is a different call to search by the first flavor that was returned from the previous call: %s...\n", firstFlavor)
	strains, err := client.SearchStrainsByFlavor(firstFlavor)
	if err != nil {
		log.Fatalf("There was a problem searching for strains by flavor: %s", err)
	}

	fmt.Printf("Found %d strains for flavor %s\n (which is the only flavor returned):", len(strains), firstFlavor)
	for index, strain := range strains {
		fmt.Printf("\tIndex: %d\t%v\n", index, strain)
	}
}

// Simple HandleResourceRequestFunc that always returns an same error, regardless of path
func alwaysReturnErrorRegardlessOfResourcePath(resourcePath string) ([]byte, error) {
	return make([]byte, 0), fmt.Errorf("Always returning an error; resource path: %s", resourcePath)
}

// Wrap the current HandleResourceRequestFunc with this one that returns a mock response for
// flavors by strain id with a specific strain id; the rest of the calls will use the previous HandleResourceRequestFunc implementation.
func wrapDefaultResourceHandlerForFlavorsByStrainIDRequest(client strainapiclient.Client, strainID int) strainapiclient.HandleResourceRequestFunc {

	originalFunc := client.SetHandleResourceRequestFunc(nil)

	var wrapperFuncForStrainFlavor strainapiclient.HandleResourceRequestFunc
	wrapperFuncForStrainFlavor = func(resourcePath string) ([]byte, error) {
		if strings.Contains(resourcePath, fmt.Sprintf("/strains/data/flavors/%d", strainID)) {
			returnValue := []byte("[\"Snot\", \"Tar\", \"Unrealized Dreams\"]")
			return returnValue, nil
		}

		return originalFunc(resourcePath)
	}

	_ = client.SetHandleResourceRequestFunc(wrapperFuncForStrainFlavor)

	return originalFunc
}
