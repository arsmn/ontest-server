package handler

import "github.com/go-chi/chi/v5"

func (a Handler) authRouter(r chi.Router) {
	r.Post("/signin", nil)
	r.Post("/signup", nil)
}
