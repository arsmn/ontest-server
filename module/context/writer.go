package context

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func (ctx *Context) JSON(data interface{}) error {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(data); err != nil {
		return err
	}

	ctx.response.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, err := ctx.response.Write(buf.Bytes())
	return err
}

func (ctx *Context) Status(status int) *Context {
	ctx.response.WriteHeader(status)
	return ctx
}

func (ctx *Context) SendStatus(status int) error {
	ctx.Status(status)
	http.Error(ctx.response, http.StatusText(status), status)
	return nil
}

func (ctx *Context) Created(location string, data interface{}) error {
	ctx.response.Header().Set("Location", location)
	ctx.Status(http.StatusCreated)
	return ctx.JSON(data)
}
