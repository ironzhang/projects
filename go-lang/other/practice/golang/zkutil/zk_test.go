package zkutil

import (
	"testing"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

var servers = []string{"172.17.0.7"}

func TestZk(t *testing.T) {
	conn, _, err := zk.Connect(servers, time.Second)
	if err != nil {
		t.Fatalf("failed to connect to zk server, %v", err)
	}

	CreateRecursive(conn, "/zk/test", "", 0, DefaultDirACLs())
	//WatchChildren(conn, "/zk/test", func(children []string) { t.Logf("%v", children) })
}
