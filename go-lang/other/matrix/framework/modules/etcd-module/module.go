package etcd_module

import (
	"fmt"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/ironzhang/matrix/framework"
	"github.com/ironzhang/matrix/jsoncfg"
)

var Config = &C{
	Endpoints:   []string{"localhost:2379", "localhost:22379", "localhost:32379"},
	DialTimeout: jsoncfg.Duration(5 * time.Second),
}

var Module = &M{}

func init() {
	framework.Register(Module, nil, Config)
}

type C struct {
	Endpoints        []string
	AutoSyncInterval jsoncfg.Duration
	DialTimeout      jsoncfg.Duration
	Username         string
	Password         string
}

type M struct {
	client *clientv3.Client
}

func (m *M) Name() string {
	return "etcd-module"
}

func (m *M) Init() (err error) {
	cfg := clientv3.Config{
		Endpoints:        Config.Endpoints,
		AutoSyncInterval: time.Duration(Config.AutoSyncInterval),
		DialTimeout:      time.Duration(Config.DialTimeout),
		Username:         Config.Username,
		Password:         Config.Password,
	}
	m.client, err = clientv3.New(cfg)
	if err != nil {
		return fmt.Errorf("new client: %v", err)
	}
	return nil
}

func (m *M) Fini() error {
	return m.client.Close()
}

func (m *M) Client() *clientv3.Client {
	return m.client
}
