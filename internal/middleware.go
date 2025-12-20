package internal

import (
	"bytes"
	"log"
	"net/http"
	"time"
)

type customResponse struct {
	http.ResponseWriter
	status int
	body   *bytes.Buffer
}

func newCustomResponse(w http.ResponseWriter) *customResponse {
	return &customResponse{ResponseWriter: w, body: &bytes.Buffer{}}
}

func (rw *customResponse) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

func (rw *customResponse) Write(data []byte) (int, error) {
	rw.body.Write(data)
	return rw.ResponseWriter.Write(data)
}

// Middlewares
func Logging(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := newCustomResponse(w)

		h.ServeHTTP(w, r)

		log.Printf("method=%s path=%s status=%d time=%v response=%s\n",
			r.Method, r.URL.Path, rw.status, time.Since(start), rw.body)
	})
}

func NameResponseServer(h http.Handler, nameServer string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", nameServer)
		h.ServeHTTP(w, r)
	})
}

// Basic Auth
func BasicAuth(uusername, password string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, p, ok := r.BasicAuth()

		if !ok || u != uusername || p != password {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			w.WriteHeader(401)
			w.Write([]byte("Unauthorised.\n"))
			return
		}
		next.ServeHTTP(w, r)
	})
}
