package http

import (
	"github.com/dlutxx/goblin/utils"
)

const (
	CfgKeyRoutes             = "routes"
	CfgKeyHandler404         = "handler404"
	CfgKeyHandler500         = "handler500"
	CfgKeyRequestMiddleware  = "requestMiddlewares"
	CfgKeyResponseMiddleware = "responseMiddlewares"
)

func parseRouter(s utils.Dict) *Router {
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

func getHandler(s utils.Dict, key string) func(*Context) {
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

func parseHandler404(s utils.Dict) func(*Context) {
	return getHandler(s, CfgKeyHandler404)
}

func parseHandler500(s utils.Dict) func(*Context) {
	return getHandler(s, CfgKeyHandler500)
}

func getHandlerSlice(s utils.Dict, key string) []func(*Context) {
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

func parseRequestMiddlewares(s utils.Dict) []func(*Context) {
	return getHandlerSlice(s, CfgKeyRequestMiddleware)
}

func parseResponseMiddlewares(s utils.Dict) []func(*Context) {
	return getHandlerSlice(s, CfgKeyResponseMiddleware)
}
