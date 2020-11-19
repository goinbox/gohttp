package router

import (
	"fmt"
	"testing"
)

func TestSimpleRouter(t *testing.T) {
	r := NewSimpleRouter().SetParseRoutePathFunc(parseRoutePathFunc)

	r.FindRoute("/demo/get")
}

func parseRoutePathFunc(path string) *RouteGuide {
	fmt.Println(path)

	return nil
}
