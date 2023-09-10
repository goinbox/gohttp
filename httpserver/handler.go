package httpserver

import (
	"net/http"

	"github.com/goinbox/pcontext"
	"github.com/goinbox/router"
)

type RoutePathFunc func(r *http.Request) string

type handler[T pcontext.Context] struct {
	router router.Router

	rpf RoutePathFunc
}

func NewHandler[T pcontext.Context](r router.Router) http.Handler {
	s := &handler[T]{
		router: r,
	}

	s.rpf = s.routePath

	return s
}

func (h *handler[T]) SetRoutePathFunc(f RoutePathFunc) *handler[T] {
	h.rpf = f

	return h
}

func (h *handler[T]) routePath(r *http.Request) string {
	return r.URL.Path
}

func (h *handler[T]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route := h.router.FindRoute(h.rpf(r))
	if route == nil {
		http.NotFound(w, r)
		return
	}

	action := route.NewActionFunc.Call(nil)[0].Interface().(Action[T])
	ctx := action.Init(r, w, route.Args)

	defer func() {
		if e := recover(); e != nil {
			var ok bool
			action, ok = e.(Action[T])
			if !ok {
				panic(e)
			}

			h.runAction(ctx, action)
		}

		_, _ = w.Write(action.ResponseBody())
		action.Destruct(ctx)
	}()

	h.runAction(ctx, action)
}

func (h *handler[T]) runAction(ctx T, action Action[T]) {
	action.Before(ctx)
	action.Run(ctx)
	action.After(ctx)
}
