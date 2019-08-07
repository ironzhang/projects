package service

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
)

func newTestEtcdDiscover() (*etcdDiscover, error) {
	namespace := "test-env"
	servers := []string{"http://127.0.0.1:2379"}

	cfg := client.Config{
		Endpoints: servers,
		Transport: client.DefaultTransport,
	}
	c, err := client.New(cfg)
	if err != nil {
		return nil, err
	}
	if err = c.Sync(context.Background()); err != nil {
		return nil, err
	}
	return newEtcdDiscover(client.NewKeysAPI(c), newNamespace(namespace)), nil
}

func regist(t *testing.T) {
	r, err := newTestEtcdRegister()
	if err != nil {
		t.Fatalf("new test etcd register failed, err[%v]", err)
		return
	}
	if err = r.RegistEndpoint("TestEtcd", "192.168.0.12:7200"); err != nil {
		t.Fatalf("regist endpoint failed, err[%v]", err)
		return
	}
	if err = r.RegistEndpoint("TestEtcd", "192.168.0.12:7201"); err != nil {
		t.Fatalf("regist endpoint failed, err[%v]", err)
		return
	}
	if err = r.UnregistEndpoint("TestEtcd", "192.168.0.12:7200"); err != nil {
		t.Fatalf("unregist endpoint failed, err[%v]", err)
		return
	}
	//r.UnregistAll()
	//time.Sleep(3 * time.Second)
}

func TestEtcdDiscover(t *testing.T) {
	d, err := newTestEtcdDiscover()
	if err != nil {
		t.Fatalf("new test etcd discover, err[%v]", err)
		return
	}

	d.Start()
	defer d.Stop()
	regist(t)
	s, ok := d.GetNamespace().GetService("TestEtcd")
	if !ok {
		t.Fatalf("get service failed")
		return
	}
	for i := 0; i < 5; i++ {
		points := s.GetEndpoints()
		fmt.Printf("%v\n", points)
		runtime.Gosched()
	}
}
