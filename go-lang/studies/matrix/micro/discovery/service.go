package discovery

import (
	"sort"
	"sync"
)

type Service interface {
	Name() string
	Addrs() []string
}

func newService(name string, refreshs ...Refresh) *service {
	return &service{
		name:     name,
		refreshs: refreshs,
		done:     make(chan struct{}),
	}
}

type service struct {
	name     string
	refreshm sync.RWMutex
	refreshs []Refresh
	done     chan struct{}
	ok       <-chan struct{}
	addrm    sync.RWMutex
	addrs    []string
}

func (s *service) Name() string {
	return s.name
}

func (s *service) Addrs() []string {
	s.addrm.RLock()
	addrs := s.addrs
	s.addrm.RUnlock()
	return addrs
}

func (s *service) SetAddrs(addrs []string) {
	s.addrm.Lock()
	s.addrs = addrs
	s.addrm.Unlock()
}

func (s *service) Refresh(kvs map[string][]byte) {
	addrs := make([]string, 0, len(kvs))
	for k, _ := range kvs {
		addrs = append(addrs, k)
	}
	sort.Strings(addrs)
	s.SetAddrs(addrs)

	s.refreshm.RLock()
	defer s.refreshm.RUnlock()
	for _, refresh := range s.refreshs {
		refresh(addrs)
	}
}

func (s *service) AddRefreshs(refreshs []Refresh) {
	if len(refreshs) > 0 {
		s.refreshm.Lock()
		s.refreshs = append(s.refreshs, refreshs...)
		s.refreshm.Unlock()

		addrs := s.Addrs()
		for _, refresh := range refreshs {
			refresh(addrs)
		}
	}
}

func (s *service) Unwatch() {
	close(s.done)
	<-s.ok
}
