package goblin

import (
    "net/http"
    "bytes"
    "time"
    "io"
)

type Data map[string]interface{}

func (d Data) MustInt(key string) int {
    val, ok := d[key].(int)
    if !ok {
        panic("Data[" + key + "] is not an int")
    }
    return val
}

func (d Data) MustString(key string) string {
    val, ok := d[key].(string)
    if !ok {
        panic("Data[" + key + "] is not a string")
    }
    return val
}

type ResponseWriter struct {
    Body *bytes.Buffer
    isStreaming bool
    sent bool
    w http.ResponseWriter
    req *http.Request
    Extra map[string]interface{}
}

func NewResponseWriter(w http.ResponseWriter, req *http.Request) *ResponseWriter {
    return &ResponseWriter{
        Body: new(bytes.Buffer),
        isStreaming: false,
        sent: false,
        w: w,
        req: req,
        Extra: make(Data, 0),
    }
}

func (res *ResponseWriter) Sent() bool {
    return res.sent
}

func (res *ResponseWriter) markSent() {
    res.sent = true
}

func (res *ResponseWriter) Header() http.Header {
    return res.w.Header()
}

func (res *ResponseWriter) SetStreaming(isStreaming bool) {
    if isStreaming && res.Body.Len() > 0{
        res.Flush()
    }
    res.isStreaming = isStreaming
}

func (res *ResponseWriter) Streaming() bool {
    return res.isStreaming
}

func (res *ResponseWriter) Write(bs []byte) (int, error) {
    if res.Streaming() {
        res.markSent()
        return res.w.Write(bs)
    }
    return res.Body.Write(bs)
}

func (res *ResponseWriter) Flush() (int, error) {
    res.markSent()
    n, e := res.w.Write(res.Body.Bytes())
    res.Body.Reset()
    return n, e
}

func (res *ResponseWriter) WriteHeader(status int) {
    res.markSent()
    res.w.WriteHeader(status)
}

// convinience func

func (res *ResponseWriter) Error(err string, code int) {
    res.markSent()
    http.Error(res.w, err, code)
}

func (res *ResponseWriter) NotFound() {
    res.markSent()
    http.NotFound(res.w, res.req)
}

func (res *ResponseWriter) Redirect(urlStr string, code int) {
    res.markSent()
    http.Redirect(res.w, res.req, urlStr, code)
}

func (res *ResponseWriter) ServeContent(name string, modtime time.Time, content io.ReadSeeker) {
    res.markSent()
    http.ServeContent(res.w, res.req, name, modtime, content)
}

func (res *ResponseWriter) ServeFile(name string) {
    res.markSent()
    http.ServeFile(res.w, res.req, name)
}

func (res *ResponseWriter) SetCookie(cookie *http.Cookie) {
    http.SetCookie(res.w, cookie)
}
