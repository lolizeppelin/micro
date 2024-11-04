package micro

import (
	"context"
	"net/url"
)

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
*/
type Component interface {
	Init()
	AfterInit()
	BeforeShutdown()
	Shutdown()

	Name() string
	Collection() string

	// Hooks pre execute hook, nil able
	Hooks(method string) []func(context.Context, url.Values, any) (context.Context, error)
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
func (*ComponentBase) Hooks(method string) []func(context.Context, url.Values, any) (context.Context, error) {
	return nil
}
