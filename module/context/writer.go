package context

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
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
