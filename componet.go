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
	Internal() bool
}
