package command

import (
	"context"
	"testing"
	"time"

	"github.com/ironzhang/zoutil"
)

const (
	TestAddress     = "localhost:2000"
	TestNodeName    = "command-0"
	TestServiceName = "command"
)

var (
	TestClient *Client
)

func OpenTestZerone(node string) (*zoutil.Zerone, error) {
	opts := zoutil.DefaultOptions
	return zoutil.Open(node, opts)
}

func ListenAndServe(node, service, addr string) {
	z, err := OpenTestZerone(node)
	if err != nil {
		panic(err)
	}
	s, err := z.NewServer(service)
	if err != nil {
		panic(err)
	}
	if err = RegisterCommand(s, z); err != nil {
		panic(err)
	}
	if err := s.ListenAndServe("tcp", addr, node); err != nil {
		panic(err)
	}
}

func DialTestClient(node, service string) {
	z, err := OpenTestZerone(node)
	if err != nil {
		panic(err)
	}
	c, err := Dial(z, service)
	if err != nil {
		panic(err)
	}
	TestClient = c
	time.Sleep(10 * time.Millisecond)
}

func TestMain(m *testing.M) {
	go ListenAndServe(TestNodeName, TestServiceName, TestAddress)
	DialTestClient(TestNodeName, TestServiceName)
	m.Run()
}

func TestGetNodeName(t *testing.T) {
	c := TestClient
	node, err := c.GetNodeName(context.Background())
	if err != nil {
		t.Fatalf("get node name: %v", err)
	}
	if got, want := node, TestNodeName; got != want {
		t.Errorf("node name: got %v, want %v", got, want)
	} else {
		t.Logf("node name: got %v", got)
	}
}

func TestGetSetTimeout(t *testing.T) {
	c := TestClient

	var err error
	var get, set time.Duration

	set = 5 * time.Second
	if err = c.SetTimeout(context.Background(), set); err != nil {
		t.Fatalf("set timeout: %v", err)
	}
	if get, err = c.GetTimeout(context.Background()); err != nil {
		t.Fatalf("get timeout: %v", err)
	}
	if got, want := get, set; got != want {
		t.Errorf("timeout: got %v, want %v", got, want)
	} else {
		t.Logf("timeout: got %v", got)
	}
}

func TestGetSetLogLevel(t *testing.T) {
	c := TestClient

	var err error
	var get, set string

	set = "ERROR"
	if err = c.SetLogLevel(context.Background(), set); err != nil {
		t.Fatalf("set log level: %v", err)
	}
	if get, err = c.GetLogLevel(context.Background()); err != nil {
		t.Fatalf("get log level: %v", err)
	}
	if got, want := get, set; got != want {
		t.Errorf("log level: got %v, want %v", got, want)
	} else {
		t.Logf("log level: got %v", got)
	}
}

func TestGetSetTraceVerbose(t *testing.T) {
	c := TestClient

	tests := []struct {
		object  TraceObject
		verbose int
	}{
		{
			object:  TraceObject{Type: "S"},
			verbose: 0,
		},
		{
			object:  TraceObject{Type: "C"},
			verbose: 1,
		},
		{
			object:  TraceObject{Type: "S", Service: TestServiceName},
			verbose: -1,
		},
	}
	for i, tt := range tests {
		if err := c.SetTraceVerbose(context.Background(), tt.object, tt.verbose); err != nil {
			t.Fatalf("%d: set trace verbose: %v", i, err)
		}
		verbose, err := c.GetTraceVerbose(context.Background(), tt.object)
		if err != nil {
			t.Fatalf("%d: get trace verbose: %v", i, err)
		}
		if got, want := verbose, tt.verbose; got != want {
			t.Errorf("%d: verbose: got %v, want %v", i, got, want)
		} else {
			t.Logf("%d: verbose: got %v", i, got)
		}
	}
}
