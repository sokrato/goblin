package main

import (
    "log"
    "flag"
    // "fmt"
    "strings"
    g "github.com/dlutxx/goblin"
)

var addr = flag.String("addr", ":8888", "server address")

type HttpBin struct {
}

func (hb *HttpBin) Home(ctx *g.Context) {
    ctx.Res.SetHeader("Content-Type", "text/html; charset=utf-8")
    ctx.Res.WriteString(`This is a demo project to mimic the famous <a href="http://httpbin.org">httpbin.org</a>`)
}

func (hb *HttpBin) returnJSON(res *g.Response, v interface{}) {
    res.SetHeader("Content-Type", "application/json; charset=utf-8")
    res.WriteJSON(v)
}

func (hb *HttpBin) getIP(req *g.Request) string {
    parts := strings.Split(req.RemoteAddr, ":")
    return parts[0]
}

func (hb *HttpBin) IP(ctx *g.Context) {
    data := map[string]string {
        "origin": hb.getIP(ctx.Req),
    }
    hb.returnJSON(ctx.Res, data)
}

func (hb *HttpBin) getHeaders(req *g.Request) map[string]string {
    headers := map[string]string{}
    for k, vals := range req.Header {
        headers[k] = strings.Join(vals, ", ")
    }
    return headers
}

func (hb *HttpBin) Headers(ctx *g.Context) {
    hb.returnJSON(ctx.Res, hb.getHeaders(ctx.Req))
}

func (hb *HttpBin) UserAgent(ctx *g.Context) {
    data := map[string]string {
        "user-agent": ctx.Req.Header["User-Agent"][0],
    }
    hb.returnJSON(ctx.Res, data)
}

func (hb *HttpBin) getArgs(req *g.Request) map[string]interface{} {
    query := req.URL.Query()
    args := map[string]interface{}{}
    for k, vals := range query {
        if len(vals) > 1 {
            args[k] = vals
        } else {
            args[k] = vals[0]
        }
    }
    return args
}

func (hb *HttpBin) Get(ctx *g.Context) {
    data := map[string]interface{} {
        "url": ctx.Req.URL.String(),
        "args": hb.getArgs(ctx.Req),
        "origin": hb.getIP(ctx.Req),
        "headers": hb.getHeaders(ctx.Req),
    }
    hb.returnJSON(ctx.Res, data)
}

func (hb *HttpBin) Status(ctx *g.Context) {
    code, _ := ctx.Params.Int("code")
    ctx.Res.WriteHeader(code)
}

func main() {
    flag.Parse()

    hb := &HttpBin{}
    cfg := g.Settings{
        g.CfgKeyRoutes: map[string]interface{}{
            `^$`: g.HF(hb.Home),
            `^ip$`: g.HF(hb.IP),
            `^headers$`: g.HF(hb.Headers),
            `^user-agent$`: g.HF(hb.UserAgent),
            `^get$`: g.HF(hb.Get),
            `^status/(?P<code>\d{3})$`: g.HF(hb.Status),
        },
    }
    app := g.NewApp(cfg)
    log.Fatalln(app.ListenAndServe(*addr))
}
