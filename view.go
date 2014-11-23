package goblin

import (
    "net/http"
    "errors"
    "strconv"
    "fmt"
    "os"
)

var (
    ErrParamNotSet = errors.New("param not set")
)

type Context struct {
    Res *Response
    Req *Request
    App *App // the main app
    Params Params // request params
    err interface{} // internal error, or nil
    Extra Extra
}

func (ctx *Context) Err() interface{} {
    return ctx.err
}

// extra data bound to a Context instance
type Extra map[string]interface{}

// named groups in url pattern
type Params map[string]string

func (p Params) Int(key string) (int, error) {
    val, ok := p[key]
    if !ok {
        return 0, ErrParamNotSet
    }
    return strconv.Atoi(val)
}

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
