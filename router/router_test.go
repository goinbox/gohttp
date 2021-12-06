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

func (c *DemoController) AddAction() *addAction {
	return new(addAction)
}

func (c *DemoController) DelAction() *delAction {
	return new(delAction)
}

func (c *DemoController) EditAction() *editAction {
	return new(editAction)
}

func (c *DemoController) GetAction() *getAction {
	return new(getAction)
}

type baseAction interface {
	Name() string
	Run()
}

type addAction struct {
}

func (a *addAction) Name() string {
	return "add"
}

func (a *addAction) Run() {
	fmt.Println("run add")
}

type delAction struct {
}

func (a *delAction) Name() string {
	return "del"
}

func (a *delAction) Run() {
	fmt.Println("run del")
}

type editAction struct {
}

func (a *editAction) Name() string {
	return "edit"
}

func (a *editAction) Run() {
	fmt.Println("run edit")
}

type getAction struct {
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
		action := vs[0].Interface().(baseAction)
		t.Log("action", action.Name())
		action.Run()
	}
}
