package goblin

import (
    "regexp"
    "errors"
)

type Route struct {
    regx *regexp.Regexp
    handler func(*Context)
    router *Router
}

type Router struct {
    routes []Route
}

func NewRouter(config map[string]interface{}) (*Router, error) {
    routes := make([]Route, 0, 64)
    for path, handler := range config {
        regx, err := regexp.Compile(path)
        if err != nil {
            return nil, err
        }
        route := Route{
            regx: regx,
            handler: nil,
            router: nil,
        }
        switch v := handler.(type) {
        case func(*Context):
            route.handler = v
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

func (r *Router) Match(path string, params Params) func(*Context) {
    for _, route := range r.routes {
        matches := route.regx.FindStringSubmatch(path)
        if matches == nil {
            continue
        }
        // fill params
        if len(matches) > 1 {
            names := route.regx.SubexpNames()
            for i := 1; i < len(matches); i++ {
                params[names[i]] = matches[i]
            }
        }
        if route.handler != nil {
            return route.handler
        }
        if route.router != nil {
            index := len(matches[0])
            return route.router.Match(path[index: ], params)
        }
    }
    return nil
}
