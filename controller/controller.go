package controller

import (
	"net/http"

	"github.com/goinbox/gomisc"
)

type ActionContext interface {
	Request() *http.Request
	ResponseWriter() http.ResponseWriter

	ResponseBody() []byte
	SetResponseBody(body []byte)
	AppendResponseBody(body []byte)

	SetValue(key string, value interface{})
	Value(key string) interface{}

	BeforeAction()
	AfterAction()
	Destruct()
}

type BaseContext struct {
	req        *http.Request
	respWriter http.ResponseWriter

	respBody []byte
	data     map[string]interface{}
}

func NewBaseContext(req *http.Request, respWriter http.ResponseWriter) *BaseContext {
	return &BaseContext{
		req:        req,
		respWriter: respWriter,

		data: make(map[string]interface{}),
	}
}

func (bc *BaseContext) Request() *http.Request {
	return bc.req
}

func (bc *BaseContext) ResponseWriter() http.ResponseWriter {
	return bc.respWriter
}

func (bc *BaseContext) ResponseBody() []byte {
	return bc.respBody
}

func (bc *BaseContext) SetResponseBody(body []byte) {
	bc.respBody = body
}

func (bc *BaseContext) AppendResponseBody(body []byte) {
	bc.respBody = gomisc.AppendBytes(bc.respBody, body)
}

func (bc *BaseContext) SetValue(key string, value interface{}) {
	bc.data[key] = value
}

func (bc *BaseContext) Value(key string) interface{} {
	return bc.data[key]
}

func (bc *BaseContext) BeforeAction() {
}

func (bc *BaseContext) AfterAction() {
}

func (bc *BaseContext) Destruct() {
}

type Controller interface {
	NewActionContext(req *http.Request, respWriter http.ResponseWriter) ActionContext
}
