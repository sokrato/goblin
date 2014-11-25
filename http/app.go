package http

import (
	"github.com/dlutxx/goblin/signal"
	"github.com/dlutxx/goblin/utils"
	"net/http"
)

var (
	RequestStarted  = signal.New("context")
	RequestFinished = signal.New("context")
)

type App struct {
	Settings            utils.Dict
	router              *Router
	handler404          func(*Context)
	handler500          func(*Context)
	requestMiddlewares  []func(*Context)
	responseMiddlewares []func(*Context)
}

func (app *App) catchInternalError(ctx *Context) {
	defer func() {
		if err := recover(); err != nil {
			ctx.err = err
			handle500(ctx)
		}
	}()

	if err := recover(); err != nil {
		ctx.err = err
		if app.handler500 != nil {
			app.handler500(ctx)
		} else {
			handle500(ctx)
		}
	}
}

func (app *App) createContext(w http.ResponseWriter, req *http.Request) *Context {
	Req := NewRequest(req, utils.Dict{})
	return &Context{
		Res:   NewResponse(w, Req),
		Req:   Req,
		App:   app,
		Extra: utils.Dict{},
	}
}

func (app *App) Handle404(ctx *Context) {
	if app.handler404 != nil {
		app.handler404(ctx)
	} else {
		handle404(ctx)
	}
	return
}

func (app *App) applyRequestMiddlewares(ctx *Context) {
	for _, mw := range app.requestMiddlewares {
		mw(ctx)
	}
}

func (app *App) applyResponseMiddlewares(ctx *Context) {
	for _, mw := range app.responseMiddlewares {
		mw(ctx)
	}
}

// Beware that response may have been flush in middlewares.
func (app *App) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := app.createContext(w, req)
	defer app.catchInternalError(ctx)
	RequestStarted.Send(utils.Dict{"context": ctx})

	view := app.router.Match(req.URL.Path[1:], ctx.Req.Params)
	if view == nil { // 404
		app.Handle404(ctx)
	} else {
		app.applyRequestMiddlewares(ctx)
		if !ctx.Res.HeaderSent() {
			view(ctx)
		}
		app.applyResponseMiddlewares(ctx)
	}
	RequestFinished.Send(utils.Dict{"context": ctx})
}

// It will panic if settings is invalid.
func NewApp(config utils.Dict) *App {
	return &App{
		router:              parseRouter(config),
		handler404:          parseHandler404(config),
		handler500:          parseHandler500(config),
		Settings:            config,
		requestMiddlewares:  parseRequestMiddlewares(config),
		responseMiddlewares: parseResponseMiddlewares(config),
	}
}

func (app *App) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, app)
}
