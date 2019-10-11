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
		url    string
		fName  string
		lName  string
		expUrl string
	}{
		{"http://x.y", "jason", "bourne", "http://x.y?firstName=jason&lastName=bourne"},
		{"http://y.z", "Barry", "Allen", "http://y.z?firstName=Barry&lastName=Allen"},
	}

	for _, tt := range tests {
		t.Run(tt.fName+" "+tt.lName, func(t *testing.T) {
			u, err := url.Parse(tt.url)
			assert.Nil(err)
			assert.Equal(tt.expUrl, addNameParams(*u, tt.fName, tt.lName))
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
			jc.test = true
			joke, err := jc.JokeWithCustomName("", "")
			assert.Nil(err)
			assert.Equal(tt.joke, joke)
		})
	}
}
