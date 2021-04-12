package service

type Discover interface {
	Start()
	Stop()
	GetService(service string) *Service
}
