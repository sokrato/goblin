package goblin

import (
    "net/http"
    "errors"
    "strconv"
    "fmt"
    "os"
)

var (
    ErrParamsNotSet = errors.New("params not set")
)

type Context struct {
    Res *ResponseWriter
    Req *Request
    App *App // the main app
    Params Params // request params
    Err interface{} // internal error, or nil
    Extra Extra
}

type Extra map[string]interface{}

type Params map[string]string

func (p Params) Int(key string) (int, error) {
    val, ok := p[key]
    if !ok {
        return 0, ErrParamsNotSet
    }
    return strconv.Atoi(val)
}

type Handler interface {
    Handle(*Context)
}

type SimpleHandler struct {
    fn func(*Context)
}

func (sv SimpleHandler) Handle(ctx *Context) {
    sv.fn(ctx)
}

func HF(fn func(*Context)) Handler {
    return SimpleHandler{fn}
}

func handle404(ctx *Context) {
    ctx.Res.NotFound()
}

func handle500(ctx *Context) {
    ctx.Res.Error(http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    fmt.Fprintln(os.Stderr, ctx.Err)
}
