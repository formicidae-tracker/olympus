package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

func JSONify(w http.ResponseWriter, obj interface{}) {
	data, err := json.Marshal(obj)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func RecoverWrap(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer func() {
			r := recover()
			if r != nil {
				switch t := r.(type) {
				case string:
					err = errors.New(t)
				case error:
					err = t
				default:
					err = errors.New("Unknown error")
				}
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}()
		h.ServeHTTP(w, r)
	})
}

func HTTPLogWrap(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[http]: %s %s from %s %s:", r.Method, r.RequestURI, r.RemoteAddr, r.UserAgent())
		h.ServeHTTP(w, r)
	})
}

type GracefulServer interface {
	Run() error
	Close() error
}

type gracefulServer struct {
	quit, done chan struct{}
	idle       chan error
	server     *http.Server
}

func NewGracefulServer(server *http.Server) GracefulServer {
	return &gracefulServer{
		server: server,
		quit:   make(chan struct{}),
		done:   make(chan struct{}),
		idle:   make(chan error),
	}
}

func (s *gracefulServer) Run() error {
	defer close(s.done)
	go func() {
		<-s.quit
		s.idle <- s.server.Shutdown(context.Background())
		close(s.idle)
	}()

	if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return <-s.idle
}

func (s *gracefulServer) Close() (err error) {
	defer func() {
		if recover() != nil {
			err = fmt.Errorf("already closed")
		}
		<-s.done
	}()
	close(s.quit)
	return nil
}
