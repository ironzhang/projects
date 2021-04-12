package service

import "sync"

const (
	ENDPOINT_STATUS_ONLINE  = "online"
	ENDPOINT_STATUS_OFFLINE = "offline"
)

type Endpoint struct {
	URL     string
	Status  string
	Service string
}

type Service struct {
	name string
	mu   sync.RWMutex
	l    []Endpoint
}

func newService(name string) *Service {
	return &Service{
		name: name,
		l:    make([]Endpoint, 0, 5),
	}
}

func (s *Service) setEndpoint(url, status string) {
	p := Endpoint{
		URL:     url,
		Status:  status,
		Service: s.name,
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, v := range s.l {
		if v.URL == url {
			s.l[i] = p
			return
		}
	}
	s.l = append(s.l, p)
}

func (s *Service) delEndpoint(url string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, v := range s.l {
		if v.URL == url {
			//s.l = append(s.l[:i], s.l[i+1:]...)
			s.l[i] = s.l[len(s.l)-1]
			s.l = s.l[:len(s.l)-1]
			return
		}
	}
}

func (s *Service) GetEndpoints() []Endpoint {
	s.mu.RLock()
	defer s.mu.RUnlock()
	l := make([]Endpoint, 0, len(s.l))
	for _, p := range s.l {
		l = append(l, p)
	}
	return l
}

type Namespace struct {
	name string
	mu   sync.Mutex
	m    map[string]*Service
}

func newNamespace(name string) *Namespace {
	return &Namespace{
		name: name,
		m:    make(map[string]*Service),
	}
}

func (ns *Namespace) service(name string) *Service {
	ns.mu.Lock()
	defer ns.mu.Unlock()
	s, ok := ns.m[name]
	if ok {
		return s
	}
	s = newService(name)
	ns.m[name] = s
	return s
}

func (ns *Namespace) updateEndpoint(service, url string) {
	ns.service(service).setEndpoint(url, ENDPOINT_STATUS_ONLINE)
}

func (ns *Namespace) deleteEndpoint(service, url string) {
	ns.service(service).delEndpoint(url)
}

func (ns *Namespace) expireEndpoint(service, url string) {
	ns.service(service).setEndpoint(url, ENDPOINT_STATUS_OFFLINE)
}

func (ns *Namespace) GetService(service string) (*Service, bool) {
	ns.mu.Lock()
	defer ns.mu.Unlock()
	s, ok := ns.m[service]
	return s, ok
}
