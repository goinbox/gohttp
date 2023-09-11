package httpserver

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goinbox/router"
	"go.opentelemetry.io/otel/trace"
)

func startTrace(ctx *context, spanName string, opts ...trace.SpanStartOption) (*context, trace.Span) {
	fmt.Println("---------- start trace", spanName)
	_, span := trace.NewNoopTracerProvider().Tracer("").Start(ctx, spanName, opts...)
	return ctx, span
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
	r.MapRouteItems(new(indexController))

	for _, path := range []string{"index", "jump"} {
		header, content, err := runHandler(NewHandler[*context](r).SetStartTraceFunc(startTrace),
			fmt.Sprintf("http://127.0.0.1/index/%s", path))
		t.Log(path, err, header, string(content))
	}
}
