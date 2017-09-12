package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

var (
	stops = []os.Signal{syscall.SIGTERM, os.Interrupt, os.Kill}
	port  = flag.Int("port", 8080, "set serving port number")
)

func main() {
	flag.Parse()
	if err := serve(context.Background(), *port, stops...); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func serve(c context.Context, port int, stop ...os.Signal) error {

	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(time.Minute))

	// Routing
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Serving
	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: r,
	}

	eg, ctx := errgroup.WithContext(c)
	eg.Go(srv.ListenAndServe)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, stop...)

	catch := <-sig

	eg.Go(func() error {
		return srv.Shutdown(ctx)
	})

	if err := eg.Wait(); err != nil {
		return err
	}

	if catch == syscall.SIGTERM {
		return nil
	}

	return errors.Errorf("main: unexpected signal - signal = %+v", catch)
}
