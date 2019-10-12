package jokesontap

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

var (
	ErrTooManyNameRequests = errors.New("too many name requests within the last minute")
	ErrOverNameApiBudget   = errors.New("name API request budget exceeded")
)

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

// BudgetNameReq represents a budgeted names API requester which will make no more requests than the
// external API will tolerate.
type BudgetNameReq struct {
	// reqTime keeps track of when queries were made.  The size of the array should be set to the maximum
	// number of requests that can be made within the minDiff time.
	reqTime [7]time.Time
	// pos is the current position in reqTime
	pos int
	// minDiff is the minimum amount of time between now and the current position in reqTime
	// before we are "over budget", after which we cannot make any more requests.
	//
	// When creating a budget where x operations can be run in y time, this should be set as the y value.
	minDiff time.Duration

	// NameClient is used to make queries to the names API.
	NameClient NameRequester
	// NameChan is populated with the results of each names API request.
	NameChan chan Name
}

// RequestOften gets new names from the names API and pushes them to the names channel, as often as possible.
func (b *BudgetNameReq) RequestOften() {
	// FIXME: not working yet
	now := time.Now()
	diff := now.Sub(b.reqTime[b.pos])
	for diff < b.minDiff {
		// not under budget yet, wait and try again
		fmt.Println("sleeping now...") // FIXME: testing
		time.Sleep(diff + time.Second)
	}
	names, err := b.NameClient.Names()
	if err != nil {
		// we do not expect to get errors under normal operation since we budget our calls, but
		// if the API or network is seeing issues this may happen
		log.WithError(err).Error("unable to get names from names client")
	}
	for _, name := range names {
		b.NameChan <- name
	}
	b.reqTime[b.pos] = now
	b.incPos()
}

// incPos increases the position counter, dropping back to 0 when the
// end of the reqTime tracking array is reached.
func (b *BudgetNameReq) incPos() {
	if b.pos >= len(b.reqTime)-1 {
		b.pos = 0
	} else {
		b.pos++
	}
}

type NameRequester interface {
	Names() ([]Name, error)
}
