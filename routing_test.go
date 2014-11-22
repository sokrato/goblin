package goblin

import (
    "testing"
    // "net/http"
)

func TestNewRouterWithNilOrInvalidHandler(t *testing.T) {
    flag := false
    // route config with immediate nil Handler
    router, err := NewRouter(map[string]interface{}{
        "^nil$": nil,
    })
    if router != nil || err == nil {
        t.FailNow()
    }
    // route config with nested nil Handler
    router, err = NewRouter(map[string]interface{}{
        "^non-nil$": HandlerFromFunc(func(ctx *Context) {
            flag = true
        }),
        "^prefx/": map[string]interface{} {
            "^index$": nil,
        },
    })
    if router != nil || err == nil {
        t.FailNow()
    }
    // route config with invalid Handler
    router, err = NewRouter(map[string]interface{}{
        "^int-view$": 123,
        "^string-view$": "view",
    })
    if router != nil || err == nil {
        t.FailNow()
    }
}

func TestRouterMatch(t *testing.T) {
    handlerCalled := false
    handler := HandlerFromFunc(func(ctx *Context) {
        handlerCalled = true
    })
    router, err := NewRouter(map[string]interface{}{
        "^handler1$": handler,
        "^api/": map[string]interface{}{
            "^doc$": handler,
            `^user/(?P<uid>\d+)$`: handler,
        },
    })
    if err != nil {
        t.FailNow()
    }
    params := Params{}
    handler = router.Match("handler1", params)
    if handler == nil {
        t.FailNow()
    }
    if handler.Handle(nil); !handlerCalled {
        t.FailNow()
    }
    if handler = router.Match("aboutwhat", params); handler != nil {
        t.FailNow()
    }
    if handler = router.Match("api/", params); handler != nil {
        t.FailNow()
    }
    handlerCalled = false
    handler = router.Match("api/doc", params)
    if handler == nil {
        t.FailNow()
    }
    if handler.Handle(nil); !handlerCalled {
        t.FailNow()
    }
    handlerCalled = false
    handler = router.Match("api/user/123", params)
    if handler.Handle(nil); !handlerCalled{
        t.FailNow()
    } else if uid, err := params.Int("uid"); err!=nil || uid != 123 {
        t.FailNow()
    }
}
