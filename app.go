package goblin

import (
    "net/http"
    "errors"
)

type App struct {
    EventEmitter
    Router *Router
    Handler404 View
    Handler500 View
    requestMiddlewares []View
    responseMiddlewares []View
}

func (app *App) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    view := app.Router.Find(req.URL.Path[1: ])
    res := NewResponseWriter(w, req)
    defer func() {
        defer func() {
            if err := recover(); err != nil {
                handle500(res, req)
            }
        }()

        if err := recover(); err != nil {
            app.Emit("500", res, req)
            app.Handler500(res, req)
        }
    }()

    if view == nil {
        // 404
        app.Emit("404", res, req)
        app.Handler404(res, req)
        return
    }
    // apply request middlewares
    for _, mw := range app.requestMiddlewares {
        if mw(res, req); res.Sent() {
            break
        }
    }
    // response middlewares
    for _, mv := range app.responseMiddlewares {
        if mv(res, req); res.Sent() {
            break
        }
    }
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

    var hdl404, hdl500 View
    hdl404Cfg, ok := settings["handle404"]
    if ok {
        hdl404, ok = hdl404Cfg.(View)
        if !ok {
            return nil, errors.New("invalid 404 handler")
        }
    } else {
        hdl404 = handle404
    }
    hdl500Cfg, ok := settings["handle500"]
    if ok {
        hdl500, ok = hdl500Cfg.(View)
        if !ok {
            return nil, errors.New("invalid 500 handler")
        }
    } else {
        hdl500 = handle500
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
