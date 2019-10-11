package jokesontap

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type Cacher interface {
	Add(interface{}, interface{}) bool
}

type Name struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
}

type NameClient struct {
	// ApiUrl is the full URL of the names server from which we can request new names.
	ApiUrl url.URL
	// HttpClient is a http client which can be reused across multiple requests.
	HttpClient *http.Client
	// Cache holds names that we may want to use if the external API is unavailable.
	Cache Cacher
}

// NewNameClient creates a NameClient with default values where baseUrl is the API URL to query.
func NewNameClient(baseUrl url.URL) *NameClient {
	return &NameClient{
		ApiUrl: baseUrl,
		HttpClient: &http.Client{
			Timeout: 5 * time.Second,
			Transport: &http.Transport{
				DisableCompression: true,
				MaxIdleConns:       10,
				IdleConnTimeout:    30 * time.Second,
			},
		},
	}
}

// Names gets several names from the names API.
func (c *NameClient) Names() ([]Name, error) {
	req, err := http.NewRequest("GET", c.ApiUrl.String(), nil)
	if err != nil {
		log.WithError(err).Errorf("unable to create new http request with URL '%s'", c.ApiUrl.String())
	}
	req.Header.Set("Accept", "application/json")
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		log.WithError(err).Errorf("unable to get new name from '%s'", c.ApiUrl.String())
		return []Name{}, err
	}
	defer resp.Body.Close()

	var names []Name
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithError(err).Error("unable to read names API response body")
		return []Name{}, err
	}
	if err := json.Unmarshal(body, &names); err != nil {
		log.WithError(err).Error("unable to unmarshal names API response body")
		return []Name{}, err
	}
	return names, nil
}

// CachedName gets a previously used name from the cache.
func (c *NameClient) CachedName() string {
	// TODO: implement a LRU cache so we have the option to pull from cached names if we run out of unique names
	return ""
}
