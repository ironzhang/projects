package main

import (
	"net"
	"net/http"
	"time"

	"git.xiaojukeji.com/disfv4/commlib/httputils"
	"git.xiaojukeji.com/pearls/tlog"
)

func RunLookupHost() {
	for i := 0; i < 10; i++ {
		ips, err := net.LookupHost("www.baidu.com")
		if err != nil {
			tlog.Errorw("net lookup IP", "error", err)
			return
		}

		for _, ip := range ips {
			tlog.Infof("%s", ip)
		}
	}
}

func RunLookupPort() {
	port, err := net.LookupPort("tcp", "thrift")
	if err != nil {
		tlog.Errorw("lookup port", "error", err)
		return
	}
	tlog.Infof("port=%d", port)
}

func RunLookupSRV() {
	_, addrs, err := net.LookupSRV("thrift", "tcp", "www.superdns.com")
	if err != nil {
		tlog.Errorw("lookup srv", "error", err)
		return
	}
	tlog.Infof("addrs=%v", addrs)
}

func RunDial() {
	conn, err := net.DialTimeout("tcp", "http.baidu.com:thrift", 3*time.Second)
	if err != nil {
		tlog.Errorw("dial timeout", "error", err)
		return
	}
	conn.Close()
}

func RunHTTP() {
	http.DefaultClient.Transport = &httputils.AccessLogRoundTripper{}
	http.DefaultClient.Timeout = 2 * time.Second
	_, err := http.Get("http://http.baidu.com")
	if err != nil {
		tlog.Errorw("http get", "error", err)
		return
	}
}

func main() {
	RunLookupSRV()
}
