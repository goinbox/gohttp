package httpserver

import (
	"net/http"

	"github.com/goinbox/gomisc"
)

type Action interface {
	Name() string

	Request() *http.Request
	ResponseWriter() http.ResponseWriter

	ResponseBody() []byte
	SetResponseBody(body []byte)
	AppendResponseBody(body []byte)

	SetValue(key string, value interface{})
	Value(key string) interface{}

	Before()
	Run()
	After()
	Destruct()
}

type BaseAction struct {
	req        *http.Request
	respWriter http.ResponseWriter

	respBody []byte
	data     map[string]interface{}

	Args []string
}

func NewBaseAction(r *http.Request, w http.ResponseWriter, args []string) *BaseAction {
	return &BaseAction{
		req:        r,
		respWriter: w,
		data:       make(map[string]interface{}),
		Args:       args,
	}
}

func (a *BaseAction) Request() *http.Request {
	return a.req
}

func (a *BaseAction) ResponseWriter() http.ResponseWriter {
	return a.respWriter
}

func (a *BaseAction) ResponseBody() []byte {
	return a.respBody
}

func (a *BaseAction) SetResponseBody(body []byte) {
	a.respBody = body
}

func (a *BaseAction) AppendResponseBody(body []byte) {
	a.respBody = gomisc.AppendBytes(a.respBody, body)
}

func (a *BaseAction) SetValue(key string, value interface{}) {
	a.data[key] = value
}

func (a *BaseAction) Value(key string) interface{} {
	return a.data[key]
}

func (a *BaseAction) Before() {
}

func (a *BaseAction) After() {
}

func (a *BaseAction) Destruct() {
}
