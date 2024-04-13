package micro

type Registry interface {
	Register(*Service) error
	Deregister(*Service) error
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
	Version   string            `json:"version"` // 当前版本号
	Metadata  map[string]string `json:"metadata"`
	Nodes     []*Node           `json:"nodes"`
	Endpoints []*Endpoint       `json:"endpoints"`
}

type Node struct {
	Id       string            `json:"id"`
	Version  string            `json:"version"` // 节点版本号
	Address  string            `json:"address"`
	Metadata map[string]string `json:"metadata"`
}

type Endpoint struct {
	Name     string            `json:"name"`
	Internal bool              `json:"internal,omitempty"` // 是否内部方法(不可通过外部网关转发)
	Metadata map[string]string `json:"metadata"`           // 元数据
}
