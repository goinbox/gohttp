package httpserver

import (
	"net/http"
	"reflect"
)

type RoutePathFunc func(r *http.Request) string

type Server struct {
	router Router

	rpf RoutePathFunc
}

func NewServer(r Router) *Server {
	s := &Server{
		router: r,
	}

	s.rpf = s.routePath

	return s
}

func (s *Server) SetRoutePathFunc(f RoutePathFunc) *Server {
	s.rpf = f

	return s
}

func (s *Server) routePath(r *http.Request) string {
	return r.URL.Path
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	action.BeforeRun()
	action.Run()
	action.AfterRun()
}

func (s *Server) makeArgsValues(r *http.Request, w http.ResponseWriter, args []string) []reflect.Value {
	return []reflect.Value{
		reflect.ValueOf(r),
		reflect.ValueOf(w),
		reflect.ValueOf(args),
	}
}

type redirectAction struct {
	*BaseAction

	code int
	url  string
}

func (a *redirectAction) Name() string {
	return "redirect"
}

func (a *redirectAction) Run() {
	http.Redirect(a.ResponseWriter(), a.Request(), a.url, a.code)
}

func Redirect(r *http.Request, w http.ResponseWriter, code int, url string) {
	panic(&redirectAction{
		BaseAction: NewBaseAction(r, w, nil),
		code:       code,
		url:        url,
	})
}
