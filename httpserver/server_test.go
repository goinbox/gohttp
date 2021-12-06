package httpserver

import (
	"net/http"
	"testing"

	"github.com/goinbox/gohttp/gracehttp"
)

type IndexController struct {
}

func (c *IndexController) Name() string {
	return "index"
}

func (c *IndexController) IndexAction(r *http.Request, w http.ResponseWriter, args []string) *IndexAction {
	return &IndexAction{NewBaseAction(r, w, args)}
}

func (c *IndexController) JumpAction(r *http.Request, w http.ResponseWriter, args []string) *JumpAction {
	return &JumpAction{NewBaseAction(r, w, args)}
}

type IndexAction struct {
	*BaseAction
}

func (a *IndexAction) Name() string {
	return "index"
}

func (a *IndexAction) Run() {
	a.AppendResponseBody([]byte("index action"))
}

type JumpAction struct {
	*BaseAction
}

func (a *JumpAction) Name() string {
	return "jump"
}

func (a *JumpAction) Run() {
	Redirect(a.Request(), a.ResponseWriter(), 302, "https://github.com/goinbox")
}

func TestServer(t *testing.T) {
	r := NewRouter()
	r.MapRouteItems(new(IndexController))

	_ = gracehttp.ListenAndServe(":8010", NewServer(r))

}
