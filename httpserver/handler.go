package httpserver

import (
	"net/http"
	"reflect"

	"github.com/goinbox/gohttp/router"
)

type RoutePathFunc func(r *http.Request) string

type handler struct {
	router router.Router

	rpf RoutePathFunc
}

func NewHandler(r router.Router) http.Handler {
	s := &handler{
		router: r,
	}

	s.rpf = s.routePath

	return s
}

func (s *handler) SetRoutePathFunc(f RoutePathFunc) *handler {
	s.rpf = f

	return s
}

func (s *handler) routePath(r *http.Request) string {
	return r.URL.Path
}

func (s *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route := s.router.FindRoute(s.rpf(r))
	if route == nil {
		http.NotFound(w, r)
		return
	}

	action := route.NewActionFunc.Call(s.makeArgsValues(r, w, route.Args))[0].Interface().(Action)

	defer func() {
		if e := recover(); e != nil {
			a, ok := e.(Action)
			if !ok {
				panic(e)
			}
			a.Run()
		}

		_, _ = w.Write(action.ResponseBody())
		action.Destruct()
	}()

	action.Before()
	action.Run()
	action.After()
}

func (s *handler) makeArgsValues(r *http.Request, w http.ResponseWriter, args []string) []reflect.Value {
	return []reflect.Value{
		reflect.ValueOf(r),
		reflect.ValueOf(w),
		reflect.ValueOf(args),
	}
}
