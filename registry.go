package micro

import "context"

type Registry interface {
	Register(context.Context, *Service) error
	Deregister(context.Context, *Service) error
	GetService(service string) ([]*Service, error)
	ListServices() ([]*Service, error)
	Watch(service string) (Watcher, error)
	String() string
}

type Watcher interface {
	// Next is a blocking call
	Next() (*Result, error)
	Stop()
}

type Result struct {
	Service *Service
	Action  string
}

type Service struct {
	Name      string            `json:"name"`
	Version   int               `json:"version"` // 服务主版本号
	Metadata  map[string]string `json:"metadata"`
	Nodes     []*Node           `json:"nodes"`
	Endpoints []*Endpoint       `json:"endpoints"`
}

type Node struct {
	Id       string            `json:"id"`
	Version  Version           `json:"version"` // 节点版本号
	Max      *Version          `json:"max"`     // 节点版本兼容上限
	Min      *Version          `json:"min"`     // 节点版本兼容下限
	Address  string            `json:"address"`
	Metadata map[string]string `json:"metadata"`
}

type Endpoint struct {
	Name       string            `json:"name"`
	Metadata   map[string]string `json:"metadata"`           // 元数据
	PrimaryKey bool              `json:"pk,omitempty"`       // 是否需要主键
	Internal   bool              `json:"internal,omitempty"` // 是否内部rpc
}
