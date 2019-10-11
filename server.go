package jokesontap

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type Server struct {
	Port int32
}

func (s *Server) ListenAndServe() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", HndlFunc)
	httpSrv := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.Port),
		Handler:      mux,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}
	return httpSrv.ListenAndServe()
}

func HndlFunc(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hellp")
}
