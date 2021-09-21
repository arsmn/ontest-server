package handler

import (
	"net/http"

	"github.com/arsmn/ontest-server/module/generate"
	"github.com/arsmn/ontest-server/module/oauth"
	"github.com/arsmn/ontest-server/module/xlog"
	"github.com/go-chi/chi/v5"
	"golang.org/x/oauth2"
)

func (h *Handler) oauthHandler(r chi.Router) {
	r.Handle("/google", h.clown(h.google))
	r.Handle("/google/callback", h.clown(h.googleCallback))
	r.Handle("/github", h.clown(h.github))
	r.Handle("/github/callback", h.clown(h.githubCallback))
	r.Handle("/linkedin", h.clown(h.linkedin))
	r.Handle("/linkedin/callback", h.clown(h.linkedinCallback))
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
	s := h.dx.Settings().OAuth
	state := generate.RandomString(16, generate.AlphaNum)
	ctx.SetCookie(s.StateCookie,
		state, int(s.CookieLifespan.Seconds()), "/", h.dx.Settings().Serve.Domain,
		true, true, http.SameSiteDefaultMode)
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
	sett := h.dx.Settings()

	state, err := ctx.Cookie(sett.OAuth.StateCookie)
	ctx.RemoveCookie(sett.OAuth.StateCookie)

	if ctx.Request().FormValue("state") != state || err != nil {
		h.dx.Logger().Error("invalid state", xlog.String("state", state))
		return ctx.TemporaryRedirect("/")
	}

	h.dx.Logger().Info("", xlog.String("code", ctx.Request().FormValue("code")))
	token, err := cfg.Config().Exchange(ctx.Request().Context(), ctx.Request().FormValue("code"))
	if err != nil {
		h.dx.Logger().Error("exchange error", xlog.Err(err))
		return ctx.TemporaryRedirect("/")
	}

	req, err := cfg.FetchData(ctx.Request().Context(), token.AccessToken)
	if err != nil {
		return ctx.TemporaryRedirect("/")
	}

	sess, err := h.dx.App().OAuthIssueSession(ctx.Request().Context(), req)
	if err != nil {
		return err
	}

	ctx.SetSecureCookie(sett.Session.Cookie, sess.Token, int(sett.Session.Lifespan.Seconds()), "/", sett.Serve.Domain)

	return ctx.TemporaryRedirect(sett.Client.WebURL)
}
