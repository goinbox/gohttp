package httpserver

import (
	"net/http"

	"github.com/goinbox/gomisc"
)

type Action[T Context] interface {
	Name() string

	Init(r *http.Request, w http.ResponseWriter, args []string) T

	Before(ctx T) error
	Run(ctx T) error
	After(ctx T, err error)

	Request() *http.Request
	ResponseWriter() http.ResponseWriter
	Args() []string

	ResponseBody() []byte
	SetResponseBody(body []byte)
	AppendResponseBody(body []byte)
}

type BaseAction[T Context] struct {
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

func (a *BaseAction[T]) Before(ctx T) error {
	return nil
}

func (a *BaseAction[T]) After(ctx T, err error) {
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
