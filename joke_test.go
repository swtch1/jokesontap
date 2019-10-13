package jokesontap

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestEncodingUrlParameters(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	tests := []struct {
		url      string
		fName    string
		lName    string
		category string
		expUrl   string
	}{
		{"http://x.y", "jason", "bourne", "nerdy", "http://x.y?firstName=jason&lastName=bourne&limitTo=%5Bnerdy%5D"},
		{"http://y.z", "Barry", "Allen", "nerdy", "http://y.z?firstName=Barry&lastName=Allen&limitTo=%5Bnerdy%5D"},
	}

	for _, tt := range tests {
		t.Run(tt.fName+" "+tt.lName, func(t *testing.T) {
			u, err := url.Parse(tt.url)
			assert.Nil(err)
			assert.Equal(tt.expUrl, addParams(*u, tt.fName, tt.lName, "nerdy"))
		})
	}
}

func TestGettingJokesFromJokesClient(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	tests := []struct {
		joke string
	}{
		{"John Smith had a joke."},
		{"Steve wilson had another one."},
	}

	for _, tt := range tests {
		t.Run(tt.joke, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, err := fmt.Fprintf(w, `{"type": "success", "value": { "joke": "%s" }}`, tt.joke)
				assert.Nil(err)
			}))
			defer ts.Close()

			u, err := url.Parse(ts.URL)
			assert.Nil(err)
			jc := NewJokeClient(*u)

			// setting testing explicitly here ensures we can avoid transforming the URL, maintaining the raw
			// URL from the test server
			joke, err := jc.Joke()
			assert.Nil(err)
			assert.Equal(tt.joke, joke)
		})
	}
}

func TestInvalidJokesResponse(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	tests := []struct {
		invalidResponse string
		expErrContains  string
	}{
		{`{"invalid"`, "unexpected end of JSON input"},
		{`{"invalid"`, "unable to unmarshal jokes API response"},
	}

	for _, tt := range tests {
		t.Run(tt.expErrContains, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, err := fmt.Fprintf(w, tt.invalidResponse)
				assert.Nil(err)
			}))
			defer ts.Close()

			u, err := url.Parse(ts.URL)
			assert.Nil(err)
			jc := NewJokeClient(*u)
			_, err = jc.JokeWithCustomName("john", "smith")
			assert.Contains(err.Error(), tt.expErrContains)
		})
	}
}
