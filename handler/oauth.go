package handler

import (
	"net/http"

	"github.com/arsmn/ontest-server/module/generate"
	"github.com/arsmn/ontest-server/module/oauth"
	"github.com/arsmn/ontest-server/transport"
	"github.com/go-chi/chi/v5"
	"golang.org/x/oauth2"
)

func (h *Handler) oauthHandler(r chi.Router) {
	r.Get("/google", h.clown(h.google))
	r.Get("/google/callback", h.clown(h.googleCallback))
	r.Get("/github", h.clown(h.github))
	r.Get("/github/callback", h.clown(h.githubCallback))
	r.Get("/linkedin", h.clown(h.linkedin))
	r.Get("/linkedin/callback", h.clown(h.linkedinCallback))
}

func (h *Handler) google(ctx *Context) error {
	return h.handleOAuth(ctx, h.dx.OAuther(oauth.GoogleType).Config())
}

func (h *Handler) github(ctx *Context) error {
	return h.handleOAuth(ctx, h.dx.OAuther(oauth.GitHubType).Config())
}

func (h *Handler) linkedin(ctx *Context) error {
	return h.handleOAuth(ctx, h.dx.OAuther(oauth.LinkedInType).Config())
}

func (h *Handler) handleOAuth(ctx *Context, cfg *oauth2.Config) error {
	s := h.dx.Settings().OAuth()
	state := generate.RandomString(16, generate.AlphaNum)
	ctx.SetCookie(s.StateCookie, state, int(s.CookieLifespan.Seconds()), "/", h.dx.Settings().Domain(), false, false, http.SameSiteLaxMode)
	url := cfg.AuthCodeURL(state, oauth2.AccessTypeOffline)
	return ctx.TemporaryRedirect(url)
}

func (h *Handler) googleCallback(ctx *Context) error {
	return h.handleCallback(ctx, h.dx.OAuther(oauth.GoogleType))
}

func (h *Handler) githubCallback(ctx *Context) error {
	return h.handleCallback(ctx, h.dx.OAuther(oauth.GitHubType))
}

func (h *Handler) linkedinCallback(ctx *Context) error {
	return h.handleCallback(ctx, h.dx.OAuther(oauth.LinkedInType))
}

func (h *Handler) handleCallback(ctx *Context, cfg oauth.OAuther) error {
	state, err := ctx.Cookie(h.dx.Settings().OAuth().StateCookie)
	if ctx.Request().FormValue("state") != state || err != nil {
		return ctx.TemporaryRedirect("/")
	}

	sess, err := h.dx.App().OAuthIssueSession(ctx.Request().Context(), &transport.OAuthSignRequest{})
	if err != nil {
		return err
	}

	s := h.dx.Settings().Session()
	ctx.SetSecureCookie(s.Cookie, sess.Token, int(s.Lifespan.Seconds()), "/", h.dx.Settings().Domain())

	return ctx.OK(success)
}
