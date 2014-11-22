package goblin

import (
    "net/http"
    "time"
    "io"
    "encoding/json"
)

type Response struct {
    status int
    headerSent bool
    w http.ResponseWriter
    req *http.Request
}

func NewResponse(w http.ResponseWriter, req *http.Request) *Response {
    return &Response{
        status: http.StatusOK,
        w: w,
        req: req,
    }
}

// get status code of current response
func (res *Response) Status() int {
    return res.status
}

func (res *Response) Header() http.Header {
    return res.w.Header()
}

func (res *Response) HeaderSent() bool {
    return res.headerSent
}

func (res *Response) markHeaderSent() {
    res.headerSent = true
}

func (res *Response) Write(bs []byte) (int, error) {
    res.markHeaderSent()
    return res.w.Write(bs)
}

func (res *Response) WriteString(str string) (int, error) {
    return res.Write([]byte(str))
}

// auto set content-type
func (res *Response) WriteJSON(v interface{}) (int, error) {
    data, err := json.Marshal(v)
    if err != nil {
        return 0, err
    }
    return res.Write(data)
}

// convinience wrappers for http.ResponseWriter

func (res *Response) WriteHeader(status int) {
    res.w.WriteHeader(status)
    res.status = status
    res.markHeaderSent()
}

func (res *Response) Error(err string, code int) {
    http.Error(res.w, err, code)
    res.status = code
    res.markHeaderSent()
}

func (res *Response) NotFound() {
    http.NotFound(res.w, res.req)
    res.status = http.StatusNotFound
    res.markHeaderSent()
}

func (res *Response) Redirect(urlStr string, code int) {
    http.Redirect(res.w, res.req, urlStr, code)
    res.status = code
    res.markHeaderSent()
}

func (res *Response) ServeContent(name string, modtime time.Time, content io.ReadSeeker) {
    http.ServeContent(res.w, res.req, name, modtime, content)
    res.markHeaderSent()
}

func (res *Response) ServeFile(name string) {
    http.ServeFile(res.w, res.req, name)
    res.markHeaderSent()
}

func (res *Response) SetCookie(cookie *http.Cookie) {
    http.SetCookie(res.w, cookie)
}
