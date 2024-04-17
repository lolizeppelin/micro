package registry

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lolizeppelin/micro"
	"github.com/lolizeppelin/micro/log"
	"go.uber.org/zap"

	"golang.org/x/exp/maps"
	"net"
	"path"
	"sort"
	"strings"
	"sync"
	"time"

	hash "github.com/mitchellh/hashstructure/v2"
	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
	clientV3 "go.etcd.io/etcd/client/v3"
)

var (
	prefix               = "/micro/registry/"
	DefaultRegisterCheck = func(ctx context.Context) error { return nil }
)

type etcdRegistry struct {
	client  *clientV3.Client
	options Options

	sync.RWMutex
	register map[string]uint64
	leases   map[string]clientV3.LeaseID
}

func NewEtcdRegistry(opts ...Option) (micro.Registry, error) {
	e := &etcdRegistry{
		options:  Options{},
		register: make(map[string]uint64),
		leases:   make(map[string]clientV3.LeaseID),
	}
	err := configure(e, opts...)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func configure(e *etcdRegistry, opts ...Option) error {

	config := clientV3.Config{
		Endpoints: []string{"127.0.0.1:2379"},
		Logger:    zap.NewNop(),
	}
	for _, o := range opts {
		o(&e.options)
	}
	if e.options.Timeout == 0 {
		e.options.Timeout = 5 * time.Second
	}
	config.DialTimeout = e.options.Timeout

	if e.options.TLSConfig != nil {
		tlsConfig := e.options.TLSConfig
		if tlsConfig == nil {
			tlsConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}
		config.TLS = tlsConfig
	}

	var cAddress []string

	for _, address := range e.options.Address {
		if len(address) == 0 {
			continue
		}
		addr, port, err := net.SplitHostPort(address)
		var ae *net.AddrError
		if errors.As(err, &ae) && ae.Err == "missing port in address" {
			port = "2379"
			addr = address
			cAddress = append(cAddress, net.JoinHostPort(addr, port))
		}
	}

	if len(cAddress) > 0 {
		config.Endpoints = cAddress
	}

	cli, err := clientV3.New(config)
	if err != nil {
		return err
	}
	//cli.Get()

	e.client = cli
	return nil
}

func encode(s *micro.Service) string {
	b, _ := json.Marshal(s)
	return string(b)
}

func decode(ds []byte) *micro.Service {
	var s *micro.Service
	_ = json.Unmarshal(ds, &s)
	return s
}

func nodePath(s, id string) string {
	service := strings.Replace(s, "/", "-", -1)
	node := strings.Replace(id, "/", "-", -1)
	return path.Join(prefix, service, node)
}

func servicePath(s string) string {
	return path.Join(prefix, strings.Replace(s, "/", "-", -1))
}

func (e *etcdRegistry) Client() *clientV3.Client {
	return e.client
}

func (e *etcdRegistry) registerNode(s *micro.Service, node *micro.Node) error {
	if len(s.Nodes) == 0 {
		return errors.New("require at least one node")
	}

	// check existing lease cache
	e.RLock()
	leaseID, ok := e.leases[s.Name+node.Id]
	e.RUnlock()

	if !ok {
		// missing lease, check if the key exists
		ctx, cancel := context.WithTimeout(context.Background(), e.options.Timeout)
		defer cancel()

		// look for the existing key
		rsp, err := e.client.Get(ctx, nodePath(s.Name, node.Id), clientV3.WithSerializable())
		if err != nil {
			return err
		}

		// get the existing lease
		for _, kv := range rsp.Kvs {
			if kv.Lease > 0 {
				leaseID = clientV3.LeaseID(kv.Lease)

				// decode the existing node
				srv := decode(kv.Value)
				if srv == nil || len(srv.Nodes) == 0 {
					continue
				}

				// create hash of service; uint64
				h, err := hash.Hash(srv.Nodes[0], hash.FormatV2, nil)
				if err != nil {
					continue
				}

				// save the info
				e.Lock()
				e.leases[s.Name+node.Id] = leaseID
				e.register[s.Name+node.Id] = h
				e.Unlock()

				break
			}
		}
	}

	var leaseNotFound bool

	// renew the lease if it exists
	if leaseID > 0 {
		log.Tracef("Renewing existing lease for %s %d", s.Name, leaseID)
		if _, err := e.client.KeepAliveOnce(context.TODO(), leaseID); err != nil {
			if !errors.Is(err, rpctypes.ErrLeaseNotFound) {
				return err
			}

			log.Tracef("Lease not found for %s %d", s.Name, leaseID)
			// lease not found do register
			leaseNotFound = true
		}
	}

	// create hash of service; uint64
	h, err := hash.Hash(node, hash.FormatV2, nil)
	if err != nil {
		return err
	}

	// get existing hash for the service node
	e.Lock()
	v, ok := e.register[s.Name+node.Id]
	e.Unlock()

	// the service is unchanged, skip registering
	if ok && v == h && !leaseNotFound {
		log.Tracef("Service %s node %s unchanged skipping registration", s.Name, node.Id)
		return nil
	}

	service := &micro.Service{
		Name:      s.Name,
		Version:   s.Version,
		Metadata:  s.Metadata,
		Endpoints: s.Endpoints,
		Nodes:     []*micro.Node{node},
	}

	ctx, cancel := context.WithTimeout(context.Background(), e.options.Timeout)
	defer cancel()

	var lgr *clientV3.LeaseGrantResponse
	if e.options.TTL.Seconds() > 0 {
		// get a lease used to expire keys since we have a ttl
		lgr, err = e.client.Grant(ctx, int64(e.options.TTL.Seconds()))
		if err != nil {
			return err
		}
		log.Tracef("Registering %s id %s with lease %v and leaseID %v and ttl %v",
			service.Name, node.Id, lgr, lgr.ID, e.options.TTL)

	}
	// create an entry for the node
	if lgr != nil {
		_, err = e.client.Put(ctx, nodePath(service.Name, node.Id), encode(service), clientV3.WithLease(lgr.ID))
	} else {
		_, err = e.client.Put(ctx, nodePath(service.Name, node.Id), encode(service))
	}
	if err != nil {
		return err
	}

	e.Lock()
	// save our hash of the service
	e.register[s.Name+node.Id] = h
	// save our leaseID of the service
	if lgr != nil {
		e.leases[s.Name+node.Id] = lgr.ID
	}
	e.Unlock()

	return nil
}

func (e *etcdRegistry) Deregister(s *micro.Service) error {
	if len(s.Nodes) == 0 {
		return errors.New("require at least one node")
	}
	for _, node := range s.Nodes {
		e.Lock()
		// delete our hash of the service
		delete(e.register, s.Name+node.Id)
		// delete our lease of the service
		delete(e.leases, s.Name+node.Id)
		e.Unlock()

		ctx, cancel := context.WithTimeout(context.Background(), e.options.Timeout)
		defer cancel()

		log.Tracef("deregister %s id %s", s.Name, node.Id)
		_, err := e.client.Delete(ctx, nodePath(s.Name, node.Id))
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *etcdRegistry) Register(s *micro.Service) error {
	if len(s.Nodes) == 0 {
		return errors.New("require at least one node")
	}

	var _err error

	// register each node individually
	for _, node := range s.Nodes {
		err := e.registerNode(s, node)
		if err != nil {
			_err = err
		}
	}

	return _err
}

func (e *etcdRegistry) GetService(name string) ([]*micro.Service, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.options.Timeout)
	defer cancel()

	rsp, err := e.client.Get(ctx, servicePath(name)+"/", clientV3.WithPrefix(), clientV3.WithSerializable())
	if err != nil {
		return nil, err
	}

	if len(rsp.Kvs) == 0 {
		return nil, micro.ErrServiceNotFound
	}

	serviceMap := map[string]*micro.Service{}

	for _, n := range rsp.Kvs {
		if sn := decode(n.Value); sn != nil {
			s, ok := serviceMap[sn.Version]
			if !ok {
				s = &micro.Service{
					Name:      sn.Name,
					Version:   sn.Version,
					Metadata:  sn.Metadata,
					Endpoints: sn.Endpoints,
				}
				serviceMap[s.Version] = s
			}

			s.Nodes = append(s.Nodes, sn.Nodes...)
		}
	}

	return maps.Values(serviceMap), nil
}

func (e *etcdRegistry) ListServices() ([]*micro.Service, error) {
	versions := make(map[string]*micro.Service)

	ctx, cancel := context.WithTimeout(context.Background(), e.options.Timeout)
	defer cancel()

	rsp, err := e.client.Get(ctx, prefix, clientV3.WithPrefix(), clientV3.WithSerializable())
	if err != nil {
		return nil, err
	}

	if len(rsp.Kvs) == 0 {
		return []*micro.Service{}, nil
	}

	for _, n := range rsp.Kvs {
		sn := decode(n.Value)
		if sn == nil {
			continue
		}
		key := fmt.Sprintf("%s-%d", sn.Name, sn.Version)
		v, ok := versions[key]
		if !ok {
			versions[key] = sn
			continue
		}
		// append to service:version nodes
		v.Nodes = append(v.Nodes, sn.Nodes...)
	}

	services := make([]*micro.Service, 0, len(versions))
	for _, service := range versions {
		services = append(services, service)
	}

	// sort the services
	sort.Slice(services, func(i, j int) bool { return services[i].Name < services[j].Name })

	return services, nil
}

func (e *etcdRegistry) Watch(service string) (micro.Watcher, error) {
	return newEtcdWatcher(e, e.options.Timeout, service)
}

func (e *etcdRegistry) String() string {
	return "etcd"
}
