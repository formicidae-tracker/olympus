package main

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
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

func Golangify[T any](r *http.Request) (*T, error) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		return nil, fmt.Errorf("invalid Content-Type %s", contentType)
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read request body: %s", err)
	}

	res := new(T)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, fmt.Errorf("could not parse JSON: %s", err)
	}

	return res, nil
}

func RecoverWrap(logger *log.Logger) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
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
					logger.Printf("panic: %s", err)

					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}()
			h.ServeHTTP(w, r)
		})
	}
}

type loggingResponseWriter struct {
	http.ResponseWriter
	status int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (w *loggingResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func HTTPLogWrap(logger *log.Logger) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lw := newLoggingResponseWriter(w)
			h.ServeHTTP(lw, r)
			logger.Printf("%s %s from %s %s: %d",
				r.Method, r.RequestURI, r.RemoteAddr, r.UserAgent(), lw.status)
		})
	}
}

func EnableCORS(origins []string) func(http.Handler) http.Handler {
	allowed := make(map[string]bool)
	for _, origin := range origins {
		allowed[origin] = true
	}

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			//concurrent read access to map is ok.
			if allowed[origin] == true {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}
			h.ServeHTTP(w, r)
		})
	}
}

func CacheControl(maxAge time.Duration) func(http.Handler) http.Handler {
	control := fmt.Sprintf("public max-age=%d", int(maxAge.Seconds()))
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", control)
			h.ServeHTTP(w, r)
		})
	}
}

type CSRFHandler struct {
	secret []byte
	logger *log.Logger
}

func NewCSRFHandler(secret []byte) (*CSRFHandler, error) {
	if len(secret) == 0 {
		return nil, errors.New("missing server secret")
	}
	return &CSRFHandler{
		secret: secret,
		logger: log.New(os.Stderr, "[csrf]: ", log.LstdFlags),
	}, nil
}

func (h *CSRFHandler) setCSRFCookie(w http.ResponseWriter) error {
	nonceBytes := make([]byte, 128)
	_, err := rand.Read(nonceBytes)
	if err != nil {
		return err
	}

	messageMac := hmac.New(sha256.New, h.secret).Sum(nonceBytes)

	token := fmt.Sprintf("%s.%s", base64.URLEncoding.EncodeToString(nonceBytes), base64.URLEncoding.EncodeToString(messageMac))

	http.SetCookie(w, &http.Cookie{
		Name:  "XSRF-TOKEN",
		Value: token,

		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
	return nil
}

func (h *CSRFHandler) checkXSRFToken(r *http.Request) error {
	cookie, err := r.Cookie("XSRF-TOKEN")
	if err != nil {
		return err
	}
	token := r.Header.Get("X-XSRF-TOKEN")
	if cookie.Value != token {
		return errors.New("cookie and token doesn't match")
	}

	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return fmt.Errorf("invalid token format: %s", token)
	}

	nonce, err := base64.URLEncoding.DecodeString(parts[0])
	if err != nil {
		return fmt.Errorf("could not decode nonce: %s", err)
	}

	nonceMAC, err := base64.URLEncoding.DecodeString(parts[1])
	if err != nil {
		return fmt.Errorf("could not decode HMAC: %s", err)
	}

	expectedMAC := hmac.New(sha256.New, h.secret).Sum(nonce)
	if hmac.Equal(nonceMAC, expectedMAC) == false {
		return fmt.Errorf("token HMAC is invalid")
	}

	return nil
}

func (hh *CSRFHandler) SetCSRFCookie(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := hh.setCSRFCookie(w); err != nil {
			hh.logger.Printf("could not generate token: %s", err)
			http.Error(w, "could not generate CSRF Token", http.StatusInternalServerError)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func (hh *CSRFHandler) CheckCSRFCookie(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := hh.checkXSRFToken(r); err != nil {
			hh.logger.Printf("invalid token on %s %s from %s: %s", r.Method, r.RequestURI, r.RemoteAddr, err)
			http.Error(w, "invalid CSRF token", http.StatusBadRequest)
			return
		}

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
