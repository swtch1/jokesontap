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

var ErrUnmarshalingNamesAPI = errors.New("unable to unmarshal names API response body, possible rate limiting from name service")

// TODO: complete implementation of this interface
//// Cacher represents a cache.
//type Cacher interface {
//	Add(interface{}, interface{}) bool
//}

type Name struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
}

// NameClient can request random names from a names server.
type NameClient struct {
	// ApiUrl is the full URL of the names server from which we can request new names.
	ApiUrl url.URL
	// HttpClient is a http client which can be reused across multiple requests.
	HttpClient *http.Client
	//// Cache holds names that we may want to use if the external API is unavailable.
	//Cache Cacher

	// TODO: this can be removed when proper metrics are implemented
	// nameReqRate is the per-minute rate at which the name client is being requested.
	nameReqRate float64 // FIXME: testing
}

// NewNameClient creates a NameClient with default values where baseUrl is the API URL to query.
func NewNameClient(baseUrl url.URL) *NameClient {
	// TODO: the client values here _may_ be too detailed for the command line, but could be taken in thorough a more
	// TODO: detailed config file or env vars.  just would be better to be able to dynamically configure.
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
	req, err := http.NewRequest("GET", c.ApiUrl.String(), nil)
	if err != nil {
		return []Name{}, errors.Wrapf(err, "unable to create new http request with URL '%s'", c.ApiUrl.String())
	}
	req.Header.Set("Accept", "application/json")
	log.Tracef("getting names from name server")
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
		return []Name{}, errors.Wrap(err, fmt.Sprintf("%s", ErrUnmarshalingNamesAPI)) // FIXME: testing
	}
	return names, nil
}

//// CachedName gets a previously used name from the cache.
//func (c *NameClient) CachedName() string {
//	// TODO: implement a LRU cache so we have the option to pull from cached names if we run out of unique names
//	// TODO: this would also need to be taken as a flag, probably through a Cache-Control header in the request
//	return ""
//}

// TODO: this implementation is fairly specific to the problem it solves.  We could break this out to be a more general
// TODO: budget executor, but probably not necessary until we have to use this pattern in more tha one place.
// BudgetNameReq is a budgeted names API requester which will make no more requests than the
// external API will tolerate.
type BudgetNameReq struct {
	// requests keeps track of when requests were made.  The size of the array should be set to the maximum
	// number of requests (the budget) that can be made within the MinDiff time.
	//
	// When creating a budget where x operations can be run in y time, the size of this array should
	// be set as the x value.
	requests [6]time.Time
	// pos is the current position in requests, a tracker for getting the oldest names API request.
	pos int
	// MinDiff is the minimum amount of time between now and the oldest names API request allowed before
	// we are "over budget", after which we cannot make any more requests.
	//
	// When creating a budget where x operations can be run in y time, this should be set as the y value.
	MinDiff time.Duration

	// NameClient is used to request new names.
	NameClient NameRequester
	// NameChan is populated with the results of each names API request.
	NameChan chan Name
}

// RequestOften gets new names from the names API and pushes them to the names channel, as often as possible.
// If the timestamp of the oldest call is more than MinDiff then we wait until we expect to successfully
// make the next call.
func (b *BudgetNameReq) RequestOften() {
	var now time.Time
	var diff time.Duration

	for {
		now = time.Now()
		diff = now.Sub(b.oldestRequest())

		// wait until the time between now and the oldest request is within set bounds
		for diff < b.MinDiff {
			now = time.Now()
			diff = now.Sub(b.oldestRequest())
		}

		nameChanFull := len(b.NameChan) == cap(b.NameChan)
		if nameChanFull {
			log.Trace("names channel is full, skipping attempt to get new names")
			time.Sleep(time.Second * 1)
			continue
		}

		b.pushNamesFromAPI()
		b.updateRequestTime(now)
	}
}

// pushNamesFromAPI pushes a new batch of names from into the name channel.
func (b *BudgetNameReq) pushNamesFromAPI() {
	names, err := b.NameClient.Names()
	if err != nil {
		switch errors.Cause(err).(type) {
		case *json.SyntaxError:
			// TODO: this could be converted into a full circuit breaker pattern instead of this basic limit
			// TODO: something like exponential back-off could be more appropriate as we do not want to continue
			// TODO: to get limited and pay the 1 minute penalty
			// receiving a json SyntaxError could mean we are not able to unmarshal the response and are likely
			// being rate limited as the names API returns HTML when throttling
			time.Sleep(time.Second * 5)
		}
		log.WithError(err).Error("unable to get names from names client")
	}
	for _, name := range names {
		b.NameChan <- name
	}
}

func (b *BudgetNameReq) oldestRequest() time.Time {
	return b.requests[b.pos]
}

func (b *BudgetNameReq) updateRequestTime(t time.Time) {
	b.requests[b.pos] = t
	b.incPos()
}

// incPos increases the position counter, dropping back to 0 when the
// end of the requests tracking array is reached.
func (b *BudgetNameReq) incPos() {
	if b.pos >= len(b.requests)-1 {
		b.pos = 0
	} else {
		b.pos++
	}
}

type NameRequester interface {
	Names() ([]Name, error)
}
