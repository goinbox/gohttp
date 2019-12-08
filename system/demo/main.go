package main

import (
	"github.com/goinbox/gohttp/controller"
	"github.com/goinbox/gohttp/gracehttp"
	"github.com/goinbox/gohttp/router"
	"github.com/goinbox/gohttp/system"

	"net/http"
)

func main() {
	dcl := new(DemoController)
	r := router.NewSimpleRouter()

	r.DefineRouteItem("^/g/([0-9]+)$", dcl, "get")
	r.MapRouteItems(new(IndexController), dcl)

	sys := system.NewSystem(r)

	gracehttp.ListenAndServe(":8010", sys)
}

type BaseActionContext struct {
	Req        *http.Request
	RespWriter http.ResponseWriter
	RespBody   []byte
}

func (bac *BaseActionContext) Request() *http.Request {
	return bac.Req
}

func (bac *BaseActionContext) ResponseWriter() http.ResponseWriter {
	return bac.RespWriter
}

func (bac *BaseActionContext) ResponseBody() []byte {
	return bac.RespBody
}

func (bac *BaseActionContext) SetResponseBody(body []byte) {
	bac.RespBody = body
}

func (bac *BaseActionContext) BeforeAction() {
	bac.RespBody = append(bac.RespBody, []byte(" index before ")...)
}

func (bac *BaseActionContext) AfterAction() {
	bac.RespBody = append(bac.RespBody, []byte(" index after ")...)
}

func (bac *BaseActionContext) Destruct() {
	println(" index destruct ")
}

type IndexController struct {
}

func (ic *IndexController) NewActionContext(req *http.Request, respWriter http.ResponseWriter) controller.ActionContext {
	return &BaseActionContext{
		Req:        req,
		RespWriter: respWriter,
	}
}

func (ic *IndexController) IndexAction(context *BaseActionContext) {
	context.RespBody = append(context.RespBody, []byte(" index action ")...)
}

func (ic *IndexController) RedirectAction(context *BaseActionContext) {
	print("here")
	system.Redirect302("https://github.com/goinbox")
}

type DemoActionContext struct {
	*BaseActionContext
}

func (dac *DemoActionContext) BeforeAction() {
	dac.RespBody = append(dac.RespBody, []byte(" demo before ")...)
}

func (dac *DemoActionContext) AfterAction() {
	dac.RespBody = append(dac.RespBody, []byte(" demo after ")...)
}

func (dac *DemoActionContext) Destruct() {
	println(" demo destruct ")
}

type DemoController struct {
}

func (dc *DemoController) NewActionContext(req *http.Request, respWriter http.ResponseWriter) controller.ActionContext {
	return &DemoActionContext{
		&BaseActionContext{
			Req:        req,
			RespWriter: respWriter,
		},
	}
}

func (dc *DemoController) DemoAction(context *DemoActionContext) {
	context.RespBody = append(context.RespBody, []byte(" demo action ")...)
}

func (dc *DemoController) GetAction(context *DemoActionContext, id string) {
	context.RespBody = append(context.RespBody, []byte(" get action id = "+id)...)
}
