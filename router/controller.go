package router

type Controller interface {
	Name() string
}

type Action interface {
	Name() string

	Before()
	Run()
	After()
	Destruct()
}

type BaseAction struct {
}

func (b BaseAction) Before() {
}

func (b BaseAction) After() {
}

func (b BaseAction) Destruct() {
}
