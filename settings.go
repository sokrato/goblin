package goblin

import (
    "errors"
)

var (
    ErrInvalidSettings = errors.New("goblin: Invalid Settings")

    /* defaultSettings = Settings{
        "debug": false,
        "env": "dev",
    } */
)

const (
    CfgKeyRoutes = "routes"
    CfgKeyHandler404 = "handler404"
    CfgKeyHandler500 = "handler500"
    CfgKeyRequestMiddleware = "requestMiddlewares"
    CfgKeyResponseMiddleware = "responseMiddlewares"
)


type Settings map[string]interface{}


func (s Settings) Router() *Router {
    routeCfg, ok := s[CfgKeyRoutes]
    if !ok {
        panic("goblin: routes not found in settings")
    }
    routes, ok := routeCfg.(map[string]interface{})
    if !ok {
        panic("goblin: invalid routes settings")
    }
    router, err := NewRouter(routes)
    if err != nil {
        panic(err)
    }
    return router
}

func (s Settings) getHandler(key string) func(*Context) {
    val, ok := s[key]
    if ok {
        handler, ok := val.(func(*Context))
        if !ok {
            panic("goblin: invalid settings for " + key)
        }
        return handler
    }
    return nil
}

func (s Settings) Handler404() func(*Context) {
    return s.getHandler(CfgKeyHandler404)
}

func (s Settings) Handler500() func(*Context) {
    return s.getHandler(CfgKeyHandler500)
}

func (s Settings) getHandlerSlice(key string) []func(*Context) {
    val, ok := s[key]
    if !ok {
        return nil
    }
    switch v := val.(type) {
    case []func(*Context):
        return v
    case func(*Context):
        return []func(*Context){v}
    default:
        panic("goblin: invalid settings for " + key)
    }
}

func (s Settings) RequestMiddlewares() []func(*Context) {
    return s.getHandlerSlice(CfgKeyRequestMiddleware)
}

func (s Settings) ResponseMiddlewares() []func(*Context) {
    return s.getHandlerSlice(CfgKeyResponseMiddleware)
}

func (s Settings) String(key string) (string, error) {
    val, ok := s[key]
    if !ok {
        return "", ErrInvalidSettings
    }
    valStr, ok := val.(string)
    if !ok {
        return "", ErrInvalidSettings
    }
    return valStr, nil
}

func (s Settings) Int(key string) (int, error) {
    val, ok := s[key]
    if !ok {
        return 0, ErrInvalidSettings
    }
    valInt, ok := val.(int)
    if !ok {
        return 0, ErrInvalidSettings
    }
    return valInt, nil
}
