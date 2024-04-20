package transport

import (
	"github.com/google/uuid"
	"sync"
	"time"
)

type Pool struct {
	tr Transport

	conns map[string][]*Conn
	size  int
	ttl   time.Duration

	sync.Mutex
}

type Conn struct {
	created time.Time
	Client
	id string
}

func (p *Conn) Close() error {
	return nil
}

func (p *Conn) Id() string {
	return p.id
}

func (p *Conn) Created() time.Time {
	return p.created
}

func NewPool(size int, ttl time.Duration, transport Transport) *Pool {
	return &Pool{
		size:  size,
		tr:    transport,
		ttl:   ttl,
		conns: make(map[string][]*Conn),
	}
}

func (p *Pool) Close() error {
	p.Lock()
	defer p.Unlock()

	var err error

	for k, c := range p.conns {
		for _, conn := range c {
			if nerr := conn.Client.Close(); nerr != nil {
				err = nerr
			}
		}

		delete(p.conns, k)
	}

	return err
}

func (p *Pool) Get(addr string, timeout time.Duration) (*Conn, error) {
	p.Lock()
	conns := p.conns[addr]

	// While we have conns check age and then return one
	// otherwise we'll create a new conn
	for len(conns) > 0 {
		conn := conns[len(conns)-1]
		conns = conns[:len(conns)-1]
		p.conns[addr] = conns

		// If conn is old kill it and move on
		if d := time.Since(conn.Created()); d > p.ttl {
			if err := conn.Client.Close(); err != nil {
				p.Unlock()
				return nil, err
			}

			continue
		}

		// We got a good conn, lets unlock and return it
		p.Unlock()

		return conn, nil
	}

	p.Unlock()

	// create new conn
	c, err := p.tr.Dial(addr, timeout)
	if err != nil {
		return nil, err
	}

	return &Conn{
		Client:  c,
		id:      uuid.New().String(),
		created: time.Now(),
	}, nil
}

func (p *Pool) Release(conn *Conn, err error) error {
	// don't store the conn if it has errored
	if err != nil {
		return conn.Client.Close()
	}

	// otherwise put it back for reuse
	p.Lock()
	defer p.Unlock()

	conns := p.conns[conn.Remote()]
	if len(conns) >= p.size {
		return conn.Client.Close()
	} else {
		if err = conn.Client.CloseSend(); err != nil {
			return err
		}
	}
	p.conns[conn.Remote()] = append(conns, conn)

	return nil
}
