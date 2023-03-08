package service

import (
	"testing"
	"time"

	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
)

func newTestEtcdRegister() (*etcdRegister, error) {
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
	return newEtcdRegister(client.NewKeysAPI(c), namespace, 2*time.Second, 3), nil
}

func TestEtcdRegister(t *testing.T) {
	//runtime.GOMAXPROCS(1)
	r, err := newTestEtcdRegister()
	if err != nil {
		t.Fatalf("new test etcd register, err[%v]", err)
		return
	}
	if err = r.RegistEndpoint("TestEtcdRegister", "192.168.0.12:7200"); err != nil {
		t.Fatalf("regist endpoint, err[%d]", err)
		return
	}
	time.Sleep(3 * time.Second)
	//r.UnregistEndpoint("TestEtcdRegister", "192.168.0.12:7200")
	r.UnregistAll()
}
