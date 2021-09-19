package context

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/arsmn/ontest-server/module/httplib"
)

type Map map[string]interface{}

func (ctx *Context) Json(status int, data interface{}) error {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(data); err != nil {
		return err
	}

	ctx.response.Header().Set("Content-Type", "application/json; charset=utf-8")
	ctx.response.WriteHeader(status)
	_, err := ctx.response.Write(buf.Bytes())
	return err
}

func (ctx *Context) String(status int, data string) error {
	ctx.response.Header().Set("Content-Type", "text/plain; charset=utf-8")
	ctx.response.WriteHeader(status)
	_, err := fmt.Fprintln(ctx.response, data)
	return err
}

func (ctx *Context) SendStatus(status int) error {
	ctx.response.Header().Set("X-Content-Type-Options", "nosniff")
	return ctx.String(status, http.StatusText(status))
}

func (ctx *Context) OK(data interface{}) error {
	return ctx.Json(http.StatusOK, data)
}

func (ctx *Context) Created(location string, data interface{}) error {
	ctx.response.Header().Set("Location", location)
	return ctx.Json(http.StatusCreated, data)
}

func (ctx *Context) SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool, sameSite http.SameSite) {
	if path == "" {
		path = "/"
	}
	host, _ := httplib.ParseAddr(domain)

	http.SetCookie(ctx.response, &http.Cookie{
		Name:     name,
		Value:    url.QueryEscape(value),
		MaxAge:   maxAge,
		Path:     path,
		Domain:   host,
		SameSite: sameSite,
		Secure:   secure,
		HttpOnly: httpOnly,
	})
}

func (ctx *Context) SetSecureCookie(name, value string, maxAge int, path, domain string) {
	ctx.SetCookie(name, value, maxAge, path, domain, true, true, http.SameSiteNoneMode)
}

func (ctx *Context) RemoveCookie(name string) {
	ctx.SetCookie(name, "", -1, "/", ctx.request.Host, false, false, http.SameSiteDefaultMode)
}

func (ctx *Context) Cookie(name string) (string, error) {
	cookie, err := ctx.request.Cookie(name)
	if err != nil {
		return "", err
	}

	val, err := url.QueryUnescape(cookie.Value)
	if err != nil {
		return "", err
	}

	return val, nil
}

func (ctx *Context) TemporaryRedirect(url string) error {
	http.Redirect(ctx.response, ctx.request, url, http.StatusTemporaryRedirect)
	return nil
}

func (ctx *Context) PermanentRedirect(url string) error {
	http.Redirect(ctx.response, ctx.request, url, http.StatusPermanentRedirect)
	return nil
}
