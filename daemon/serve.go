package daemon

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/arsmn/ontest/driver"
	"github.com/arsmn/ontest/handler"
	"github.com/ory/graceful"
)

func ServePublic(ctx context.Context, r driver.Registry, wg *sync.WaitGroup, args []string) {
	defer wg.Done()

	s := r.Settings()
	l := r.Logger()
	h := handler.New(r)

	var handler http.Handler = h

	server := graceful.WithDefaults(&http.Server{
		Addr:    s.PublicListenOn(),
		Handler: handler,
	})

	l.Info(fmt.Sprintf("Starting the public httpd on: %s", server.Addr))
	if err := graceful.Graceful(server.ListenAndServe, server.Shutdown); err != nil {
		l.Fatal("Failed to gracefully shutdown public httpd")
	}

	l.Info("Public httpd was shutdown gracefully")
}

func ServeAll(ctx context.Context, d driver.Registry) func(args []string) {
	return func(args []string) {
		var wg sync.WaitGroup
		wg.Add(1)
		go ServePublic(ctx, d, &wg, args)
		wg.Wait()
	}
}
