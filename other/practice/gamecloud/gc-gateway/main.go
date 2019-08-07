package main

import (
	"bufio"
	"fmt"
	"log"
	"sync"

	"github.com/ironzhang/gamecloud/gc-gateway/connsvr"
	"github.com/ironzhang/golang/lru"
)

type MsgParser struct {
	mu    sync.Mutex
	cache lru.Cache
}

func NewMsgParser(maxEntries int) *MsgParser {
	p := MsgParser{}
	p.cache.Init(maxEntries, nil)
	return &p
}

func (p *MsgParser) newReader(c *connsvr.Conn) *bufio.Reader {
	p.mu.Lock()
	defer p.mu.Unlock()
	if a, ok := p.cache.Get(c); ok {
		return a.(*bufio.Reader)
	}
	r := bufio.NewReader(c)
	p.cache.Add(c, r)
	return r
}

func (p *MsgParser) Parse(c *connsvr.Conn) ([]byte, error) {
	r := p.newReader(c)
	line, _, err := r.ReadLine()
	if err != nil {
		return nil, err
	}
	return line, nil
}

type GatewayServer struct {
	connsvr.Server
}

func NewGatewayServer(net, addr string) *GatewayServer {
	return new(GatewayServer).Init(net, addr)
}

func (s *GatewayServer) Init(net, addr string) *GatewayServer {
	s.Server.Init(net, addr, NewMsgParser(2000), s)
	return s
}

func (s *GatewayServer) DownProcess() {
}

func (s *GatewayServer) OnConnect(c *connsvr.Conn) error {
	return nil
}

func (s *GatewayServer) OnDisconnect(c *connsvr.Conn, kickout bool) {
}

func (s *GatewayServer) OnRecvMessage(c *connsvr.Conn, msg []byte) {
	cmd := string(msg)
	switch cmd {
	case "kickout":
		c.Kickout()
	default:
		log.Println(cmd)
		fmt.Fprintln(c, cmd)
	}
}

func main() {
	s := NewGatewayServer("tcp", "localhost:8000")
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
