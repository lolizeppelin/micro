package registry

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lolizeppelin/micro"
	"github.com/lolizeppelin/micro/log"
	"github.com/lolizeppelin/micro/utils"
	"golang.org/x/exp/maps"
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
	options := Options{
		Timeout: 5 * time.Second,
	}

	for _, o := range opts {
		o(&options)
	}

	if options.Client == nil {
		return nil, fmt.Errorf("etcd v3 client not found")
	}

	e := &etcdRegistry{
		options:  options,
		register: make(map[string]uint64),
		leases:   make(map[string]clientV3.LeaseID),
		client:   options.Client,
	}
	return e, nil
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

func (e *etcdRegistry) registerNode(ctx context.Context, s *micro.Service, node *micro.Node) error {
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
		log.Debugf(ctx, "Renewing existing lease for %s %d", s.Name, leaseID)
		if _, err := e.client.KeepAliveOnce(context.TODO(), leaseID); err != nil {
			if !errors.Is(err, rpctypes.ErrLeaseNotFound) {
				return err
			}

			log.Debugf(ctx, "Lease not found for %s %d", s.Name, leaseID)
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
		log.Debugf(ctx, "Service %s node %s unchanged skipping registration", s.Name, node.Id)
		return nil
	}

	service := &micro.Service{
		Name:      s.Name,
		Version:   s.Version,
		Metadata:  s.Metadata,
		Endpoints: s.Endpoints,
		Nodes:     []*micro.Node{node},
	}

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, e.options.Timeout)
	defer cancel()

	var lgr *clientV3.LeaseGrantResponse
	if e.options.TTL.Seconds() > 0 {
		// get a lease used to expire keys since we have a ttl
		lgr, err = e.client.Grant(ctx, int64(e.options.TTL.Seconds()))
		if err != nil {
			return err
		}
		log.Debugf(ctx, "Registering %s id %s with lease %v and leaseID %v and ttl %v",
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

func (e *etcdRegistry) Deregister(ctx context.Context, s *micro.Service) error {
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

		log.Debugf(ctx, "deregister %s id %s", s.Name, node.Id)

		f := func() error {
			_ctx, cancel := context.WithTimeout(ctx, e.options.Timeout)
			defer cancel()
			_, err := e.client.Delete(_ctx, nodePath(s.Name, node.Id))
			if err != nil {
				return err
			}
			return nil
		}

		if err := f(); err != nil {
			return err
		}

	}

	return nil
}

func (e *etcdRegistry) Register(ctx context.Context, s *micro.Service) error {
	if len(s.Nodes) == 0 {
		return errors.New("require at least one node")
	}

	var _err error

	// register each node individually
	for _, node := range s.Nodes {
		err := e.registerNode(ctx, s, node)
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
		if service := decode(n.Value); service != nil {
			main := utils.UnsafeToString(service.Version)
			s, ok := serviceMap[main]
			if !ok {
				s = &micro.Service{
					Name:      service.Name,
					Version:   service.Version,
					Metadata:  service.Metadata,
					Endpoints: service.Endpoints,
				}
				serviceMap[main] = s
			}

			s.Nodes = append(s.Nodes, service.Nodes...)
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
