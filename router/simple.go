package router

import (
	"github.com/goinbox/gohttp/controller"

	"reflect"
	"regexp"
	"strings"
)

type actionItem struct {
	argsNum int
	value   *reflect.Value
}

type routeItem struct {
	cl  controller.Controller
	clv *reflect.Value
	clt reflect.Type

	controllerName string
	actionMap      map[string]*actionItem
}

type routeDefined struct {
	regex *regexp.Regexp

	controllerName string
	actionName     string
}

type RouteGuide struct {
	ControllerName string
	ActionName     string
	ActionArgs     []string
}

type ParseRoutePathFunc func(path string) *RouteGuide

type SimpleRouter struct {
	emptyControllerName string
	emptyActionName     string

	defaultRouteGuide *RouteGuide

	cregex *regexp.Regexp
	aregex *regexp.Regexp

	routeDefinedList []*routeDefined
	routeTable       map[string]*routeItem

	prpf ParseRoutePathFunc
}

func NewSimpleRouter() *SimpleRouter {
	s := &SimpleRouter{
		emptyActionName:     "index",
		emptyControllerName: "index",

		cregex: regexp.MustCompile("([A-Z][A-Za-z0-9_]*)Controller$"),
		aregex: regexp.MustCompile("^([A-Z][A-Za-z0-9_]*)Action$"),

		routeTable: make(map[string]*routeItem),
	}

	s.prpf = s.parseRoutePath

	return s
}

func (s *SimpleRouter) SetEmptyControllerName(name string) *SimpleRouter {
	s.emptyControllerName = name

	return s
}

func (s *SimpleRouter) SetEmptyActionName(name string) *SimpleRouter {
	s.emptyActionName = name

	return s
}

func (s *SimpleRouter) SetDefaultRoute(controllerName, actionName string) *SimpleRouter {
	s.defaultRouteGuide = &RouteGuide{
		ControllerName: controllerName,
		ActionName:     actionName,
	}

	return s
}

func (s *SimpleRouter) SetParseRoutePathFunc(prpf ParseRoutePathFunc) *SimpleRouter {
	s.prpf = prpf

	return s
}

func (s *SimpleRouter) MapRouteItems(cls ...controller.Controller) {
	for _, cl := range cls {
		s.mapRouteItem(cl)
	}
}

func (s *SimpleRouter) mapRouteItem(cl controller.Controller) {
	ri := s.getRouteItem(cl)
	if ri == nil {
		return
	}

	for i := 0; i < ri.clv.NumMethod(); i++ {
		am := ri.clt.Method(i)
		actionName := s.getActionName(am.Name)
		if actionName == "" {
			continue
		}
		_, ok := ri.actionMap[actionName]
		if ok {
			continue
		}
		actionArgsNum := s.getActionArgsNum(am, ri.clt)
		if actionArgsNum == -1 {
			continue
		}

		av := ri.clv.Method(i)
		ri.actionMap[actionName] = &actionItem{
			argsNum: actionArgsNum,
			value:   &av,
		}
	}
}

func (s *SimpleRouter) DefineRouteItem(pattern string, cl controller.Controller, actionName string) {
	methodName := strings.Title(actionName) + "Action"
	actionName = strings.ToLower(methodName)
	if actionName == "" {
		return
	}

	ri := s.getRouteItem(cl)
	if ri == nil {
		return
	}

	am, ok := ri.clt.MethodByName(methodName)
	if !ok {
		return
	}
	actionArgsNum := s.getActionArgsNum(am, ri.clt)
	if actionArgsNum == -1 {
		return
	}

	av := ri.clv.MethodByName(methodName)
	ri.actionMap[actionName] = &actionItem{
		argsNum: actionArgsNum,
		value:   &av,
	}

	s.routeDefinedList = append(s.routeDefinedList, &routeDefined{
		regex: regexp.MustCompile(pattern),

		controllerName: strings.ToLower(ri.controllerName),
		actionName:     strings.ToLower(actionName),
	})
}

func (s *SimpleRouter) getRouteItem(cl controller.Controller) *routeItem {
	v := reflect.ValueOf(cl)
	t := v.Type()

	controllerName := s.getControllerName(t.String())
	if controllerName == "" {
		return nil
	}

	ri, ok := s.routeTable[controllerName]
	if !ok {
		ri = &routeItem{
			cl:  cl,
			clv: &v,
			clt: t,

			controllerName: controllerName,
			actionMap:      make(map[string]*actionItem),
		}
		s.routeTable[controllerName] = ri
	}

	return ri
}

