package registry

import (
	"context"
	"testing"
	"time"

	"github.com/coreos/etcd/clientv3"
)

func NewClient(t *testing.T) *clientv3.Client {
	c, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379", "localhost:22379", "localhost:32379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func TestPinger(t *testing.T) {
	c := NewClient(t)
	p := newPinger(c, 2*time.Second, 5, "TestPinger/Key", "1")

	var err error
	var exist bool

	if err = p.Setup(); err != nil {
		t.Fatalf("setup: %v", err)
	}

	for i := 0; i <= 10; i++ {
		if exist, err = p.Exist(context.Background()); err != nil {
			t.Fatalf("exist: %v", err)
		} else if !exist {
			t.Errorf("after %d second: key not exist", i)
		} else {
			t.Logf("after %d second: key exist", i)
		}
		time.Sleep(time.Second)
	}

	if err = p.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}
	if exist, err = p.Exist(context.Background()); err != nil {
		t.Fatalf("exist: %v", err)
	} else if exist {
		t.Errorf("after close: key exist")
	} else {
		t.Logf("after close: key not exist")
	}
}
