package connsvr

import (
	"log"
	"net"
	"sync"
	"sync/atomic"
)

type Parser interface {
	Parse(c *Conn) (msg []byte, err error)
}

type Handler interface {
	OnConnect(c *Conn) error
	OnDisconnect(c *Conn, kickout bool)
	OnRecvMessage(c *Conn, msg []byte)
}

type Conn struct {
	net.Conn
	id      int
	kickout int32
	value   atomic.Value
}

func newConn(id int, nc net.Conn) *Conn {
	return &Conn{
		Conn: nc,
		id:   id,
	}
}

func (c *Conn) Id() int {
	return c.id
}

func (c *Conn) isKickout() bool {
	return atomic.LoadInt32(&c.kickout) != 0
}

func (c *Conn) Kickout() error {
	atomic.StoreInt32(&c.kickout, 1)
	return c.Conn.Close()
}

func (c *Conn) Value() interface{} {
	return c.value.Load()
}

func (c *Conn) SetValue(v interface{}) {
	c.value.Store(v)
}

type Server struct {
	net     string
	addr    string
	parser  Parser
	handler Handler

	mu    sync.RWMutex
	conns map[int]*Conn
}

func NewServer(net, addr string, parser Parser, handler Handler) *Server {
	return new(Server).Init(net, addr, parser, handler)
}

func (s *Server) Init(net, addr string, parser Parser, handler Handler) *Server {
	s.net = net
	s.addr = addr
	s.parser = parser
	s.handler = handler
	s.conns = make(map[int]*Conn)
	return s
}

func (s *Server) ListenAndServe() error {
	ln, err := net.Listen(s.net, s.addr)
	if err != nil {
		log.Printf("listen: %v", err)
		return err
	}

	connid := 0
	for {
		nc, err := ln.Accept()
		if err != nil {
			log.Printf("accept: %v", err)
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				continue
			}
			break
		}

		connid++
		c := newConn(connid, nc)
		go s.serveConn(c)
	}

	return nil
}

func (s *Server) serveConn(c *Conn) {
	defer c.Conn.Close()

	s.addConn(c)
	s.handler.OnConnect(c)
	for {
		msg, err := s.parser.Parse(c)
		if err != nil {
			log.Printf("recv message: %v", err)
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				continue
			}
			break
		}
		log.Printf("recv message: %d", msg)
		s.handler.OnRecvMessage(c, msg)
	}
	s.handler.OnDisconnect(c, c.isKickout())
	s.delConn(c.id)
}

func (s *Server) addConn(c *Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.conns[c.id] = c
}

func (s *Server) delConn(id int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.conns, id)
}

func (s *Server) GetConn(id int) (*Conn, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.conns[id]
	return c, ok
}

func (s *Server) ListConns() []*Conn {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var conns []*Conn
	for _, c := range s.conns {
		conns = append(conns, c)
	}
	return conns
}
