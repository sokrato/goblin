package http

import (
	"github.com/dlutxx/goblin/utils"
	"net/http"
)

type Request struct {
	*http.Request
	Params utils.Dict
}

func NewRequest(req *http.Request, params utils.Dict) *Request {
	return &Request{
		Request: req,
		Params:  params,
	}
}
