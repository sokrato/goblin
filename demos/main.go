package main

import (
	"fmt"
	"log"
	"time"

	g "github.com/dlutxx/goblin/http"
	"github.com/dlutxx/goblin/utils"
)

var (
	demo = &Demo{
		Name: "demo",
	}

	settings = utils.Dict{
		g.CfgKeyRoutes: map[string]interface{}{
			`^echo/(?P<msg>.+)$`: demo.Echo,
			`^files/.*$`:         g.FileServer("/home/xu/", "/files/"),
			`^book/(?P<bookid>\d+)/`: map[string]interface{}{
				"^read$":               demo.ReadBook,
				`^buy/(?P<price>\d+)$`: demo.BuyBook,
			},
		},
		g.CfgKeyHandler404: demo.Handle404,
		g.CfgKeyRequestMiddleware: []func(*g.Context){
			demo.RequestMDW,
		},
		g.CfgKeyResponseMiddleware: []func(*g.Context){
			demo.ResponseMDW,
		},
	}
)

type Demo struct {
	Name string
}

func (d *Demo) BuyBook(ctx *g.Context) {
	res := ctx.Res
	res.WriteString("buy book-" + ctx.Req.Params.MustString("bookid") + " with " + ctx.Req.Params.MustString("price"))
}

func (d *Demo) ReadBook(ctx *g.Context) {
	ctx.Res.WriteString("reading book-" + ctx.Req.Params.MustString("bookid"))
}

func (d *Demo) ResponseMDW(ctx *g.Context) {
	startTime := ctx.Extra["startTime"].(time.Time)
	duration := time.Now().Sub(startTime)
	log.Println(ctx.Res.StatusCode(), ctx.Req.Method, ctx.Req.URL, duration)
}

func (d *Demo) RequestMDW(ctx *g.Context) {
	ctx.Extra["startTime"] = time.Now()
	if ctx.Req.Params.MustString("msg") == "500" {
		ctx.Res.Error("you asked for error", 500)
	}
}

func (d *Demo) Echo(ctx *g.Context) {
	msg := ctx.Req.Params.MustString("msg")
	ctx.Res.WriteString(msg + "\n")
}

func (d *Demo) RecordRequest(args utils.Dict) {
	req := args["context"].(*g.Context).Req
	log.Println("New Request: ", req.Method, req.URL)
}

func (d *Demo) Handle404(ctx *g.Context) {
	ctx.Res.WriteHeader(404)
	ctx.Res.WriteString(ctx.Req.RequestURI + " dose not exist\n")
}

func main() {
	app := g.NewApp(settings)
	g.RequestStarted.Connect(demo.RecordRequest)
	fmt.Println("try: curl http://localhost:8888/echo/hello%20Goblin!")
	log.Fatalln(app.ListenAndServe(":8888"))
}
