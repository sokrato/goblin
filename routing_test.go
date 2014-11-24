package goblin

import (
    "testing"
    "net/http"
)

func TestNewRouterWithNilOrInvalidView(t *testing.T) {
    flag := false
    // route config with immediate nil view
    router, err := NewRouter(map[string]interface{}{
        "^nil$": nil,
    })
    if router != nil || err == nil {
        t.FailNow()
    }
    // route config with nested nil view
    router, err = NewRouter(map[string]interface{}{
        "^non-nil$": HandlerFromFunc(func(res *ResponseWriter, req *http.Request) {
            flag = true
        }),
        "^prefx/": map[string]interface{} {
            "^index$": nil,
        },
    })
    if router != nil || err == nil {
        t.FailNow()
    }
    // route config with invalid view
    router, err = NewRouter(map[string]interface{}{
        "^int-view$": 123,
        "^string-view$": "view",
    })
    if router != nil || err == nil {
        t.FailNow()
    }
}

func TestRouterFind(t *testing.T) {
    viewCalled := false
    handler := HandlerFromFunc(func(res *ResponseWriter, req *http.Request) {
        viewCalled = true
    })
    router, err := NewRouter(map[string]interface{}{
        "^view1$": handler,
        "^api/": map[string]interface{}{
            "^doc$": handler,
            `^user/(\d+)$`: handler,
        },
    })
    if err != nil {
        t.FailNow()
    }
    view := router.Find("view1")
    if view == nil {
        t.FailNow()
    }
    if view.Handle(nil, nil); !viewCalled {
        t.FailNow()
    }
    if view = router.Find("aboutwhat"); view != nil {
        t.FailNow()
    }
    if view = router.Find("api"); view != nil {
        t.FailNow()
    }
    viewCalled = false
    view = router.Find("api/doc")
    if view == nil {
        t.FailNow()
    }
    if view.Handle(nil, nil); !viewCalled {
        t.FailNow()
    }
    viewCalled = false
    view = router.Find("api/user/123")
    if view.Handle(nil, nil); !viewCalled {
        t.Fail()
    }
}
