package discovery_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/ironzhang/matrix/micro/discovery"
	"github.com/ironzhang/matrix/tlog"
)

func TestDiscovery(t *testing.T) {
	tlog.Init(tlog.Config{DisableStderr: true})

	c := discovery.NewClient(t)
	c.Delete(context.Background(), "TestDiscovery/", clientv3.WithPrefix())

	d := discovery.New(c, discovery.Options{Namespace: "TestDiscovery"})
	s, err := d.Watch("Service")
	if err != nil {
		t.Fatalf("watch: %v", err)
	}
	if got, want := s.Name(), "Service"; got != want {
		t.Errorf("service name: %q != %q", got, want)
	}

	if _, err = c.Put(context.Background(), "TestDiscovery/Service/127.0.0.1:2001", "1"); err != nil {
		t.Fatal(err)
	}
	if _, err = c.Put(context.Background(), "TestDiscovery/Service/127.0.0.1:2002", "1"); err != nil {
		t.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)
	if got, want := s.Addrs(), []string{"127.0.0.1:2001", "127.0.0.1:2002"}; !reflect.DeepEqual(got, want) {
		t.Errorf("addrs: %v != %v", got, want)
	}

	if _, err = c.Put(context.Background(), "TestDiscovery/Service/127.0.0.1:2000", "1"); err != nil {
		t.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)
	if got, want := s.Addrs(), []string{"127.0.0.1:2000", "127.0.0.1:2001", "127.0.0.1:2002"}; !reflect.DeepEqual(got, want) {
		t.Errorf("addrs: %v != %v", got, want)
	}

	if _, err = c.Delete(context.Background(), "TestDiscovery/Service/127.0.0.1:2001"); err != nil {
		t.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)
	if got, want := s.Addrs(), []string{"127.0.0.1:2000", "127.0.0.1:2002"}; !reflect.DeepEqual(got, want) {
		t.Errorf("addrs: %v != %v", got, want)
	}

	d.UnwatchAll()
}

func TestRefresh(t *testing.T) {
	tlog.Init(tlog.Config{DisableStderr: true})

	var count int
	f := func(addrs []string) {
		count++
	}

	c := discovery.NewClient(t)
	d := discovery.New(c, discovery.Options{Namespace: "TestRefresh"})
	if _, err := d.Watch("Service", f); err != nil {
		t.Fatalf("watch: %v", err)
	}
	if _, err := d.Watch("Service", f); err != nil {
		t.Fatalf("watch: %v", err)
	}
	if got, want := count, 2; got != want {
		t.Errorf("after watch, count: %v != %v", got, want)
	}

	if _, err := c.Put(context.Background(), "TestRefresh/Service/127.0.0.1:2001", "1"); err != nil {
		t.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)
	if got, want := count, 4; got != want {
		t.Errorf("after put, count: %v != %v", got, want)
	}

	d.UnwatchAll()
}
