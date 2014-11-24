package goblin

import (
    "net/http"
    "errors"
    "log"
)

type App struct {
    EventEmitter
    Router *Router
    Handler404 Handler
    Handler500 Handler
    requestMiddlewares []Handler
    responseMiddlewares []Handler
}

func (app *App) catchInternalError(res *ResponseWriter, req *http.Request) {
    defer func() {
        if err := recover(); err != nil {
            res.Body.Reset()
            handle500(res, req)
        }
        res.Flush()
    }()

    if err := recover(); err != nil {
        app.Emit("500", res, req)
        res.Body.Reset()
        if app.Handler500 != nil {
            app.Handler500.Handle(res, req)
        } else {
            handle500(res, req)
        }
    }
}

func (app *App) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    view := app.Router.Find(req.URL.Path[1: ])
    res := NewResponseWriter(w, req)
    defer app.catchInternalError(res, req)

    if view == nil {
        // 404
        app.Emit("404", res, req)
        if app.Handler404 != nil {
            app.Handler404.Handle(res, req)
        } else {
            handle404(res, req)
        }
        res.Flush()
        return
    }
    // apply request middlewares
    for _, mw := range app.requestMiddlewares {
        if mw.Handle(res, req); res.Sent() {
            break
        }
    }
    view.Handle(res, req)
    // response middlewares
    for _, mv := range app.responseMiddlewares {
        if mv.Handle(res, req); res.Sent() {
            break
        }
    }
    log.Println(res.Sent())
    if !res.Sent() {
        res.Flush()
    }
}

func NewApp(settings Settings) (*App, error) {
    routeCfg, ok := settings["routes"]
    if !ok {
        return nil, errors.New("routes not found")
    }
    routes, ok := routeCfg.(map[string]interface{})
    if !ok {
        return nil, errors.New("invalid routes")
    }
    router, err := NewRouter(routes)
    if err != nil {
        return nil, err
    }

    var hdl404, hdl500 Handler
    hdl404Cfg, ok := settings["handle404"]
    if ok {
        hdl404, ok = hdl404Cfg.(Handler)
        if !ok {
            return nil, errors.New("invalid 404 handler")
        }
    }
    hdl500Cfg, ok := settings["handle500"]
    if ok {
        hdl500, ok = hdl500Cfg.(Handler)
        if !ok {
            return nil, errors.New("invalid 500 handler")
        }
    }
    return &App{
        EventEmitter: make(EventEmitter, 0),
        Router: router,
        Handler404: hdl404,
        Handler500: hdl500,
        requestMiddlewares: nil, //reqMiddlewares,
        responseMiddlewares: nil, //resMiddlewares,
    }, nil
}

func (app *App) ListenAndServe(addr string) error {
    return http.ListenAndServe(addr, app)
}
