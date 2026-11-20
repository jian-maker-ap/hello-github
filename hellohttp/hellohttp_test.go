package hellohttp_test

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/jian-maker-ap/hello-github/hellohttp"
)

// wantContentType is the Content-Type of every hello world response.
const wantContentType = "text/plain; charset=utf-8"

// serve serves one request through h and returns the recorded response. The
// response body is closed when the test finishes.
func serve(t *testing.T, h http.Handler, method, target string) *http.Response {
	t.Helper()
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(method, target, nil))
	res := rec.Result()
	t.Cleanup(func() { res.Body.Close() })
	return res
}

func TestHandlerGet(t *testing.T) {
	for _, target := range []string{"/", "/?name=world"} {
		t.Run(target, func(t *testing.T) {
			res := serve(t, hellohttp.Handler(), http.MethodGet, target)

			if got, want := res.StatusCode, http.StatusOK; got != want {
				t.Fatalf("GET %s: status = %d, want %d", target, got, want)
			}
			if got, want := res.Header.Get("Content-Type"), wantContentType; got != want {
				t.Fatalf("GET %s: Content-Type = %q, want %q", target, got, want)
			}
			body, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("GET %s: read body: %v", target, err)
			}
			if got, want := string(body), hellohttp.Message; got != want {
				t.Fatalf("GET %s: body = %q, want %q", target, got, want)
			}
		})
	}
}

func TestHandlerHead(t *testing.T) {
	res := serve(t, hellohttp.Handler(), http.MethodHead, "/")

	if got, want := res.StatusCode, http.StatusOK; got != want {
		t.Fatalf("HEAD /: status = %d, want %d", got, want)
	}
	if got, want := res.Header.Get("Content-Type"), wantContentType; got != want {
		t.Fatalf("HEAD /: Content-Type = %q, want %q", got, want)
	}
}

func TestHandlerNotFound(t *testing.T) {
	for _, target := range []string{"/hello", "/hello/", "/index.html", "/api/v1"} {
		res := serve(t, hellohttp.Handler(), http.MethodGet, target)
		if got, want := res.StatusCode, http.StatusNotFound; got != want {
			t.Errorf("GET %s: status = %d, want %d", target, got, want)
		}
	}
}

func TestHandlerMethodNotAllowed(t *testing.T) {
	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete} {
		res := serve(t, hellohttp.Handler(), method, "/")
		if got, want := res.StatusCode, http.StatusMethodNotAllowed; got != want {
			t.Errorf("%s /: status = %d, want %d", method, got, want)
		}
		if got, want := res.Header.Get("Allow"), "GET, HEAD"; got != want {
			t.Errorf("%s /: Allow = %q, want %q", method, got, want)
		}
	}
}

func TestNewServer(t *testing.T) {
	srv := hellohttp.NewServer("127.0.0.1:9000")

	if got, want := srv.Addr, "127.0.0.1:9000"; got != want {
		t.Fatalf("Addr = %q, want %q", got, want)
	}
	if srv.Handler == nil {
		t.Fatal("Handler is nil")
	}
	timeouts := map[string]time.Duration{
		"ReadHeaderTimeout": srv.ReadHeaderTimeout,
		"ReadTimeout":       srv.ReadTimeout,
		"WriteTimeout":      srv.WriteTimeout,
		"IdleTimeout":       srv.IdleTimeout,
	}
	for name, timeout := range timeouts {
		if timeout <= 0 {
			t.Errorf("%s = %v, want a positive duration", name, timeout)
		}
	}
}

func TestServerServesRealRequests(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}

	srv := hellohttp.NewServer(ln.Addr().String())
	serveErr := make(chan error, 1)
	go func() { serveErr <- srv.Serve(ln) }()
	defer func() {
		if err := srv.Close(); err != nil {
			t.Errorf("Close: %v", err)
		}
		if err := <-serveErr; err != nil && !errors.Is(err, http.ErrServerClosed) {
			t.Errorf("Serve: %v", err)
		}
	}()

	base := "http://" + ln.Addr().String()
	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get(base + "/")
	if err != nil {
		t.Fatalf("GET /: %v", err)
	}
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Fatalf("GET /: read body: %v", err)
	}
	if got, want := resp.StatusCode, http.StatusOK; got != want {
		t.Fatalf("GET /: status = %d, want %d", got, want)
	}
	if got, want := string(body), hellohttp.Message; got != want {
		t.Fatalf("GET /: body = %q, want %q", got, want)
	}
	if got, want := resp.Header.Get("Content-Type"), wantContentType; got != want {
		t.Fatalf("GET /: Content-Type = %q, want %q", got, want)
	}

	resp, err = client.Head(base + "/")
	if err != nil {
		t.Fatalf("HEAD /: %v", err)
	}
	body, err = io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Fatalf("HEAD /: read body: %v", err)
	}
	if len(body) != 0 {
		t.Fatalf("HEAD /: body = %q, want no body", body)
	}
	if got, want := resp.StatusCode, http.StatusOK; got != want {
		t.Fatalf("HEAD /: status = %d, want %d", got, want)
	}

	resp, err = client.Post(base+"/", "text/plain", strings.NewReader("hello"))
	if err != nil {
		t.Fatalf("POST /: %v", err)
	}
	resp.Body.Close()
	if got, want := resp.StatusCode, http.StatusMethodNotAllowed; got != want {
		t.Fatalf("POST /: status = %d, want %d", got, want)
	}

	resp, err = client.Get(base + "/missing")
	if err != nil {
		t.Fatalf("GET /missing: %v", err)
	}
	resp.Body.Close()
	if got, want := resp.StatusCode, http.StatusNotFound; got != want {
		t.Fatalf("GET /missing: status = %d, want %d", got, want)
	}
}

func ExampleHandler() {
	srv := httptest.NewServer(hellohttp.Handler())
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(string(body))
	// Output: hello world
}
