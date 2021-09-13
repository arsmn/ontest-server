package context

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

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

func (ctx *Context) SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool) {
	if path == "" {
		path = "/"
	}

	http.SetCookie(ctx.response, &http.Cookie{
		Name:     name,
		Value:    url.QueryEscape(value),
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		SameSite: http.SameSiteDefaultMode,
		Secure:   secure,
		HttpOnly: httpOnly,
	})
}

func (ctx *Context) SetSecureCookie(name, value string, maxAge int) {
	http.SetCookie(ctx.response, &http.Cookie{
		Name:     name,
		Value:    url.QueryEscape(value),
		MaxAge:   maxAge,
		Path:     "/",
		Domain:   ctx.request.Host,
		SameSite: http.SameSiteDefaultMode,
		Secure:   true,
		HttpOnly: true,
	})
}

func (ctx *Context) RemoveCookie(w http.ResponseWriter, r *http.Request, name string) {
	ctx.SetCookie(name, "", -1, "/", ctx.request.Host, false, false)
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
