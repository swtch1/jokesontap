package jokesontap

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestStartingServerWithNilNamesChanErrors(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	srv := Server{Port: 5000}
	assert.Equal(ErrNamesChanUninitialized, srv.ListenAndServe())
}

func TestServer(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	tests := []struct {
		name string
		joke string
	}{
		{"basic_response", "Chuck Norris wants to be Bill Murray"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, err := fmt.Fprintf(w, `{"type": "success", "value": { "joke": "%s"}}`, tt.joke)
				assert.Nil(err)
			}))
			defer ts.Close()

			jokeUrl, err := url.Parse(ts.URL)
			assert.Nil(err)
			jokeClient := NewJokeClient(*jokeUrl)
			nameChan := make(chan Name, 10)
			mockName := Name{
				Name:    "Bill",
				Surname: "Murray",
			}
			nameChan <- mockName

			srv := Server{
				Port:       0,
				JokeClient: jokeClient,
				Names:      nameChan,
			}

			req := httptest.NewRequest("GET", "http://doesnt.matter", nil)
			w := httptest.NewRecorder()
			srv.GetCustomJoke(w, req)
			body, err := ioutil.ReadAll(w.Result().Body)
			assert.Nil(err)
			// we expect a 200 response and the joke to be written to the response writer
			assert.Equal(http.StatusOK, w.Result().StatusCode)
			assert.Equal(tt.joke+"\n", string(body))
		})
	}
}
