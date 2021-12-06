package router

import (
	"reflect"
	"regexp"
	"strings"
)

type Route struct {
	C             Controller
	NewActionFunc reflect.Value
	Args          []string
}

type Router interface {
	MapRouteItems(cs ...Controller)
	DefineRouteItem(pattern string, c Controller, actionName string)

	FindRoute(path string) *Route
}

type RouteGuide struct {
	ControllerName string
	ActionName     string
	ActionArgs     []string
}

type ParseRoutePathFunc func(path string) *RouteGuide

type routeDefined struct {
	regex *regexp.Regexp

	controllerName string
	actionName     string
}

type routeItem struct {
	c  Controller
	cv reflect.Value
	ct reflect.Type

	controllerName string
	actionMap      map[string]reflect.Value
}

type router struct {
	emptyControllerName string
	emptyActionName     string

	defaultRouteGuide *RouteGuide

	cregex *regexp.Regexp
	aregex *regexp.Regexp

	routeDefinedList []*routeDefined
	routeTable       map[string]*routeItem

	prpf ParseRoutePathFunc
}

func NewRouter() *router {
	r := &router{
		emptyControllerName: "index",
		emptyActionName:     "index",

		cregex: regexp.MustCompile("([A-Z][A-Za-z0-9_]*)Controller$"),
		aregex: regexp.MustCompile("^([A-Z][A-Za-z0-9_]*)Action$"),

		routeTable: make(map[string]*routeItem),
	}

	r.prpf = r.parseRoutePath

	return r
}

func (r *router) parseRoutePath(path string) *RouteGuide {
	rg := new(RouteGuide)

	path = strings.Trim(path, "/")
	sl := strings.Split(path, "/")

	sl[0] = strings.TrimSpace(sl[0])
	if sl[0] == "" {
		rg.ControllerName = r.emptyControllerName
		rg.ActionName = r.emptyActionName
	} else {
		rg.ControllerName = sl[0]
		if len(sl) > 1 {
			sl[1] = strings.TrimSpace(sl[1])
			if sl[1] != "" {
				rg.ActionName = sl[1]
			} else {
				rg.ActionName = r.emptyActionName
			}
		} else {
			rg.ActionName = r.emptyActionName
		}
	}

	return rg
}

func (r *router) SetEmptyControllerName(name string) *router {
	r.emptyControllerName = name

	return r
}

func (r *router) SetEmptyActionName(name string) *router {
	r.emptyActionName = name

	return r
}

func (r *router) SetDefaultRoute(controllerName, actionName string) *router {
	r.defaultRouteGuide = &RouteGuide{
		ControllerName: controllerName,
		ActionName:     actionName,
	}

	return r
}

func (r *router) SetParseRoutePathFunc(f ParseRoutePathFunc) *router {
	r.prpf = f

	return r
}

func (r *router) MapRouteItems(cs ...Controller) {
	for _, c := range cs {
		r.mapRouteItem(c)
	}
}

func (r *router) mapRouteItem(c Controller) {
	ri := r.getRouteItem(c)
	if ri == nil {
		return
	}

	for i := 0; i < ri.cv.NumMethod(); i++ {
		m := ri.ct.Method(i)
		actionName := r.parseActionName(m.Name)
		if actionName == "" {
			continue
		}
		_, ok := ri.actionMap[actionName]
		if ok {
			continue
		}

		ri.actionMap[actionName] = ri.cv.Method(i)
	}
}

func (r *router) getRouteItem(c Controller) *routeItem {
	v := reflect.ValueOf(c)
	t := v.Type()
	controllerName := r.parseControllerName(t.String())
	if controllerName == "" {
		return nil
	}

	ri, ok := r.routeTable[controllerName]
	if !ok {
		ri = &routeItem{
			c:  c,
			cv: v,
			ct: t,

			controllerName: controllerName,
			actionMap:      make(map[string]reflect.Value),
		}
		r.routeTable[controllerName] = ri
	}

	return ri

}

func (r *router) parseControllerName(typeString string) string {
	matches := r.cregex.FindStringSubmatch(typeString)
	if matches == nil {
		return ""
	}

	return strings.ToLower(matches[1])
}

func (r *router) parseActionName(methodName string) string {
	matches := r.aregex.FindStringSubmatch(methodName)
	if matches == nil {
		return ""
	}

	return strings.ToLower(matches[1])
}

func (r *router) DefineRouteItem(pattern string, c Controller, actionName string) {
	if actionName == "" {
		return
	}

	ri := r.getRouteItem(c)
	if ri == nil {
		return
	}

	methodName := strings.Title(actionName) + "Action"
	actionName = strings.ToLower(actionName)

	ri.actionMap[actionName] = ri.cv.MethodByName(methodName)

	r.routeDefinedList = append(r.routeDefinedList, &routeDefined{
		regex: regexp.MustCompile(pattern),

		controllerName: ri.controllerName,
		actionName:     actionName,
	})
}

func (r *router) FindRoute(path string) *Route {
	path = strings.ToLower(path)

	rg := r.findRouteGuideByDefined(path)
	if rg == nil {
		rg = r.findRouteGuideByGeneral(path)
	}

	if rg == nil {
		return nil
	}

	ri, ok := r.routeTable[rg.ControllerName]
	if !ok {
		return nil
	}

	av, ok := ri.actionMap[rg.ActionName]
	if !ok {
		return nil
	}

	return &Route{
		C:             ri.c,
		NewActionFunc: av,
		Args:          rg.ActionArgs,
	}
}

func (r *router) findRouteGuideByDefined(path string) *RouteGuide {
	for _, d := range r.routeDefinedList {
		matches := d.regex.FindStringSubmatch(path)
		if matches == nil {
			continue
		}

		return &RouteGuide{
			ControllerName: d.controllerName,
			ActionName:     d.actionName,
			ActionArgs:     matches[1:],
		}
	}

	return nil
}

func (r *router) findRouteGuideByGeneral(path string) *RouteGuide {
	rg := r.prpf(path)
	if rg == nil {
		return nil
	}

	if r.checkIfUseDefaultRoute(rg) {
		return r.defaultRouteGuide
	}

	return rg
}

func (r *router) checkIfUseDefaultRoute(rg *RouteGuide) bool {
	if r.defaultRouteGuide == nil {
		return false
	}

	ri, ok := r.routeTable[rg.ControllerName]
	if !ok {
		return true
	}

	_, ok = ri.actionMap[rg.ActionName]
	if !ok {
		return true
	}

	return false
}
