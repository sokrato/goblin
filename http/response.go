package http

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type Response struct {
	statusCode int
	headerSent bool
	w          http.ResponseWriter
	req        *Request
}

func NewResponse(w http.ResponseWriter, req *Request) *Response {
	return &Response{
		statusCode: http.StatusOK,
		w:          w,
		req:        req,
	}
}

func (res *Response) ResponseWriter() http.ResponseWriter {
	return res.w
}

// get statusCode code of current response
func (res *Response) StatusCode() int {
	return res.statusCode
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

func (res *Response) WriteHeader(statusCode int) {
	res.w.WriteHeader(statusCode)
	res.statusCode = statusCode
	res.markHeaderSent()
}

func (res *Response) Error(err string, code int) {
	http.Error(res.w, err, code)
	res.statusCode = code
	res.markHeaderSent()
}

func (res *Response) NotFound() {
	http.NotFound(res.w, res.req.Request)
	res.statusCode = http.StatusNotFound
	res.markHeaderSent()
}

func (res *Response) Redirect(urlStr string, code int) {
	http.Redirect(res.w, res.req.Request, urlStr, code)
	res.statusCode = code
	res.markHeaderSent()
}

func (res *Response) RedirectPerm(urlStr string) {
	res.Redirect(urlStr, http.StatusMovedPermanently)
}

func (res *Response) RedirectTemp(urlStr string) {
	res.Redirect(urlStr, http.StatusFound)
}

func (res *Response) ServeContent(name string, modtime time.Time, content io.ReadSeeker) {
	http.ServeContent(res.w, res.req.Request, name, modtime, content)
	res.markHeaderSent()
}

func (res *Response) ServeFile(name string) {
	http.ServeFile(res.w, res.req.Request, name)
	res.markHeaderSent()
}

func (res *Response) SetCookie(cookie *http.Cookie) {
	http.SetCookie(res.w, cookie)
}

func (res *Response) SetHeader(key, val string) {
	res.w.Header().Set(key, val)
}

func (res *Response) AddHeader(key, val string) {
	res.w.Header().Add(key, val)
}
