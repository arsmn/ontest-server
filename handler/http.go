package handler

import (
	"net/http"

	"github.com/arsmn/ontest-server/app"
	"github.com/arsmn/ontest-server/module/context"
	"github.com/arsmn/ontest-server/module/oauth"
	"github.com/arsmn/ontest-server/module/xlog"
	"github.com/arsmn/ontest-server/settings"
	"github.com/rs/cors"

	"github.com/go-chi/chi/v5"
)

type (
	Map                 = context.Map
	Context             = context.Context
	HandleFunc          func(*Context) error
	Middleware          func(HandleFunc) HandleFunc
	handlerDependencies interface {
		app.Provider
		settings.Provider
		oauth.Provider
		xlog.Provider
	}
	Handler struct {
		dx      handlerDependencies
		handler http.Handler
		count   uint32
	}
)

func New(dx handlerDependencies) *Handler {
	h := new(Handler)

	h.dx = dx

	root := chi.NewRouter()

	root.Use(cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5004"},
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowedHeaders:   []string{"Origin", "Accept", "Content-Type", "X-Requested-With"},
		AllowCredentials: true,
	}).Handler)

	root.Route("/auth", h.authRouter)
	root.Route("/oauth", h.oauthHandler)
	root.Route("/account", h.accountRouter)

	h.handler = root

	return h
}

func (h *Handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h.handler.ServeHTTP(rw, r)
}

func (h *Handler) HandlersCount() uint32 {
	return h.count
}

func (h *Handler) TemplatesCount() uint32 {
	return 0
}

func (h *Handler) clown(fn HandleFunc, mws ...Middleware) http.HandlerFunc {
	h.count++

	handler := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		ctx := context.Acquire(rw, r)
		if err := chain(fn, mws)(ctx); err != nil {
			if catch := handleError(ctx, err); catch != nil {
				_ = ctx.SendStatus(http.StatusInternalServerError)
			}
		}
		context.Release(ctx)
	})

	return handler
}

func chain(endpoint HandleFunc, mws []Middleware) HandleFunc {
	if len(mws) == 0 {
		return endpoint
	}

	h := mws[len(mws)-1](endpoint)
	for i := len(mws) - 2; i >= 0; i-- {
		h = mws[i](h)
	}

	return h
}
