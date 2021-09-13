package daemon

import (
	"context"
	"net/http"
	"sync"

	"github.com/arsmn/ontest/driver"
	"github.com/arsmn/ontest/handler"
	"github.com/arsmn/ontest/settings"
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

	if s.StartupMessageEnabled() {
		startupMessage(server.Addr, false, h.HandlersCount(), h.TemplatesCount(), settings.ConfigFileUsed())
	}

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
