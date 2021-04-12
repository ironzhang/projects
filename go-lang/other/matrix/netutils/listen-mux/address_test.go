package listen_mux

import (
	"net"
	"testing"
)

func ListenAddrs(network string, addrs []string) ([]net.Listener, error) {
	lns := make([]net.Listener, 0, len(addrs))
	for _, addr := range addrs {
		ln, err := net.Listen(network, addr)
		if err != nil {
			return nil, err
		}
		lns = append(lns, ln)
	}
	return lns, nil
}

func TestAddress(t *testing.T) {
	tests := []struct {
		addrs []string
		net   string
		str   string
	}{
		{
			addrs: []string{":2000", ":2001", ":2002"},
			net:   "tcp",
			str:   "[::]:2000,[::]:2001,[::]:2002",
		},
		{
			addrs: []string{"127.0.0.1:2003"},
			net:   "tcp",
			str:   "127.0.0.1:2003",
		},
	}

	for i, tt := range tests {
		lns, err := ListenAddrs(tt.net, tt.addrs)
		if err != nil {
			t.Fatalf("tests[%d]: listen addrs: %v", i, err)
		}
		addr := resolvAddress(tt.net, lns)
		if got, want := addr.Network(), tt.net; got != want {
			t.Errorf("network: %s != %s", got, want)
		}
		if got, want := addr.String(), tt.str; got != want {
			t.Errorf("string: %s != %s", got, want)
		}
	}
}
