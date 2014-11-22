package goblin

import (
    "testing"
    "net/http"
)

func makeTestView(flag *bool) View{
    return func (res *ResponseWriter, req *http.Request) {
        *flag = true
    }
}

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
        "^non-nil$": makeTestView(&flag),
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
    router, err := NewRouter(map[string]interface{}{
        "^view1$": makeTestView(&viewCalled),
        "^api/": map[string]interface{}{
            "^doc$": makeTestView(&viewCalled),
            `^user/(\d+)$`: makeTestView(&viewCalled),
        },
    })
    if err != nil {
        t.FailNow()
    }
    view := router.Find("view1")
    if view == nil {
        t.FailNow()
    }
    if view(nil, nil); !viewCalled {
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
    if view(nil, nil); !viewCalled {
        t.FailNow()
    }
    viewCalled = false
    view = router.Find("api/user/123")
    if view(nil, nil); !viewCalled {
        t.Fail()
    }
}
