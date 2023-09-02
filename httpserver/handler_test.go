package httpserver

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goinbox/gohttp/router"
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

func (a *IndexAction) Before() {
	a.AppendResponseBody([]byte("before\n"))
}

func (a *IndexAction) Run() {
	a.AppendResponseBody([]byte("index action\n"))
}

func (a *IndexAction) After() {
	a.AppendResponseBody([]byte("after\n"))
}

func (a *IndexAction) Destruct() {
	fmt.Println("destruct")
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

func runHandler(handler http.Handler, target string) (http.Header, []byte, error) {
	req := httptest.NewRequest(http.MethodPost, target, nil)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	resp := w.Result()
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	return resp.Header, body, err
}

func TestHandler(t *testing.T) {
	r := router.NewRouter()
	r.MapRouteItems(new(IndexController))

	handler := NewHandler(r)
	for _, path := range []string{"index", "jump"} {
		header, content, err := runHandler(handler, fmt.Sprintf("http://127.0.0.1/index/%s", path))
		t.Log(path, err, header, string(content))
	}
}
