package goblin

import (
    "net/http"
)

const (
    Evt404 = "404"
    Evt500 = "500"
    EvtRequestNew = "request.new"
    EvtRequestFinished = "request.finished"
)

type App struct {
    EventEmitter
    Router *Router
    handler404 Handler
    handler500 Handler
    Settings Settings
    requestMiddlewares []Handler
    responseMiddlewares []Handler
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
        app.Emit(Evt500, ctx)
        if app.handler500 != nil {
            app.handler500.Handle(ctx)
        } else {
            handle500(ctx)
        }
    }
}

func (app *App) createContext(w http.ResponseWriter, req *http.Request) *Context {
    Req := &Request{req}
    return &Context{
        Res: NewResponse(w, req),
        Req: Req,
        App: app,
        Params: Params{},
        Extra: Extra{},
    }
}

func (app *App) Handle404(ctx *Context) {
    app.Emit(Evt404, ctx)
    if app.handler404 != nil {
        app.handler404.Handle(ctx)
    } else {
        handle404(ctx)
    }
    return
}

func (app *App) applyRequestMiddlewares(ctx *Context) {
    for _, mw := range app.requestMiddlewares {
        mw.Handle(ctx)
    }
}

func (app *App) applyResponseMiddlewares(ctx *Context) {
    for _, mv := range app.responseMiddlewares {
        mv.Handle(ctx)
    }
}

// Beware that response may have been flush in middlewares.
func (app *App) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    ctx := app.createContext(w, req)
    defer app.catchInternalError(ctx)
    app.Emit(EvtRequestNew, req)

    view := app.Router.Match(req.URL.Path[1: ], ctx.Params)
    if view == nil { // 404
        app.Handle404(ctx)
    }

    app.applyRequestMiddlewares(ctx)

    if view != nil && !ctx.Res.HeaderSent() {
        view.Handle(ctx)
    }

    app.applyResponseMiddlewares(ctx)
    app.Emit(EvtRequestFinished, ctx)
}

// It will panic if settings is invalid.
func NewApp(settings Settings) *App {
    return &App{
        EventEmitter: make(EventEmitter, 0),
        Router: settings.Router(),
        handler404: settings.Handler404(),
        handler500: settings.Handler500(),
        Settings: settings,
        requestMiddlewares: settings.RequestMiddlewares(),
        responseMiddlewares: settings.ResponseMiddlewares(),
    }
}

func (app *App) ListenAndServe(addr string) error {
    return http.ListenAndServe(addr, app)
}
