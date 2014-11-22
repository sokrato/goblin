package main

import (
    "log"
    "time"
    "fmt"
    "net/http"

    g "github.com/dlutxx/goblin"
)

var (
    demo = &Demo{
        Name: "demo",
    }

    settings g.Settings = g.Settings{
        g.CfgKeyRoutes: map[string]interface{}{
            `^echo/(?P<msg>.+)$`: g.HF(demo.Echo),
            `^book/(?P<bookid>\d+)/`: map[string]interface{} {
                "^read$": g.HF(demo.ReadBook),
                `^buy/(?P<price>\d+)$`: g.HF(demo.BuyBook),
            },
        },
        g.CfgKeyHandler404: g.HF(demo.Handle404),
        g.CfgKeyRequestMiddlewares: []g.Handler{
            g.HF(demo.RequestMDW),
        },
        g.CfgKeyResponseMiddlewares: []g.Handler{
            g.HF(demo.ResponseMDW),
        },
    }
)

type Demo struct {
    Name string
}

func (d *Demo) BuyBook(ctx *g.Context) {
    res := ctx.Res
    res.WriteString("buy book-" + ctx.Params["bookid"] + " with " + ctx.Params["price"])
}

func (d *Demo) ReadBook(ctx *g.Context) {
    ctx.Res.WriteString("reading book-" + ctx.Params["bookid"])
}

func (d *Demo) ResponseMDW(ctx *g.Context) {
    startTime := ctx.Extra["startTime"].(time.Time)
    duration := time.Now().Sub(startTime)
    log.Println(ctx.Req.URL, " cost ", duration)
}

func (d *Demo) RequestMDW(ctx *g.Context) {
    if ctx.Params["msg"] == "500" {
        ctx.Res.Error("you asked for error", 500)
        return
    }
    ctx.Extra["startTime"] = time.Now()
}

func (d *Demo) Echo(ctx *g.Context) {
    msg := ctx.Params["msg"]
    ctx.Res.WriteString(msg + "\n")
}

func (d *Demo) HandleEvent(evt string, args...interface{}) {
    req := args[0].(*http.Request)
    log.Println(req.Method, req.URL)
}

func (d *Demo) Handle404(ctx *g.Context) {
    ctx.Res.WriteString(ctx.Req.RequestURI + " dose not exist")
}

func main() {
    app := g.NewApp(settings)
    app.On("request.new", demo)
    fmt.Println("try: curl 'http://localhost:8888/echo/hello%20Goblin!'")
    log.Fatalln(app.ListenAndServe(":8888"))
}
