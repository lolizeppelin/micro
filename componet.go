package micro

// Module is the interface that represent a module.
type Module interface {
	Init() error
	AfterInit()
	BeforeShutdown()
	Shutdown() error
}

/*
Component 通用api组件
Restful方法名 Get/List/Create/Update/Patch/Delete
非Restful方法名 以 GET/POST/PUT/PATCH/DELETE 开头
其他方法为网关可转发方法
*/
type Component interface {
	Init()
	AfterInit()
	BeforeShutdown()
	Shutdown()

	Name() string
	Collection() string
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
