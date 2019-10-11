package jokesontap

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

var ErrUnsuccessfulJokeQuery = errors.New("general error getting new joke")

type Joke struct {
	Type  string `json:"type"`
	Value struct {
		ID   int    `json:"id"`
		Joke string `json:"joke"`
	} `json:"value"`
}

func (j *Joke) Successful() bool {
	if j.Type == "success" {
		return true
	}
	return false
}

type JokeClient struct {
	// ApiUrl is the base URL of the jokes API to query
	ApiUrl url.URL
	// HttpClient is a http client which can be reused across multiple requests.
	HttpClient *http.Client
}

// NewJokeClient creates a JokeClient with default values where baseUrl is the API URL without any parameters.
func NewJokeClient(baseUrl url.URL) *JokeClient {
	// TODO: the client values here _may_ be too detailed for the command line, but could be taken in thorough a more
	// TODO: detailed config file or env vars
	return &JokeClient{
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

func (c *JokeClient) Joke() (string, error) {
	return c.jokeFromUrl(c.ApiUrl.String())
}

// JokeWithCustomName gets a new joke using the first and last name passed in.
func (c *JokeClient) JokeWithCustomName(fName, lName string) (string, error) {
	return c.jokeFromUrl(addNameParams(c.ApiUrl, fName, lName))
}

func (c JokeClient) jokeFromUrl(apiUrl string) (string, error) {
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return "", errors.Wrapf(err, "unable to create new http request with URL '%s'", apiUrl)
	}
	req.Header.Set("Accept", "application/json")
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return "", errors.Wrapf(err, "unable to get new joke from '%s'", apiUrl)
	}
	defer resp.Body.Close()

	var joke Joke
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "unable to read jokes API response body")
	}
	if err := json.Unmarshal(body, &joke); err != nil {
		return "", errors.Wrap(err, "unable to unmarshal jokes API response body")
	}
	if !joke.Successful() {
		return "", ErrUnsuccessfulJokeQuery
	}
	return joke.Value.Joke, nil
}

// addNameParams will add the first and last names as parameters to url.
func addNameParams(baseUrl url.URL, fName, lName string) string {
	params := url.Values{}
	params.Add("firstName", fName)
	params.Add("lastName", lName)
	baseUrl.RawQuery = params.Encode()
	return baseUrl.String()
}
