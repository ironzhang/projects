package discovery

import (
	"fmt"
	"sync"
	"time"

	"github.com/coreos/etcd/clientv3"
)

type Refresh func(addrs []string)

type Options struct {
	Namespace string
	Timeout   time.Duration
}

type Discovery struct {
	client    *clientv3.Client
	timeout   time.Duration
	namespace string

	mu sync.RWMutex
	m  map[string]*service
}

func New(client *clientv3.Client, opts Options) *Discovery {
	return &Discovery{
		client:    client,
		timeout:   opts.Timeout,
		namespace: opts.Namespace,
		m:         make(map[string]*service),
	}
}

func (d *Discovery) Namespace() string {
	return d.namespace
}

func (d *Discovery) Watch(svc string, refreshs ...Refresh) (Service, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if s, ok := d.m[svc]; ok {
		s.AddRefreshs(refreshs)
		return s, nil
	}

	s := newService(svc, refreshs...)
	m := monitor{
		client:  d.client,
		prefix:  d.prefix(svc),
		kvs:     make(map[string][]byte),
		refresh: s.Refresh,
	}
	if err := m.Setup(d.timeout); err != nil {
		return nil, err
	}
	s.ok = m.Go(s.done)
	d.m[svc] = s

	return s, nil
}

func (d *Discovery) Unwatch(svc string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if s, ok := d.m[svc]; ok {
		s.Unwatch()
		delete(d.m, svc)
	}
}

func (d *Discovery) UnwatchAll() {
	d.mu.Lock()
	defer d.mu.Unlock()
	for _, s := range d.m {
		s.Unwatch()
	}
	d.m = make(map[string]*service)
}

func (d *Discovery) Service(svc string) (Service, bool) {
	d.mu.RLock()
	s, ok := d.m[svc]
	d.mu.RUnlock()
	return s, ok
}

func (d *Discovery) prefix(svc string) string {
	return fmt.Sprintf("%s/%s/", d.namespace, svc)
}
