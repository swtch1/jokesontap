package jokesontap

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestGettingNamesFromNamesClient(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	tests := []struct {
		name     string
		resp     string
		expNames []Name
	}{
		{
			"single_name",
			`[{"name": "John", "surname": "Smith"}]`, []Name{
				{Name: "John", Surname: "Smith"},
			},
		},
		{
			"multiple_names",
			`[{"name": "John", "surname": "Smith"}, {"name": "Jay", "surname": "Grey"}]`, []Name{
				{Name: "John", Surname: "Smith"},
				{Name: "Jay", Surname: "Grey"},
			},
		},
		{
			"non_english",
			`[{"name": "Ασκάλαφος", "surname": "Γιάνναρης"}]`, []Name{
				{Name: "Ασκάλαφος", Surname: "Γιάνναρης"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, err := fmt.Fprintf(w, tt.resp)
				assert.Nil(err)
			}))
			defer ts.Close()

			u, err := url.Parse(ts.URL)
			assert.Nil(err)
			ns := NewNameClient(*u)
			names, err := ns.Names()
			assert.Nil(err)
			for _, n := range tt.expNames {
				assert.True(nameInNames(n, names))
			}
		})
	}
}

func TestNewNameServerEnforcesTimeout(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	tests := []struct {
		timeout time.Duration
	}{
		{time.Nanosecond * 1},
	}

	for _, tt := range tests {
		t.Run(string(tt.timeout), func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				// sleep so that the client will timeout
				time.Sleep(tt.timeout * 3)
				_, err := fmt.Fprintf(w, `[{"name": "John", "surname": "Smith"}]`)
				assert.Nil(err)
			}))
			defer ts.Close()

			u, err := url.Parse(ts.URL)
			assert.Nil(err)
			ns := NewNameClient(*u)
			// keep the timeout short so we force a timeout
			ns.HttpClient.Timeout = tt.timeout
			_, err = ns.Names()
			assert.Contains(err.Error(), "Client.Timeout exceeded")
		})
	}
}

func nameInNames(name Name, names []Name) bool {
	for _, n := range names {
		if n == name {
			return true
		}
	}
	return false
}

func TestInvalidNameResponse(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	tests := []struct {
		invalidResponse string
		expErrContains  string
	}{
		{`{"invalid"`, "unexpected end of JSON input"},
		{`{"invalid"`, "unable to unmarshal names API response"},
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
			nc := NewNameClient(*u)
			_, err = nc.Names()
			assert.Contains(err.Error(), tt.expErrContains)
		})
	}
}

func TestBudgetedNamesPositionNeverPanics(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	nr := BudgetNameReq{}
	for i := 0; i < len(nr.reqTime)*2; i++ {
		// basically just asserting that we aren't off by 1 which would eventually panic on an invalid index
		assert.NotPanics(nr.incPos)
	}
}

type MockNameClient struct {
	NamesMethodCalls int
}

func (c *MockNameClient) Names() ([]Name, error) {
	c.NamesMethodCalls++
	return []Name{
		{Name: "John", Surname: "Smith"},
		{Name: "Steve", Surname: "Wilson"},
	}, nil
}

func TestBudgetedNamesExecutesNoMoreOftenThanExpected(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	nameClient := &MockNameClient{}

	nchan := make(chan Name, 100)
	nr := BudgetNameReq{
		reqTime:    [7]time.Time{},
		NameClient: nameClient,
		NameChan:   nchan,
	}

	go nr.RequestOften()
	fmt.Println(nameClient.NamesMethodCalls)
	time.Sleep(time.Second * 2)
	fmt.Println(nameClient.NamesMethodCalls)
	time.Sleep(time.Second * 2)
	fmt.Println(nameClient.NamesMethodCalls)
	time.Sleep(time.Second * 2)
	fmt.Println(nameClient.NamesMethodCalls)
	time.Sleep(time.Second * 2)
	assert.True(false) // FIXME: start here
}
