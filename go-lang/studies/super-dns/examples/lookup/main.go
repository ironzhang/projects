package main

import (
	"net"

	"github.com/ironzhang/tlog"
)

func main() {
	for i := 0; i < 10; i++ {
		ips, err := net.LookupIP("www.baidu.com")
		if err != nil {
			tlog.Errorw("net lookup IP", "error", err)
			return
		}

		for _, ip := range ips {
			tlog.Infof("%s", ip)
		}
	}
}