func (s *SimpleRouter) getControllerName(typeString string) string {
	matches := s.cregex.FindStringSubmatch(typeString)
	if matches == nil {
		return ""
	}

	return strings.ToLower(matches[1])
}

func (s *SimpleRouter) getActionName(methodName string) string {
	matches := s.aregex.FindStringSubmatch(methodName)
	if matches == nil {
		return ""
	}

	actionName := strings.ToLower(matches[1])
	if actionName != "before" && actionName != "after" {
		return actionName
	}

	return ""
}

func (s *SimpleRouter) getActionArgsNum(actionMethod reflect.Method, controllerType reflect.Type) int {
	n := actionMethod.Type.NumIn()
	if n < 2 {
		return -1
	}

	if actionMethod.Type.In(0).String() != controllerType.String() {
		return -1
	}

	if n > 2 {
		valid := true
		for i := 2; i < n; i++ {
			if actionMethod.Type.In(i).String() != "string" {
				valid = false
				break
			}
		}
		if !valid {
			return -1
		}
	}

	return n - 2 //delete s and context
}

func (s *SimpleRouter) FindRoute(path string) *Route {
	path = strings.ToLower(path)

	rg := s.findRouteGuideByDefined(path)
	if rg == nil {
		rg = s.findRouteGuideByGeneral(path)
	}

	if rg == nil {
		return nil
	}

	ri, ok := s.routeTable[rg.ControllerName]
	if !ok {
		return nil
	}

	ai, ok := ri.actionMap[rg.ActionName]
	if !ok {
		return nil
	}

	return &Route{
		Cl:          ri.cl,
		ActionValue: ai.value,
		Args:        s.makeActionArgs(rg.ActionArgs, ai.argsNum),
	}
}

func (s *SimpleRouter) findRouteGuideByDefined(path string) *RouteGuide {
	for _, rd := range s.routeDefinedList {
		matches := rd.regex.FindStringSubmatch(path)
		if matches == nil {
			continue
		}

		return &RouteGuide{
			ControllerName: rd.controllerName,
			ActionName:     rd.actionName,
			ActionArgs:     matches[1:],
		}
	}

	return nil
}

func (s *SimpleRouter) findRouteGuideByGeneral(path string) *RouteGuide {
	rg := s.prpf(path)
	if rg == nil {
		return nil
	}

	if s.checkIfUseDefaultRoute(rg) {
		return s.defaultRouteGuide
	}

	return rg
}

func (s *SimpleRouter) parseRoutePath(path string) *RouteGuide {
	rg := new(RouteGuide)

	path = strings.Trim(path, "/")
	sl := strings.Split(path, "/")

	sl[0] = strings.TrimSpace(sl[0])
	if sl[0] == "" {
		rg.ControllerName = s.emptyControllerName
		rg.ActionName = s.emptyActionName
	} else {
		rg.ControllerName = sl[0]
		if len(sl) > 1 {
			sl[1] = strings.TrimSpace(sl[1])
			if sl[1] != "" {
				rg.ActionName = sl[1]
			} else {
				rg.ActionName = s.emptyActionName
			}
		} else {
			rg.ActionName = s.emptyActionName
		}
	}

	return rg
}

func (s *SimpleRouter) makeActionArgs(args []string, validArgsNum int) []string {
	rgArgsNum := len(args)
	missArgsNum := validArgsNum - rgArgsNum
	switch {
	case missArgsNum == 0:
	case missArgsNum > 0:
		for i := 0; i < missArgsNum; i++ {
			args = append(args, "")
		}
	case missArgsNum < 0:
		args = args[:validArgsNum]
	}

	return args
}

func (s *SimpleRouter) checkIfUseDefaultRoute(rg *RouteGuide) bool {
	if s.defaultRouteGuide == nil {
		return false
	}

	ri, ok := s.routeTable[rg.ControllerName]
	if !ok {
		return true
	}

	_, ok = ri.actionMap[rg.ActionName]
	if !ok {
		return true
	}

	return false
}
