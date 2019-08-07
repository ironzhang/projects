package listen_mux

import (
	"crypto/tls"
	"errors"
	"net"
	"sync"

	"github.com/ironzhang/matrix/tlog"
)

var Verbose bool

var (
	errClosing = errors.New("use of closed network connection")
	errUnknown = errors.New("unknown")
)

func Listen(network string, addrs []string, backlog int) (net.Listener, error) {
	listeners := make([]net.Listener, 0, len(addrs))
	for _, addr := range addrs {
		ln, err := net.Listen(network, addr)
		if err != nil {
			return nil, err
		}
		listeners = append(listeners, ln)
	}
	return NewListener(network, listeners, backlog), nil
}

func ListenTLS(network string, addrs []string, backlog int, config *tls.Config) (net.Listener, error) {
	listeners := make([]net.Listener, 0, len(addrs))
	for _, addr := range addrs {
		ln, err := tls.Listen(network, addr, config)
		if err != nil {
			return nil, err
		}
		listeners = append(listeners, ln)
	}
	return NewListener(network, listeners, backlog), nil
}

func NewListener(network string, listeners []net.Listener, backlog int) net.Listener {
	if len(listeners) == 1 {
		return listeners[0]
	}
	if backlog <= 0 {
		backlog = len(listeners)
	}
	return &listener{
		listeners: listeners,
		backlog:   backlog,
		addr:      resolvAddress(network, listeners),
	}
}

type listener struct {
	listeners []net.Listener
	backlog   int
	addr      address

	once   sync.Once
	ch     <-chan net.Conn
	closed bool
}

func (l *listener) Accept() (net.Conn, error) {
	l.once.Do(func() {
		l.ch = accept(l.listeners, l.backlog)
	})

	c, ok := <-l.ch
	if !ok {
		if l.closed {
			return nil, &net.OpError{Op: "accept", Net: l.addr.Network(), Addr: l.addr, Err: errClosing}
		} else {
			return nil, &net.OpError{Op: "accept", Net: l.addr.Network(), Addr: l.addr, Err: errUnknown}
		}
	}
	return c, nil
}

func (l *listener) Close() error {
	l.closed = true
	for _, ln := range l.listeners {
		ln.Close()
	}
	return nil
}

func (l *listener) Addr() net.Addr {
	return l.addr
}

func accept(ls []net.Listener, backlog int) <-chan net.Conn {
	ch := make(chan net.Conn, backlog)
	go func() {
		doAccept(ls, ch)
		close(ch)
	}()
	return ch
}

func doAccept(ls []net.Listener, ch chan<- net.Conn) {
	var wg sync.WaitGroup
	for _, ln := range ls {
		wg.Add(1)
		go func(ln net.Listener) {
			defer wg.Done()
			loopAccept(ln, ch)
		}(ln)
	}
	wg.Wait()
}

func loopAccept(ln net.Listener, ch chan<- net.Conn) {
	log := tlog.Std().Sugar().With("address", ln.Addr().String())
	for {
		c, err := ln.Accept()
		if err != nil {
			if Verbose {
				log.Debugw("accept", "error", err)
			}
			break
		}
		ch <- c
	}
}
