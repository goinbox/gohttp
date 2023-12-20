package httpserver

import (
	"github.com/goinbox/pcontext"
)

type Context interface {
	pcontext.Context

	SetController(controller string)
	Controller() string

	SetAction(action string)
	Action() string
}

type BaseContext struct {
	pcontext.Context

	controller string
	action     string
}

func (c *BaseContext) SetController(controller string) {
	c.controller = controller
}

func (c *BaseContext) Controller() string {
	return c.controller
}

func (c *BaseContext) SetAction(action string) {
	c.action = action
}

func (c *BaseContext) Action() string {
	return c.action
}
