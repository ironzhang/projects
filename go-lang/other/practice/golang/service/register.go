package service

type Register interface {
	RegistEndpoint(service, url string) error
	UnregistEndpoint(service, url string) error
	UnregistAll()
}
