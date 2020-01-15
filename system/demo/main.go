package main

import (
	"net/http"

	"github.com/goinbox/gohttp/controller"
	"github.com/goinbox/gohttp/gracehttp"
	"github.com/goinbox/gohttp/router"
	"github.com/goinbox/gohttp/system"
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
	*controller.BaseContext
}

func (bac *BaseActionContext) BeforeAction() {
	bac.AppendResponseBody([]byte(" index before "))
}

func (bac *BaseActionContext) AfterAction() {
	bac.AppendResponseBody([]byte(" index after "))
}

func (bac *BaseActionContext) Destruct() {
	println(" index destruct ")
}

type IndexController struct {
}

func (ic *IndexController) NewActionContext(req *http.Request, respWriter http.ResponseWriter) controller.ActionContext {
	return &BaseActionContext{
		controller.NewBaseContext(req, respWriter),
	}
}

func (ic *IndexController) IndexAction(context *BaseActionContext) {
	context.AppendResponseBody([]byte(" index action "))
}

func (ic *IndexController) RedirectAction(context *BaseActionContext) {
	print("here")
	system.Redirect302("https://github.com/goinbox")
}

type DemoActionContext struct {
	*BaseActionContext
}

func (dac *DemoActionContext) BeforeAction() {
	dac.AppendResponseBody([]byte(" demo before "))
}

func (dac *DemoActionContext) AfterAction() {
	dac.AppendResponseBody([]byte(" demo after "))
}

func (dac *DemoActionContext) Destruct() {
	println(" demo destruct ")
}

type DemoController struct {
}

func (dc *DemoController) NewActionContext(req *http.Request, respWriter http.ResponseWriter) controller.ActionContext {
	return &DemoActionContext{
		&BaseActionContext{
			controller.NewBaseContext(req, respWriter),
		},
	}
}

func (dc *DemoController) DemoAction(context *DemoActionContext) {
	context.AppendResponseBody([]byte(" demo action "))
}

func (dc *DemoController) GetAction(context *DemoActionContext, id string) {
	context.AppendResponseBody([]byte(" get action id = " + id))
}
