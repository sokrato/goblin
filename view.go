package goblin

import "net/http"

type View func (*ResponseWriter, *http.Request)

func handle404(res *ResponseWriter, req *http.Request) {
    res.NotFound()
}

func handle500(res *ResponseWriter, req *http.Request) {
    res.Error(http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
