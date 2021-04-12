package model

type Endpoint struct {
}

type Address struct {
	Addr     string
	Metadata interface{}
}

type Tags map[string][]string
