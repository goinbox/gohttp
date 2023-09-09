package httpserver

import (
	"net/http"

	"github.com/goinbox/gomisc"
	"github.com/goinbox/pcontext"
	"github.com/goinbox/router"
)

type Action[T pcontext.Context] interface {
	router.Action[T]

	Init(r *http.Request, w http.ResponseWriter, args []string) T

	Request() *http.Request
	ResponseWriter() http.ResponseWriter
	Args() []string

	ResponseBody() []byte
	SetResponseBody(body []byte)
	AppendResponseBody(body []byte)
}

type BaseAction[T pcontext.Context] struct {
	router.BaseAction[T]

	req        *http.Request
	respWriter http.ResponseWriter
	args       []string

	respBody []byte
}

func (a *BaseAction[T]) Init(r *http.Request, w http.ResponseWriter, args []string) {
	a.req = r
	a.respWriter = w
	a.args = args
}

func (a *BaseAction[T]) Request() *http.Request {
	return a.req
}

func (a *BaseAction[T]) ResponseWriter() http.ResponseWriter {
	return a.respWriter
}

func (a *BaseAction[T]) Args() []string {
	return a.args
}

func (a *BaseAction[T]) ResponseBody() []byte {
	return a.respBody
}

func (a *BaseAction[T]) SetResponseBody(body []byte) {
	a.respBody = body
}

func (a *BaseAction[T]) AppendResponseBody(body []byte) {
	a.respBody = gomisc.AppendBytes(a.respBody, body)
}
