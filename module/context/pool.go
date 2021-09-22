package context

import (
	"net/http"
	"sync"
)

var pool = sync.Pool{
	New: func() interface{} {
		return new(Context)
	},
}

func Acquire(rw http.ResponseWriter, r *http.Request) *Context {
	c := pool.Get().(*Context)
	c.request = r
	c.response = rw
	c.user = nil
	return c
}

func Release(c *Context) {
	c.request = nil
	c.response = nil
	c.user = nil
	pool.Put(c)
}
