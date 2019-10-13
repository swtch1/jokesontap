package jokesontap

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

var (
	ErrNamesChanUninitialized = errors.New("the server's names channel is uninitialized, please submit an issue")
	ErrNoNamesAvailable       = errors.New("the server has no names to provide")
)

type Server struct {
	// Port is the port where the server will listen.
	Port int32
	// JokeClient requests new jokes using a customized name, if given.
	JokeClient *JokeClient
	// Names is a buffered channel where names will be retrieved from.  The server expects for this
	// to be populated ahead of time by another thread.  We are basically using this as a queue, but the
	// implementation is more simple and more easily supports handling timeouts.
	Names chan Name
}

func (s *Server) ListenAndServe() error {
	if s.Names == nil {
		return ErrNamesChanUninitialized
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.GetCustomJoke)
	httpSrv := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.Port),
		Handler:      mux,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}
	return httpSrv.ListenAndServe()
}

func (s *Server) GetCustomJoke(w http.ResponseWriter, req *http.Request) {
	select {
	case name := <-s.Names:
		joke, err := s.JokeClient.JokeWithCustomName(name.Name, name.Surname)
		if err != nil {
			log.WithError(err).Error("failed to get joke with custom name")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err, "\n")
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, joke, "\n")
	case <-time.After(time.Second * 5):
		log.WithError(ErrNoNamesAvailable).Error("timeout getting name")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, ErrNoNamesAvailable, "\n")
	}
}
