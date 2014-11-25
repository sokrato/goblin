package http

import (
	"github.com/dlutxx/goblin/utils"
)

type Context struct {
	Res   *Response
	Req   *Request
	App   *App // the main app
	Extra utils.Dict
	err   interface{} // internal error, or nil
}

func (ctx *Context) Err() interface{} {
	return ctx.err
}
