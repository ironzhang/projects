package registry_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/ironzhang/matrix/micro/registry"
	"github.com/ironzhang/matrix/tlog"
)

func MakeKey(namespace string, p registry.Endpoint) string {
	return fmt.Sprintf("%s/%s/%s", namespace, p.Service, p.Addr)
}

func KeyIsExist(c *clientv3.Client, key string) bool {
	log := tlog.Std().Sugar().With("key", key)
	resp, err := c.Get(context.Background(), key)
	if err != nil {
		log.Errorw("client get", "error", err)
		return false
	}
	if len(resp.Kvs) != 1 {
		return false
	}
	if string(resp.Kvs[0].Key) != key {
		return false
	}
	return true
}

func TestRegistry(t *testing.T) {
	c := registry.NewClient(t)
	r := registry.New(c, registry.Options{Namespace: "TestRegistry"})

	ns := r.Namespace()
	p1 := registry.Endpoint{"S1", "A1"}
	p2 := registry.Endpoint{"S2", "A2"}

	if key := MakeKey(ns, p1); KeyIsExist(c, key) {
		t.Fatalf("before register: key(%s) is existed", key)
	}
	if key := MakeKey(ns, p2); KeyIsExist(c, key) {
		t.Fatalf("before register: key(%s) is existed", key)
	}

	var err error
	if err = r.Register(p1); err != nil {
		t.Fatalf("register p1: %v", err)
	}
	if err = r.Register(p2); err != nil {
		t.Fatalf("register p2: %v", err)
	}
	if key := MakeKey(ns, p1); !KeyIsExist(c, key) {
		t.Fatalf("after register: key(%s) is not existed", key)
	}
	if key := MakeKey(ns, p2); !KeyIsExist(c, key) {
		t.Fatalf("after register: key(%s) is not existed", key)
	}

	time.Sleep(11 * time.Second)
	if key := MakeKey(ns, p1); !KeyIsExist(c, key) {
		t.Fatalf("after sleep: key(%s) is not existed", key)
	}
	if key := MakeKey(ns, p2); !KeyIsExist(c, key) {
		t.Fatalf("after sleep: key(%s) is not existed", key)
	}

	if err = r.Unregister(p1); err != nil {
		t.Fatalf("unregister p1: %v", err)
	}
	if key := MakeKey(ns, p1); KeyIsExist(c, key) {
		t.Fatalf("after unregister p1: key(%s) is existed", key)
	}
	if key := MakeKey(ns, p2); !KeyIsExist(c, key) {
		t.Fatalf("after unregister p1: key(%s) is not existed", key)
	}

	r.UnregisterAll()
	if key := MakeKey(ns, p1); KeyIsExist(c, key) {
		t.Fatalf("after unregister all: key(%s) is existed", key)
	}
	if key := MakeKey(ns, p2); KeyIsExist(c, key) {
		t.Fatalf("after unregister all: key(%s) is existed", key)
	}
}
