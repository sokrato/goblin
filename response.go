package goblin

import (
    "net/http"
    "bytes"
    "time"
    "io"
    "errors"
)


var (
    ErrResetAfterFlush = errors.New("cannot reset response after content flushed")
)

type ResponseWriter struct {
    Buffer *bytes.Buffer
    isStreaming bool
    flushed bool
    written bool
    w http.ResponseWriter
    req *http.Request
}

func NewResponseWriter(w http.ResponseWriter, req *http.Request) *ResponseWriter {
    return &ResponseWriter{
        Buffer: new(bytes.Buffer),
        isStreaming: false,
        flushed: false,
        written: false,
        w: w,
        req: req,
    }
}

func (res *ResponseWriter) Written() bool {
    return res.written || res.flushed
}

func (res *ResponseWriter) Flushed() bool {
    return res.flushed
}

func (res *ResponseWriter) Flush() (int, error) {
    n, e := res.w.Write(res.Buffer.Bytes())
    res.markFlushed()
    res.Buffer.Reset()
    return n, e
}

func (res *ResponseWriter) markFlushed() {
    res.flushed = true
}

func (res *ResponseWriter) Header() http.Header {
    return res.w.Header()
}

func (res *ResponseWriter) SetStreaming() {
    res.isStreaming = true
    if res.Buffer.Len() > 0 {
        res.Flush()
    }
}

func (res *ResponseWriter) IsStreaming() bool {
    return res.isStreaming
}

func (res *ResponseWriter) Write(bs []byte) (int, error) {
    if res.IsStreaming() {
        res.markFlushed()
        return res.w.Write(bs)
    }
    res.written = true
    return res.Buffer.Write(bs)
}

func (res *ResponseWriter) WriteString(str string) (int, error) {
    return res.Write([]byte(str))
}

func (res *ResponseWriter) Reset() error {
    if res.Flushed() {
        return ErrResetAfterFlush
    }
    res.Buffer.Reset()
    res.written = false
    return nil
}

func (res *ResponseWriter) WriteHeader(status int) {
    res.markFlushed()
    res.w.WriteHeader(status)
}

// convinience func

func (res *ResponseWriter) Error(err string, code int) {
    res.markFlushed()
    http.Error(res.w, err, code)
}

func (res *ResponseWriter) NotFound() {
    res.markFlushed()
    http.NotFound(res.w, res.req)
}

func (res *ResponseWriter) Redirect(urlStr string, code int) {
    res.markFlushed()
    http.Redirect(res.w, res.req, urlStr, code)
}

func (res *ResponseWriter) ServeContent(name string, modtime time.Time, content io.ReadSeeker) {
    res.markFlushed()
    http.ServeContent(res.w, res.req, name, modtime, content)
}

func (res *ResponseWriter) ServeFile(name string) {
    res.markFlushed()
    http.ServeFile(res.w, res.req, name)
}

func (res *ResponseWriter) SetCookie(cookie *http.Cookie) {
    http.SetCookie(res.w, cookie)
}
