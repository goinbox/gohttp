package httpserver

import (
	"net/http"

	"github.com/goinbox/gohttp/router"
	"github.com/goinbox/gomisc"
)

type Action interface {
	router.Action

	Request() *http.Request
	ResponseWriter() http.ResponseWriter

	ResponseBody() []byte
	SetResponseBody(body []byte)
	AppendResponseBody(body []byte)

	SetValue(key string, value interface{})
	Value(key string) interface{}
}

type BaseAction struct {
	router.BaseAction

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

type redirectAction struct {
	*BaseAction

	code int
	url  string
}

func (a *redirectAction) Name() string {
	return "redirect"
}

func (a *redirectAction) Run() {
	http.Redirect(a.ResponseWriter(), a.Request(), a.url, a.code)
}

func Redirect(r *http.Request, w http.ResponseWriter, code int, url string) {
	panic(&redirectAction{
		BaseAction: NewBaseAction(r, w, nil),
		code:       code,
		url:        url,
	})
}
