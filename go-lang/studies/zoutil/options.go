package zoutil

import (
	"fmt"
	"time"

	"github.com/coreos/etcd/client"
	"github.com/ironzhang/x-pearls/config"
	"github.com/ironzhang/x-pearls/govern/etcdv2"
	"github.com/ironzhang/zerone"
)

type EtcdOptions struct {
	Endpoints []string
	Username  string
	Password  string
}

type Options struct {
	Zerone        string
	Filename      string
	Namespace     string
	Driver        string
	Etcd          EtcdOptions
	Timeout       config.Duration
	ClientVerbose int
	ServerVerbose int
}

var DefaultOptions = Options{
	Zerone:    "DZerone",
	Namespace: "zerone",
	Driver:    etcdv2.DriverName,
	Etcd:      EtcdOptions{Endpoints: []string{"http://localhost:2379"}},
	Timeout:   config.Duration(10 * time.Second),
}

func (o Options) zerone() (zerone.Options, error) {
	switch o.Zerone {
	case "SZerone":
		return o.szerone()
	case "DZerone":
		return o.dzerone()
	default:
		return nil, fmt.Errorf("zerone %s is unknown", o.Zerone)
	}
}

func (o Options) szerone() (zerone.SOptions, error) {
	return zerone.SOptions{Filename: o.Filename}, nil
}

func (o Options) dzerone() (zerone.DOptions, error) {
	switch o.Driver {
	case etcdv2.DriverName:
		return zerone.DOptions{
			Namespace: o.Namespace,
			Driver:    o.Driver,
			Config:    client.Config{Endpoints: o.Etcd.Endpoints, Username: o.Etcd.Username, Password: o.Etcd.Password},
		}, nil
	default:
		return zerone.DOptions{}, fmt.Errorf("driver %s is unknown", o.Driver)
	}
}
