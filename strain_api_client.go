// Package strainapiclient is a Go client module for calling
// The Strain API (learn more at https://http://strains.evanbusse.com).
package strainapiclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const baseURLHost string = "strainapi.evanbusse.com"
const baseURL string = "https://" + baseURLHost

// Client represents the interface a Client must implemenet
type Client interface {
	ListAllEffects()
	ListAllFlavors()
	ListAllStrains()
}

// DefaultClient is the default implementation of a Client for The Strain API
type DefaultClient struct {
	apiKey string
}

// NewDefaultClient creates a new DefaultClient with the apiKey passed in.
func NewDefaultClient(apiKey string) *DefaultClient {
	client := &DefaultClient{apiKey: apiKey}

	return client
}

// simpleHTTPGet is just a simple wrapper for getting basic
// byte slices from an HTTP GET call.
// It uses the base url of the API and appends the string
// passed in to the path (you must add a leading '/').
func (c *DefaultClient) simpleHTTPGet(restOfURLPath string) ([]byte, error) {
	req, err := http.NewRequest("GET", baseURL+"/"+c.apiKey+restOfURLPath, nil)
	req.Header.Set("Host", baseURLHost)
	req.Header.Set("User-Agent", "Tony fun")

	client := http.Client{
		Timeout: 0,
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("There was a problem connecting to the api: %s", err)
		return make([]byte, 0), err
	}

	defer resp.Body.Close()

	body, bodyErr := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return make([]byte, 0), fmt.Errorf("Status: %d - %s", resp.StatusCode, string(body))
	}

	if bodyErr != nil || err != nil {
		log.Printf("There was a problem reading the body of the response: %s", err)
		return make([]byte, 0), err
	}

	return body, nil
}

// CanConnect simply hits the root of the API with your API Key
// and makes sure it gets back the default response from the API.
func (c *DefaultClient) CanConnect() bool {
	// Expected response: Seems legit to me man...
	body, _ := c.simpleHTTPGet("")
	return string(body) == "Seems legit to me man..."
}

// Effect represents the effects that can be experienced when
// consuming a strain.
type Effect struct {
	Name string     `json:"effect"`
	Type EffectType `json:"type"`
}

// EffectType represents the possible types effects can be.
type EffectType string

// The valid vlaues of EffectType
const (
	// EffectTypePositive represents positive effects
	EffectTypePositive EffectType = "positive"
	// EffectTypeNegative represents negative effects
	EffectTypeNegative = "negative"
	// EffectTypeMedical represents possible medical-related effects
	EffectTypeMedical = "medical"
)

// ListAllEffects returns a slice of Effect elements that
// represents all effects that can be experienced.
func (c *DefaultClient) ListAllEffects() ([]Effect, error) {
	effects := make([]Effect, 0)

	allEffectsJSONBytes, err := c.simpleHTTPGet("/searchdata/effects")
	if err != nil {
		return effects, err
	}

	marshallErr := json.Unmarshal(allEffectsJSONBytes, &effects)
	return effects, marshallErr
}

// Flavor represents a componenet of strain flavor.
type Flavor string

// ListAllFlavors returns a slice of Flavor elements that
// represents all flavors of a strain.
func (c *DefaultClient) ListAllFlavors() ([]Flavor, error) {
	flavors := make([]Flavor, 0)

	allFlavorsJSONBytes, err := c.simpleHTTPGet("/searchdata/flavors")
	if err != nil {
		return flavors, err
	}

	marshallErr := json.Unmarshal(allFlavorsJSONBytes, &flavors)
	return flavors, marshallErr
}

// Race indicates the type of strain (Indica, Sativa, Hybrid)
type Race string

const (
	// RaceIndica represents a Race of a strain
	RaceIndica Race = "indica"
	// RaceSativa represents a Race of a strain
	RaceSativa = "sativa"
	// RaceHybrid represents a Race of a strain
	RaceHybrid = "hybrid"
)

// StrainSummary represents a high-level description of a Strain
type StrainSummary struct {
	Name        string `json:"name"`
	ID          int    `json:"id"`
	Description string `json:"desc"`
	Race        Race   `json:"race"`
}

// Strain represents a single strain of cannabis and its properites.
type Strain struct {
	Name        string                  `json:"name"`
	ID          int                     `json:"id"`
	Description string                  `json:"desc"`
	Race        Race                    `json:"race"`
	Flavors     []Flavor                `json:"flavors"`
	Effects     map[EffectType][]string `json:"effects"`
}

const strainSearchBasePath string = "/strains/search"

// ListAllStrainsResult represents the results of a strain search
type ListAllStrainsResult map[string]Strain

// ListAllStrains gets a ListAllStrainsResult of all strains
// (please use sparingly, it is expensive to run).
func (c *DefaultClient) ListAllStrains() (ListAllStrainsResult, error) {
	strainsResults := make(ListAllStrainsResult)

	findAllURL := strainSearchBasePath + "/all"
	strainsResultsJSONBytes, err := c.simpleHTTPGet(findAllURL)

	if err != nil {
		return strainsResults, err
	}

	marshallErr := json.Unmarshal(strainsResultsJSONBytes, &strainsResults)

	populateStrainNames(strainsResults)

	return strainsResults, marshallErr
}

// Set the name on each Strain to the name of the key
func populateStrainNames(strains ListAllStrainsResult) {
	for name, strain := range strains {
		strain.Name = name
		// Have to assign it back to the map to make it stick
		strains[name] = strain
	}
}

// StrainSearchResults is a slice of StrainSummary results from a search.
type StrainSearchResults []StrainSummary

// SearchStrainsByName gets a StrainSearchResult of all strains matching
// the name passed in.
func (c *DefaultClient) SearchStrainsByName(name string) (StrainSearchResults, error) {
	strainsResults := make(StrainSearchResults, 0)

	searchURL := strainSearchBasePath + "/name/" + name
	strainsResultsJSONBytes, err := c.simpleHTTPGet(searchURL)

	if err != nil {
		return strainsResults, err
	}

	marshallErr := json.Unmarshal(strainsResultsJSONBytes, &strainsResults)

	return strainsResults, marshallErr
}

// SearchStrainsByRace gets a StrainSearchResult of all strains matching
// the Race passed in.
func (c *DefaultClient) SearchStrainsByRace(race Race) (StrainSearchResults, error) {
	strainsResults := make(StrainSearchResults, 0)

	searchURL := strainSearchBasePath + "/race/" + string(race)
	strainsResultsJSONBytes, err := c.simpleHTTPGet(searchURL)

	if err != nil {
		return strainsResults, err
	}

	marshallErr := json.Unmarshal(strainsResultsJSONBytes, &strainsResults)

	return strainsResults, marshallErr
}
