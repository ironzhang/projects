package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/coreos/etcd/client"
	"github.com/ironzhang/golang/done"
	"github.com/ironzhang/golang/log"
	"golang.org/x/net/context"
)

type etcdRegister struct {
	kapi                  client.KeysAPI
	namespace             string
	heartbeatInterval     time.Duration
	heartbeatTimeoutRound int

	dg done.DoneGroup
}

func newEtcdRegister(kapi client.KeysAPI, namespace string, heartbeatInterval time.Duration, heartbeatTimeoutRound int) *etcdRegister {
	return &etcdRegister{
		kapi:                  kapi,
		namespace:             namespace,
		heartbeatInterval:     heartbeatInterval,
		heartbeatTimeoutRound: heartbeatTimeoutRound,
	}
}

func (r *etcdRegister) RegistEndpoint(service, url string) error {
	path := r.path(service, url)
	ctx, err := r.dg.Add(path)
	if err != nil {
		return errors.New("path is registed")
	}
	p := pinger{
		api:      r.kapi,
		ttl:      r.heartbeatInterval * time.Duration(r.heartbeatTimeoutRound),
		interval: r.heartbeatInterval,
		path:     path,
	}
	if err := p.online(); err != nil {
		return err
	}
	go p.heartbeat(ctx)
	return nil
}

func (r *etcdRegister) UnregistEndpoint(service, url string) error {
	path := r.path(service, url)
	if err := r.dg.Done(path, true); err != nil {
		return err
	}
	return nil
}

func (r *etcdRegister) UnregistAll() {
	r.dg.DoneAll(true)
}

func (r *etcdRegister) path(service, url string) string {
	return fmt.Sprintf("/%s/provider/%s/%s", r.namespace, service, url)
}

type pinger struct {
	api      client.KeysAPI
	ttl      time.Duration
	interval time.Duration
	path     string
}

func (p *pinger) online() error {
	opts := &client.SetOptions{TTL: p.ttl}
	_, err := p.api.Set(context.Background(), p.path, "online", opts)
	if err != nil {
		log.Errorf("%s heartbeat ping failed, err[%v]", p.path, err)
		return err
	}
	log.Infof("%s heartbeat ping success", p.path)
	return nil
}

func (p *pinger) offline() error {
	_, err := p.api.Delete(context.Background(), p.path, &client.DeleteOptions{})
	if err != nil {
		log.Errorf("%s heartbeat offline failed, err[%v]", p.path, err)
		return err
	}
	log.Infof("%s heartbeat offline success", p.path)
	return nil
}

func (p *pinger) heartbeat(ctx *done.Context) {
	for {
		select {
		case <-time.After(p.interval):
			p.online()
		case <-ctx.Done():
			p.offline()
			ctx.OK()
			return
		}
	}
}
