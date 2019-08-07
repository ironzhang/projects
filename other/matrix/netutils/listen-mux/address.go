package listen_mux

import "net"

type address struct {
	net string
	str string
}

func resolvAddress(network string, ls []net.Listener) address {
	var s string
	for i, ln := range ls {
		if i == 0 {
			s = ln.Addr().String()
		} else {
			s = s + "," + ln.Addr().String()
		}
	}
	return address{net: network, str: s}
}

func (a address) Network() string {
	return a.net
}

func (a address) String() string {
	return a.str
}
