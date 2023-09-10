package httpserver

import (
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

type redirectAction struct {
	baseAction

	code int
	url  string
}

func (a *redirectAction) Name() string {
	return "redirect"
}

func (a *redirectAction) Run(ctx *context) error {
	http.Redirect(a.ResponseWriter(), a.Request(), a.url, a.code)
	return nil
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

func (a *baseAction) redirect(code int, url string) {
	panic(&redirectAction{
		baseAction: *a,
		code:       code,
		url:        url,
	})
}

type indexAction struct {
	baseAction
}

func (a *indexAction) Name() string {
	return "index"
}

func (a *indexAction) Before(ctx *context) error {
	a.AppendResponseBody([]byte("before index\n"))
	return nil
}

func (a *indexAction) Run(ctx *context) error {
	a.AppendResponseBody([]byte("index action\n"))
	return nil
}

func (a *indexAction) After(ctx *context, err error) {
	a.AppendResponseBody([]byte("after index\n"))
}

type jumpAction struct {
	baseAction
}

func (a *jumpAction) Name() string {
	return "jump"
}

func (a *jumpAction) Before(ctx *context) error {
	a.AppendResponseBody([]byte("before jump\n"))
	return nil
}

func (a *jumpAction) Run(ctx *context) error {
	a.redirect(302, "https://github.com/goinbox")
	return nil
}
