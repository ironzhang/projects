package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/coreos/etcd/client"
	"github.com/ironzhang/golang/done"
	"github.com/ironzhang/golang/log"
	"golang.org/x/net/context"
)

type etcdDiscover struct {
	kapi  client.KeysAPI
	space *Namespace
	dg    done.DoneGroup
}

func newEtcdDiscover(kapi client.KeysAPI, space *Namespace) *etcdDiscover {
	return &etcdDiscover{
		kapi:  kapi,
		space: space,
	}
}

func (d *etcdDiscover) Start() error {
	ctx, err := d.dg.Add("watch-service")
	if err != nil {
		return err
	}
	go watchService(ctx, d.kapi, d.space)
	return nil
}

func (d *etcdDiscover) Stop() error {
	return d.dg.Done("watch-service", true)
}

func (d *etcdDiscover) GetNamespace() *Namespace {
	return d.space
}

func watchDir(namespace string) string {
	return fmt.Sprintf("/%s/provider/", namespace)
}

func parsePath(path string, namespace string) (service, url string, err error) {
	dir := watchDir(namespace)
	if !strings.HasPrefix(path, dir) {
		log.Errorf("%s has not perfix %s", path, dir)
		return "", "", errors.New("prefix is wrong")
	}
	names := strings.Split(path[len(dir):], "/")
	if len(names) != 2 {
		log.Errorf("%s split wrong", path[len(dir):])
		return "", "", errors.New("path split wrong")
	}
	return names[0], names[1], nil
}

func update(watcher client.Watcher, space *Namespace) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	res, err := watcher.Next(ctx)
	if err != nil {
		if err != context.DeadlineExceeded {
			log.Errorf("watcher next failed, err[%v]", err)
		}
		return
	}
	node := res.Node
	if node == nil {
		log.Errorf("node is nil")
		return
	}
	if node.Dir {
		log.Infof("node is dir")
		return
	}
	service, url, err := parsePath(node.Key, space.name)
	if err != nil {
		log.Errorf("parse path failed, err[%v]", err)
		return
	}

	log.Infof("%s %s %s", res.Action, service, url)
	switch res.Action {
	case "set", "update":
		space.updateEndpoint(service, url)
	case "delete":
		space.deleteEndpoint(service, url)
	case "expire":
		space.expireEndpoint(service, url)
	}
}

func watchService(ctx *done.Context, api client.KeysAPI, space *Namespace) {
	dir := watchDir(space.name)
	opts := &client.WatcherOptions{Recursive: true}
	watcher := api.Watcher(dir, opts)

	log.Infof("watch %s start", dir)
	for {
		update(watcher, space)
		select {
		case <-ctx.Done():
			log.Infof("watch %s done", dir)
			ctx.OK()
			return
		default:
			continue
		}
	}
}
