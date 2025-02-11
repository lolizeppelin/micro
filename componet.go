package micro

import (
	"context"
	"golang.org/x/exp/slices"
	"net/url"
)

var (
	components []Component
)

type PreExecuteHook func(context.Context, url.Values, []byte) (context.Context, error)

// Module is the interface that represent a module.
type Module interface {
	Init() error
	AfterInit()
	BeforeShutdown()
	Shutdown() error
}

/*
Component 通用api组件
1. Restful方法名 Get/List/Create/Update/Patch/Delete
2. 非Restful方法名 以 GET_/POST_/PUT_/PATCH_/DELETE_开头,其余部分小写为路径  e.g  User.GET_Money, 路径为/user/money
3. 以RPC_开头的方法为内部rpc方法
4. 其他方法为注册到网关可转发方法(不可与3分割后同名)
5. 一般不建议在组件中设置生命周期方法,尽量在模块中做生命周期相关操作
*/
type Component interface {
	/*
		Init 初始化执行
	*/
	Init()
	/*
		AfterInit 初始化完成后执行
	*/
	AfterInit()
	/*
		BeforeShutdown 进程停止前时执行
	*/
	BeforeShutdown()
	/*
		Shutdown 进程停止时执行
	*/
	Shutdown()
	/*
		Name 组件对应Resource名
	*/
	Name() string
	/*
		Collection 组件对应Collection名(resource的复数形式)
	*/
	Collection() string
	/*
		Hooks pre execute hook, nil able
		@method  原始方法名
	*/
	Hooks(method string) []PreExecuteHook
}

/*
ComponentBase 通用组件继承
*/
type ComponentBase struct{}

func (*ComponentBase) Init()                           {}
func (*ComponentBase) AfterInit()                      {}
func (*ComponentBase) BeforeShutdown()                 {}
func (*ComponentBase) Shutdown()                       {}
func (*ComponentBase) Hooks(_ string) []PreExecuteHook { return nil }

func RegComponent(component Component) {
	components = append(components, component)
}

func LoadComponents() []Component {
	return slices.Clone(components)
}
