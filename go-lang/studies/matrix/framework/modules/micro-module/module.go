package micro_module

import (
	"time"

	"github.com/ironzhang/matrix/framework"
	"github.com/ironzhang/matrix/framework/modules/etcd-module"
	"github.com/ironzhang/matrix/jsoncfg"
	"github.com/ironzhang/matrix/micro/discovery"
	"github.com/ironzhang/matrix/micro/registry"
)

var Config = &C{
	Namespace: "matrix",
	Timeout:   jsoncfg.Duration(5 * time.Second),
	TTL:       10,
}

var Module = &M{}

func init() {
	framework.Register(Module, nil, Config)
}

type C struct {
	Namespace string
	Timeout   jsoncfg.Duration
	TTL       int64
}

type M struct {
	r *registry.Registry
	d *discovery.Discovery
}

func (m *M) Name() string {
	return "micro-module"
}

func (m *M) Init() error {
	c := etcd_module.Module.Client()
	m.r = registry.New(c, registry.Options{
		TTL:       Config.TTL,
		Timeout:   time.Duration(Config.Timeout),
		Namespace: Config.Namespace,
	})
	m.d = discovery.New(c, discovery.Options{
		Namespace: Config.Namespace,
		Timeout:   time.Duration(Config.Timeout),
	})
	return nil
}

func (m *M) Fini() error {
	m.r.UnregisterAll()
	m.d.UnwatchAll()
	return nil
}

func (m *M) Registry() *registry.Registry {
	return m.r
}

func (m *M) Discovery() *discovery.Discovery {
	return m.d
}
