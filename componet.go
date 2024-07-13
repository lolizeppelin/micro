package micro

// Module is the interface that represent a module.
type Module interface {
	Init() error
	AfterInit()
	BeforeShutdown()
	Shutdown() error
}

type Component interface {
	Init()
	AfterInit()
	BeforeShutdown()
	Shutdown()

	Name() string
	Collection() string
	Internal() bool
}

type ComponentBase struct {
}

func (*ComponentBase) Init() {

}
func (*ComponentBase) AfterInit() {

}
func (*ComponentBase) BeforeShutdown() {

}
func (*ComponentBase) Shutdown() {

}
func (*ComponentBase) Internal() bool {
	return true
}
