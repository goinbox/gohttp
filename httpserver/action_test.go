package httpserver

import (
	"fmt"
	"net/http"

	"github.com/goinbox/golog"
	"github.com/goinbox/pcontext"
)

type indexController struct {
}

func (c *indexController) Name() string {
	return "index"
}

func (c *indexController) IndexAction() *indexAction {
	return &indexAction{}
}

func (c *indexController) JumpAction() *jumpAction {
	return &jumpAction{}
}

type context struct {
	pcontext.Context
}

type baseAction struct {
	BaseAction[*context]
}

func (a *baseAction) Init(r *http.Request, w http.ResponseWriter, args []string) *context {
	a.BaseAction.Init(r, w, args)
	return &context{pcontext.NewSimpleContext(&golog.NoopLogger{})}
}

type indexAction struct {
	baseAction
}

func (a *indexAction) Name() string {
	return "index"
}

func (a *indexAction) Before(ctx *context) {
	a.AppendResponseBody([]byte("before index\n"))
}

func (a *indexAction) Run(ctx *context) {
	a.AppendResponseBody([]byte("index action\n"))
}

func (a *indexAction) After(ctx *context) {
	a.AppendResponseBody([]byte("after index\n"))
}

func (a *indexAction) Destruct(ctx *context) {
	fmt.Println("destruct index")
}

type jumpAction struct {
	baseAction
}

func (a *jumpAction) Name() string {
	return "jump"
}

func (a *jumpAction) Before(ctx *context) {
	a.AppendResponseBody([]byte("before jump\n"))
}

func (a *jumpAction) Run(ctx *context) {
	redirect(302, "https://github.com/goinbox")
}

func (a *jumpAction) After(ctx *context) {
	a.AppendResponseBody([]byte("after jump\n"))
}

func (a *jumpAction) Destruct(ctx *context) {
	fmt.Println("destruct jump")
}

type redirectAction struct {
	baseAction

	code int
	url  string
}

func (a *redirectAction) Name() string {
	return "redirect"
}

func (a *redirectAction) Run(ctx *context) {
	http.Redirect(a.ResponseWriter(), a.Request(), a.url, a.code)
}

func redirect(code int, url string) {
	panic(&redirectAction{
		code: code,
		url:  url,
	})
}
