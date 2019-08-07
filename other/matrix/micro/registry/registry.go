package registry

import (
	"fmt"
	"sync"
	"time"

	"github.com/coreos/etcd/clientv3"
)

type Endpoint struct {
	Service string
	Addr    string
}

type Options struct {
	TTL       int64
	Timeout   time.Duration
	Namespace string
}

type Registry struct {
	client    *clientv3.Client
	ttl       int64
	timeout   time.Duration
	namespace string

	mu      sync.Mutex
	pingers map[string]*pinger
}

func New(c *clientv3.Client, opts Options) *Registry {
	if opts.TTL <= 0 {
		opts.TTL = 10 // default ttl: 10s
	}
	return &Registry{
		client:    c,
		ttl:       opts.TTL,
		timeout:   opts.Timeout,
		namespace: opts.Namespace,
		pingers:   make(map[string]*pinger),
	}
}

func (r *Registry) Namespace() string {
	return r.namespace
}

func (r *Registry) Register(point Endpoint) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := r.key(point)
	if _, ok := r.pingers[key]; ok {
		return fmt.Errorf("key(%s) existed", key)
	}

	p := newPinger(r.client, r.timeout, r.ttl, key, "1")
	if err := p.Setup(); err != nil {
		return err
	}
	r.pingers[key] = p

	return nil
}

func (r *Registry) Unregister(point Endpoint) (err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := r.key(point)
	if p, ok := r.pingers[key]; ok {
		err = p.Close()
		delete(r.pingers, key)
	}
	return err
}

func (r *Registry) UnregisterAll() {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, p := range r.pingers {
		p.Close()
	}
	r.pingers = make(map[string]*pinger)
}

func (r *Registry) key(p Endpoint) string {
	return fmt.Sprintf("%s/%s/%s", r.namespace, p.Service, p.Addr)
}
