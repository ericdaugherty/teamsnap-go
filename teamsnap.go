package teamsnap

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

const rootURL = "https://api.teamsnap.com/v3/"

// TeamSnap is the main object used to interact with the TeamSnap API
type TeamSnap struct {
	AuthToken string
	Version   string
	RootLinks []Link
}

// Response is the response provided from the root API Query
type Response struct {
	Collection Collection `json:"collection"`
}

// Collection of URLs that the API supports
type Collection struct {
	Version string `json:"version"`
	Href    string `json:"href"`
	Rel     string `json:"rel"`
	Links   []Link `json:"links"`
	Items   []Item `json:"items"`
}

// Link provides the rel and href for an API call
type Link struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
}

// Item contains Data and Links for a specific result
type Item struct {
	Href  string `json:"href"`
	Data  []Data `json:"data"`
	Links []Link `json:"links"`
}

// Data is a name/value pair
type Data struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

// Initialize should be called before any other methods.  Queries the API to verify
// it is reachable and determine hrefs for rel calls.
func (teamSnap *TeamSnap) Initialize() error {
	result, err := teamSnap.query(rootURL)
	if err != nil {
		return err
	}
	teamSnap.Version = result.Collection.Version
	teamSnap.RootLinks = result.Collection.Links
	return nil
}

// FetchRoot queries the specified root rel and returns the result.
func (teamSnap *TeamSnap) FetchRoot(rel string) (*Response, error) {
	return teamSnap.Fetch(rel, teamSnap.RootLinks)
}

// Fetch queries the specified rel given the provided links and returns the result.
func (teamSnap *TeamSnap) Fetch(rel string, links []Link) (*Response, error) {
	href, err := teamSnap.findHref(rel, links)
	if err != nil {
		return nil, err
	}
	result, err := teamSnap.query(href)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// DataValue looks for a data item that matches the specified name and returns the value.
// Value type is unknown and returned as an interface{}
// An error is returned if no match is found.
func (item *Item) DataValue(name string) (interface{}, error) {
	for _, data := range item.Data {
		if data.Name == name {
			return data.Value, nil
		}
	}
	return "", errors.New("No match found")
}

// DataValueString looks for a data item that matches the specified name and returns the value.
// If string, value is returned, otherwise value is converted to a string.
// An error is returned if no match is found.
func (item *Item) DataValueString(name string) (string, error) {
	value, err := item.DataValue(name)
	if err != nil {
		return "", err
	}

	switch v := value.(type) {
	case string:
		return v, nil
	case bool:
		return strconv.FormatBool(v), nil
	case float64:
		return strconv.FormatFloat(v, 'f', 0, 32), nil
	case nil:
		return "", nil
	default:
		return "", errors.New("Unable to covert Data Element " + name + " to a string.")
	}
}

// DataValueInt returns the value if it can be converted to an int, othrewise an error is returned
func (item *Item) DataValueInt(name string) (int64, error) {
	value, err := item.DataValue(name)
	if err != nil {
		return 0, err
	}

	switch v := value.(type) {
	case string:
		return strconv.ParseInt(v, 10, 32)
	case float64:
		return int64(v), nil
	default:
		return 0, errors.New("Unable to convert the value into an int.")
	}
}

// DataValueTime returns the value can be converted to a time.Time, othrewise an error is returned
func (item *Item) DataValueTime(name string) (time.Time, error) {
	value, err := item.DataValue(name)
	if err != nil {
		return time.Now(), err
	}

	switch v := value.(type) {
	case string:
		return time.Parse(time.RFC3339, v)
	default:
		return time.Now(), errors.New("Unable to convert the value into a time.")
	}
}

func (teamSnap *TeamSnap) findHref(rel string, links []Link) (string, error) {
	for _, link := range links {
		if link.Rel == rel {
			return link.Href, nil
		}
	}

	return "", errors.New("No href found for rel: " + rel)
}

func (teamSnap *TeamSnap) query(href string) (*Response, error) {

	client := &http.Client{}
	req, err := http.NewRequest("GET", href, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+teamSnap.AuthToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("\n\n\n\n", "Error!", err.Error())
		return nil, err
	}

	return &response, nil
}
