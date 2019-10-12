package jokesontap

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

var ErrTooManyNameRequests = errors.New("too many name requests within the last minute")

// TODO: complete implementation of this interface
// Cacher represents a cache.
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
	// TODO: the client values here _may_ be too detailed for the command line, but could be taken in thorough a more
	// TODO: detailed config file or env vars
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

// Names gets several names from the names API.  Names is intelligent as it relates to the restrictions
// of the name API and will short circuit if too many requests are made.  If Names is called more often
// than the API will allow an ErrTooManyNameRequests error will be returned.
func (c *NameClient) Names() ([]Name, error) {
	// TODO: implemente requests budget
	req, err := http.NewRequest("GET", c.ApiUrl.String(), nil)
	if err != nil {
		return []Name{}, errors.Wrapf(err, "unable to create new http request with URL '%s'", c.ApiUrl.String())
	}
	req.Header.Set("Accept", "application/json")
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return []Name{}, errors.Wrapf(err, "unable to get new name from '%s'", c.ApiUrl.String())
	}
	defer resp.Body.Close()

	var names []Name
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []Name{}, errors.Wrapf(err, "unable to read names API response body")
	}
	if err := json.Unmarshal(body, &names); err != nil {
		return []Name{}, errors.Wrap(err, "unable to unmarshal names API response body")
	}
	return names, nil
}

// CachedName gets a previously used name from the cache.
func (c *NameClient) CachedName() string {
	// TODO: implement a LRU cache so we have the option to pull from cached names if we run out of unique names
	// TODO: this would also need to be taken as a flag, probably through a Cache-Control header in the request
	return ""
}
