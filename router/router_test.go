package router

import (
	"fmt"
	"testing"
)

type DemoController struct {
}

func (c *DemoController) Name() string {
	return "demo"
}

func (c *DemoController) AddAction() Action {
	return new(addAction)
}

func (c *DemoController) DelAction() Action {
	return new(delAction)
}

func (c *DemoController) EditAction() Action {
	return new(editAction)
}

func (c *DemoController) GetAction() Action {
	return new(getAction)
}

type addAction struct {
	BaseAction
}

func (a *addAction) Name() string {
	return "add"
}

func (a *addAction) Run() {
	fmt.Println("run add")
}

type delAction struct {
	BaseAction
}

func (a *delAction) Name() string {
	return "del"
}

func (a *delAction) Run() {
	fmt.Println("run del")
}

type editAction struct {
	BaseAction
}

func (a *editAction) Name() string {
	return "edit"
}

func (a *editAction) Run() {
	fmt.Println("run edit")
}

type getAction struct {
	BaseAction
}

func (a *getAction) Name() string {
	return "get"
}

func (a *getAction) Run() {
	fmt.Println("run get")
}

func TestRouter(t *testing.T) {
	r := NewRouter()
	r.MapRouteItems(new(DemoController))

	pathList := []string{
		"demo/add",
		"demo/del",
		"demo/edit",
		"demo/get",
	}

	for _, path := range pathList {
		route := r.FindRoute(path)
		if route == nil {
			t.Fatal(path, "not find route")
		}
		t.Log("controller", route.C.Name())

		vs := route.NewActionFunc.Call(nil)
		action := vs[0].Interface().(Action)
		t.Log("action", action.Name())
		action.Run()
	}
}
