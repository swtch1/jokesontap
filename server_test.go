package jokesontap

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStartingServerWithNilNamesChanErrors(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	srv := Server{Port: 5000}
	assert.Equal(ErrNamesChanUninitialized, srv.ListenAndServe())
}
