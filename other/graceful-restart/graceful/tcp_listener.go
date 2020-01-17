package graceful

import (
	"fmt"
	"net"
	"os"
	"syscall"
	"time"
)

func shouldGracefulRestart() bool {
	s := os.Getenv("graceful_restart")
	return s == "on"
}

func setGracefulRestartEnv() {
	os.Setenv("graceful_restart", "on")
}

func newTCPListenerFromFD(fd uintptr) (*net.TCPListener, error) {
	ln, err := net.FileListener(os.NewFile(fd, ""))
	if err != nil {
		return nil, err
	}
	sock, ok := ln.(*net.TCPListener)
	if !ok {
		return nil, fmt.Errorf("file descriptor %d is not a valid tcp socket", fd)
	}
	return sock, nil
}

func ListenTCP(network, address string, fd uintptr) (*net.TCPListener, error) {
	if shouldGracefulRestart() {
		return newTCPListenerFromFD(fd)
	}

	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, err
	}
	return net.ListenTCP(network, addr)
}

func StopTCPListeners(listeners ...*net.TCPListener) (err error) {
	for _, ln := range listeners {
		if err = ln.SetDeadline(time.Now()); err != nil {
			return err
		}
	}
	return nil
}

func ForkExec(listeners ...*net.TCPListener) (pid int, err error) {
	if err = StopTCPListeners(listeners...); err != nil {
		return 0, err
	}

	var f *os.File
	files := []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()}
	for _, ln := range listeners {
		if f, err = ln.File(); err != nil {
			return 0, err
		}
		files = append(files, f.Fd())
	}

	setGracefulRestartEnv()
	attr := syscall.ProcAttr{Env: os.Environ(), Files: files}
	pid, err = syscall.ForkExec(os.Args[0], os.Args, &attr)
	if err != nil {
		return 0, err
	}
	return pid, nil
}
