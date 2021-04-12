package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/ironzhang/tlog"
)

func main() {
	ln, err := NewTCPListener(":8081")
	if err != nil {
		tlog.Fatalw("new listener", "error", err)
		return
	}
	tlog.Infof("listen on %s", ln.Addr())

	signalc := make(chan os.Signal)
	signal.Notify(signalc, syscall.SIGHUP, syscall.SIGTERM)
	for sig := range signalc {
		switch sig {
		// 优雅关闭
		case syscall.SIGTERM:

		// 优雅重启
		case syscall.SIGHUP:
			if err := ForkTCPListener(ln); err != nil {
				tlog.Errorw("fork tcp listener", "error", err)
			}
			os.Exit(0)
		}
	}
}
