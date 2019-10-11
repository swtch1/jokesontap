package jokesontap

import (
	log "github.com/sirupsen/logrus"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// don't produce normal logs during tests
	log.SetOutput(nopWriter{})
	os.Exit(m.Run())
}

type nopWriter struct{}

func (w nopWriter) Write(p []byte) (int, error) {
	return 0, nil
}
