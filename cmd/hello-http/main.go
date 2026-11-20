// Command hello-http serves the hello world HTTP service.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jian-maker-ap/hello-github/hellohttp"
)

// shutdownTimeout bounds how long a graceful shutdown waits for in-flight
// requests.
const shutdownTimeout = 5 * time.Second

func main() {
	logger := log.New(os.Stderr, "hello-http: ", 0)

	addr := flag.String("addr", ":8080", "TCP address to listen on")
	flag.Parse()

	ln, err := net.Listen("tcp", *addr)
	if err != nil {
		logger.Fatalf("listen: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := serve(ctx, ln, logger); err != nil {
		logger.Fatal(err)
	}
}

// serve serves the hello world service on ln until ctx is done, then shuts it
// down gracefully. It returns nil on a clean shutdown.
func serve(ctx context.Context, ln net.Listener, logger *log.Logger) error {
	srv := hellohttp.NewServer(ln.Addr().String())
	serveErr := make(chan error, 1)
	go func() { serveErr <- srv.Serve(ln) }()

	logger.Printf("listening on %s", ln.Addr())

	select {
	case err := <-serveErr:
		return err
	case <-ctx.Done():
	}

	logger.Printf("shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}
	if err := <-serveErr; err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}
