package experimental

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/coreos/etcd/clientv3"
)

func get(c *clientv3.Client, key string) (string, error) {
	resp, err := c.Get(context.Background(), key)
	if err != nil {
		return "", err
	}
	if len(resp.Kvs) > 0 {
		return string(resp.Kvs[0].Value), nil
	}
	return "", nil
}

func KVAssert(t *testing.T, n string, c *clientv3.Client, key, value string) {
	got, err := get(c, key)
	if err != nil {
		t.Fatalf("%s: %v", n, err)
	}
	if want := value; got != want {
		t.Errorf("%s: key: %q, value: %q != %q", n, key, got, want)
	}
}

func NewTestClient(t *testing.T) *clientv3.Client {
	c, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379", "localhost:22379", "localhost:32379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func TestRegister(t *testing.T) {
	c := NewTestClient(t)

	key := "sample_key"
	value := "sample_value"
	ttl := int64(10)

	id, err := register(context.Background(), c, key, value, ttl)
	if err != nil {
		t.Fatal(err)
	}
	c.KeepAlive(context.Background(), id)

	for i := 0; i < 50; i++ {
		KVAssert(t, fmt.Sprintf("%d", i), c, key, value)
		time.Sleep(time.Second)
	}

	unregister(context.Background(), c, id, key)
	KVAssert(t, "unregister", c, key, "")
}

func TestClientGet(t *testing.T) {
	c := NewTestClient(t)

	var err error
	if _, err = c.Put(context.Background(), "Test/S1/A1", ""); err != nil {
		t.Fatal(err)
	}
	if _, err = c.Put(context.Background(), "Test/S1/A2", ""); err != nil {
		t.Fatal(err)
	}

	resp, err := c.Get(context.Background(), "Test/S1", clientv3.WithPrefix())
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Kvs) != 2 {
		t.Fatalf("length of kvs: %d != 2", len(resp.Kvs))
	}
	if got, want := string(resp.Kvs[0].Key), "Test/S1/A1"; got != want {
		t.Errorf("kvs[0]: %q != %q", got, want)
	}
	if got, want := string(resp.Kvs[1].Key), "Test/S1/A2"; got != want {
		t.Errorf("kvs[0]: %q != %q", got, want)
	}
}

func TestClientWatch(t *testing.T) {
	c := NewTestClient(t)

	ch := c.Watch(context.Background(), "Test/S1", clientv3.WithPrefix())
	go func() {
		for resp := range ch {
			for _, e := range resp.Events {
				fmt.Printf("%s\n", e)
			}
		}
	}()

	var err error
	if _, err = c.Put(context.Background(), "Test/S1/A1", ""); err != nil {
		t.Fatal(err)
	}
	if _, err = c.Put(context.Background(), "Test/S1/A2", ""); err != nil {
		t.Fatal(err)
	}
}
