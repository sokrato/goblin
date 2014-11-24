package goblin

import (
    "regexp"
    "errors"
)

type Route struct {
    regx *regexp.Regexp
    view Handler
    router *Router
}

type Router struct {
    routes []Route
}

func NewRouter(config map[string]interface{}) (*Router, error) {
    routes := make([]Route, 0, 64)
    for path, view := range config {
        regx, err := regexp.Compile(path)
        if err != nil {
            return nil, err
        }
        route := Route{
            regx: regx,
            view: nil,
            router: nil,
        }
        switch v := view.(type) {
        case Handler:
            route.view = v
        case map[string]interface{}:
            route.router, err = NewRouter(v)
            if err != nil {
                return nil, err
            }
        default:
            return nil, errors.New("invalid route at " + path)
        }
        routes = append(routes, route)
    }
    return &Router{routes}, nil
}

func (r *Router)Find(path string) Handler {
    for _, route := range r.routes {
        locs := route.regx.FindStringIndex(path)
        if locs == nil {
            continue
        }
        if route.view != nil {
            return route.view
        }
        if route.router != nil {
            return route.router.Find(path[locs[1]:])
        }
    }
    return nil
}
