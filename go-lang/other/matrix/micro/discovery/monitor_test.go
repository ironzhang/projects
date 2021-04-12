package discovery

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/ironzhang/matrix/tlog"
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

func TestMonitorSetup(t *testing.T) {
	tlog.Init(tlog.Config{DisableStderr: true})

	c := NewClient(t)

	var err error
	if _, err = c.Put(context.Background(), "TestMonitorSetup/k1", "v1"); err != nil {
		t.Fatal(err)
	}
	if _, err = c.Put(context.Background(), "TestMonitorSetup/k2", "v2"); err != nil {
		t.Fatal(err)
	}

	m := monitor{
		client: c,
		prefix: "TestMonitorSetup/",
		kvs:    make(map[string][]byte),
	}
	if err = m.Setup(0); err != nil {
		t.Fatalf("setup: %v", err)
	}

	want := map[string][]byte{
		"k1": []byte("v1"),
		"k2": []byte("v2"),
	}
	if got := m.kvs; !reflect.DeepEqual(got, want) {
		t.Fatalf("kvs: %v != %v", got, want)
	}
}

func TestMonitorRun(t *testing.T) {
	tlog.Init(tlog.Config{DisableStderr: true})

	c := NewClient(t)

	m := monitor{
		client: c,
		prefix: "TestMonitorRun/",
		kvs:    make(map[string][]byte),
	}
	ctx, cancel := context.WithCancel(context.Background())
	go m.Run(ctx)
	time.Sleep(100 * time.Millisecond)

	var err error
	var want = map[string][]byte{
		"k1": []byte("v1"),
		"k2": []byte("v2"),
	}

	if _, err = c.Put(context.Background(), "TestMonitorRun/k1", "v1"); err != nil {
		t.Fatal(err)
	}
	if _, err = c.Put(context.Background(), "TestMonitorRun/k2", "v2"); err != nil {
		t.Fatal(err)
	}
	if got := m.kvs; !reflect.DeepEqual(got, want) {
		t.Fatalf("kvs: %v != %v", got, want)
	}

	if _, err = c.Delete(context.Background(), "TestMonitorRun/k2"); err != nil {
		t.Fatal(err)
	}
	delete(want, "k2")
	if got := m.kvs; !reflect.DeepEqual(got, want) {
		t.Fatalf("kvs: %v != %v", got, want)
	}

	cancel()
	time.Sleep(100 * time.Millisecond)
}

func TestMonitorGo(t *testing.T) {
	tlog.Init(tlog.Config{DisableStderr: true})

	c := NewClient(t)

	var err error
	if _, err = c.Put(context.Background(), "TestMonitorGo/k1", "v1"); err != nil {
		t.Fatal(err)
	}
	if _, err = c.Put(context.Background(), "TestMonitorGo/k2", "v2"); err != nil {
		t.Fatal(err)
	}

	m := monitor{
		client: c,
		prefix: "TestMonitorGo/",
		kvs:    make(map[string][]byte),
	}

	if err = m.Setup(0); err != nil {
		t.Fatalf("setup: %v", err)
	}

	done := make(chan struct{})
	ok := m.Go(done)

	if _, err = c.Delete(context.Background(), "TestMonitorGo/k2"); err != nil {
		t.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)

	want := map[string][]byte{"k1": []byte("v1")}
	if got := m.kvs; !reflect.DeepEqual(got, want) {
		t.Fatalf("kvs: %v != %v", got, want)
	}

	close(done)
	<-ok
}
