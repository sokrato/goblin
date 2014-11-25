package http

import (
	"github.com/dlutxx/goblin/utils"
	"log"
	"testing"
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
		"^non-nil$": func(ctx *Context) {
			flag = true
		},
		"^prefx/": map[string]interface{}{
			"^index$": nil,
		},
	})
	if router != nil || err == nil {
		t.FailNow()
	}
	// route config with invalid Handler
	router, err = NewRouter(map[string]interface{}{
		"^int-view$":    123,
		"^string-view$": "view",
	})
	if router != nil || err == nil {
		t.FailNow()
	}
}

func TestRouterMatch(t *testing.T) {
	handlerCalled := false
	handler := func(ctx *Context) {
		handlerCalled = true
	}
	router, err := NewRouter(map[string]interface{}{
		"^handler1$": handler,
		"^api/": map[string]interface{}{
			"^doc$":               handler,
			`^user/(?P<uid>\d+)$`: handler,
		},
	})
	if err != nil {
		log.Println("NewRouter error", err)
		t.FailNow()
	}
	params := utils.Dict{}
	handler = router.Match("handler1", params)
	if handler == nil {
		log.Println("cannot find handler1")
		t.FailNow()
	}
	if handler(nil); !handlerCalled {
		log.Println("handler1 not called")
		t.FailNow()
	}
	if handler = router.Match("aboutwhat", params); handler != nil {
		log.Println("aboutwhat not found")
		t.FailNow()
	}
	if handler = router.Match("api/", params); handler != nil {
		log.Println("api/ should match any handler")
		t.FailNow()
	}
	handlerCalled = false
	handler = router.Match("api/doc", params)
	if handler == nil {
		log.Println("cannot find api/doc")
		t.FailNow()
	}
	if handler(nil); !handlerCalled {
		log.Println("doc handler not called")
		t.FailNow()
	}
	handlerCalled = false
	handler = router.Match("api/user/123", params)
	if handler(nil); !handlerCalled {
		log.Println("user handler not called")
		t.FailNow()
	} else if uid, err := params.String("uid"); err != nil || uid != "123" {
		log.Println("param uid not found:", err)
		t.FailNow()
	}
}
