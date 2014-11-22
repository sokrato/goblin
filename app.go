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
    Handler404 Handler
    Handler500 Handler
    Settings Settings
    requestMiddlewares []Handler
    responseMiddlewares []Handler
}

func (app *App) catchInternalError(ctx *Context) {
    defer func() {
        if err := recover(); err != nil {
            ctx.Res.Reset()
            ctx.Err = err
            handle500(ctx)
        }
        ctx.Res.Flush()
    }()

    if err := recover(); err != nil {
        ctx.Err = err
        app.Emit(Evt500, ctx)
        ctx.Res.Reset()
        if app.Handler500 != nil {
            app.Handler500.Handle(ctx)
        } else {
            handle500(ctx)
        }
    }
}

func (app *App) createContext(w http.ResponseWriter, req *http.Request) *Context {
    Req := &Request{req}
    return &Context{
        Res: NewResponseWriter(w, req),
        Req: Req,
        App: app,
        Params: Params{},
        Extra: Extra{},
    }
}

func (app *App) handle404(ctx *Context) {
    app.Emit(Evt404, ctx)
    if app.Handler404 != nil {
        app.Handler404.Handle(ctx)
    } else {
        handle404(ctx)
    }
    ctx.Res.Flush()
    return
}

func (app *App) applyRequestMiddlewares(ctx *Context) {
    for _, mw := range app.requestMiddlewares {
        if mw.Handle(ctx); ctx.Res.Written() {
            break
        }
    }
}

func (app *App) applyResponseMiddlewares(ctx *Context) {
    if ctx.Res.Flushed() {
        return
    }
    for _, mv := range app.responseMiddlewares {
        if mv.Handle(ctx); ctx.Res.Flushed() {
            break
        }
    }
    ctx.Res.Flush()
}

func (app *App) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    ctx := app.createContext(w, req)
    defer app.catchInternalError(ctx)
    app.Emit(EvtRequestNew, req)

    view := app.Router.Match(req.URL.Path[1: ], ctx.Params)
    if view == nil { // 404
        app.handle404(ctx)
        return
    }

    app.applyRequestMiddlewares(ctx)
    if !ctx.Res.Written() {
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
        Handler404: settings.Handler404(),
        Handler500: settings.Handler500(),
        Settings: settings,
        requestMiddlewares: settings.RequestMiddlewares(),
        responseMiddlewares: settings.ResponseMiddlewares(),
    }
}

func (app *App) ListenAndServe(addr string) error {
    return http.ListenAndServe(addr, app)
}
