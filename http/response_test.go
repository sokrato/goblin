package http

import (
	"github.com/dlutxx/goblin/utils"
	"net/http"
	"testing"
)

type MockResponse struct {
	header http.Header
}

func (mr *MockResponse) Write(bs []byte) (int, error) {
	return len(bs), nil
}

func (mr *MockResponse) WriteHeader(h int) {
	//
}

func (mr *MockResponse) Header() http.Header {
	return mr.header
}

func mockNewResponse() *Response {
	mr := &MockResponse{
		header: http.Header{},
	}
	req, err := http.NewRequest("GET", "https://www.google.com/search?key=go", nil)
	if err != nil {
		panic(req)
	}
	q := NewRequest(req, utils.Dict{})
	return NewResponse(mr, q)
}

func TestStatus(t *testing.T) {
	resp := mockNewResponse()
	resp.WriteHeader(200)
	if resp.StatusCode() != 200 || !resp.HeaderSent() {
		t.FailNow()
	}

	resp = mockNewResponse()
	resp.Error("teapot", 418)
	if resp.StatusCode() != 418 || !resp.HeaderSent() {
		t.FailNow()
	}

	resp = mockNewResponse()
	resp.NotFound()
	if resp.StatusCode() != http.StatusNotFound || !resp.HeaderSent() {
		t.FailNow()
	}

	resp = mockNewResponse()
	resp.RedirectTemp("/redirect")
	if resp.StatusCode() != http.StatusFound || !resp.HeaderSent() {
		t.FailNow()
	}

	resp = mockNewResponse()
	resp.RedirectPerm("/redirect")
	if resp.StatusCode() != http.StatusMovedPermanently || !resp.HeaderSent() {
		t.FailNow()
	}
}
