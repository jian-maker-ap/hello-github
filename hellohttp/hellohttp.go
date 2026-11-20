// Package hellohttp implements a hello world HTTP service.
//
// Handler returns the service's http.Handler and NewServer returns a
// ready-to-run *http.Server that serves it:
//
//	srv := hellohttp.NewServer(":8080")
//	log.Fatal(srv.ListenAndServe())
//
// The command in cmd/hello-http runs the service from the command line.
package hellohttp

import (
	"net/http"
	"time"
)

// Message is the greeting returned by the hello world service.
const Message = "hello world\n"

// Handler returns the http.Handler of the hello world service.
//
// The handler answers GET and HEAD requests for / with [Message]. Requests
// for any other path get a 404 Not Found response, and requests with any
// other method get a 405 Method Not Allowed response.
func Handler() http.Handler {
	return http.HandlerFunc(serveHelloWorld)
}

// NewServer returns an *http.Server that serves [Handler] on addr.
//
// The server has request timeouts set, so slow or stalled clients cannot hold
// a connection open indefinitely.
func NewServer(addr string) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           Handler(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
}

// serveHelloWorld serves the hello world message at /.
func serveHelloWorld(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		w.Header().Set("Allow", "GET, HEAD")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte(Message))
}
