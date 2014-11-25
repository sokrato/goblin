package http

import (
	"fmt"
	"net/http"
	"os"
)

func handle404(ctx *Context) {
	ctx.Res.NotFound()
}

func handle500(ctx *Context) {
	ctx.Res.Error(http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	fmt.Fprintln(os.Stderr, ctx.Req.URL, ctx.Err())
}

func FileServer(root, prefix string) func(*Context) {
	fs := http.Dir(root)
	httpHandler := http.StripPrefix(prefix, http.FileServer(fs))
	return func(ctx *Context) {
		httpHandler.ServeHTTP(ctx.Res.w, ctx.Req.Request)
	}
}
