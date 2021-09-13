package context

import "encoding/json"

func (ctx *Context) BindJson(dst interface{}) error {
	return json.NewDecoder(ctx.request.Body).Decode(dst)
}
