package goblin

import "net/http"

type Handler interface {
    Handle(*ResponseWriter, *http.Request)
}

type SimpleHandler struct {
    fn func(*ResponseWriter, *http.Request)
}

func (sv SimpleHandler) Handle(res *ResponseWriter, req *http.Request) {
    sv.fn(res, req)
}

func HandlerFromFunc(fn func(*ResponseWriter, *http.Request)) Handler {
    return SimpleHandler{fn}
}

func handle404(res *ResponseWriter, req *http.Request) {
    res.NotFound()
}

func handle500(res *ResponseWriter, req *http.Request) {
    res.Error(http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
