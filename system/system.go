package system

import (
	"github.com/goinbox/gohttp/controller"
	"github.com/goinbox/gohttp/router"

	"net/http"
	"reflect"
)

type RoutePathFunc func(r *http.Request) string

type System struct {
	router router.Router

	rpf RoutePathFunc
}

func NewSystem(r router.Router) *System {
	s := &System{
		router: r,
	}

	s.rpf = s.routePath

	return s
}

func (s *System) SetRoutePathFunc(rpf RoutePathFunc) *System {
	s.rpf = rpf

	return s
}

func (s *System) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route := s.router.FindRoute(s.rpf(r))
	if route == nil {
		http.NotFound(w, r)
		return
	}

	context := route.Cl.NewActionContext(r, w)

	defer func() {
		if e := recover(); e != nil {
			ji, ok := e.(*jumpItem)
			if !ok {
				panic(e)
			}
			ji.jf(context, ji.args...)
		}

		_, _ = w.Write(context.ResponseBody())
		context.Destruct()
	}()

	context.BeforeAction()
	route.ActionValue.Call(s.makeArgsValues(context, route.Args))
	context.AfterAction()
}

func (s *System) routePath(r *http.Request) string {
	return r.URL.Path
}

func (s *System) makeArgsValues(context controller.ActionContext, args []string) []reflect.Value {
	argsValues := make([]reflect.Value, len(args)+1)
	argsValues[0] = reflect.ValueOf(context)
	for i, arg := range args {
		argsValues[i+1] = reflect.ValueOf(arg)
	}

	return argsValues
}

type JumpFunc func(context controller.ActionContext, args ...interface{})

type jumpItem struct {
	jf JumpFunc

	args []interface{}
}

func JumpOutAction(jf JumpFunc, args ...interface{}) {
	ji := &jumpItem{
		jf:   jf,
		args: args,
	}

	panic(ji)
}

func Redirect302(url string) {
	JumpOutAction(redirect302, url)
}

func redirect302(context controller.ActionContext, args ...interface{}) {
	http.Redirect(context.ResponseWriter(), context.Request(), args[0].(string), 302)
}
