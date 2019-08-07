package main

import (
	"encoding/json"
	"flag"
	"log"
	"net"
	"sync"
	"time"

	"github.com/ironzhang/practice/x-test/udp/proto"
)

func DialUDP(network, address string) (*net.UDPConn, error) {
	addr, err := net.ResolveUDPAddr(network, address)
	if err != nil {
		return nil, err
	}
	return net.DialUDP(network, nil, addr)
}

func RunTest(n int, addr string, ch chan proto.Message) {
	conn, err := DialUDP("udp", addr)
	if err != nil {
		log.Fatalf("dial udp: %v", err)
		return
	}
	conn.SetReadBuffer(1024 * 1024)
	conn.SetWriteBuffer(1024 * 1024)

	buf := make([]byte, 1500)
	for i := 0; i < n; i++ {
		msg := proto.Message{ClientSend: time.Now()}
		data, err := json.Marshal(msg)
		if err != nil {
			log.Printf("marshal: %v", err)
			return
		}
		if _, err = conn.Write(data); err != nil {
			log.Printf("write: %v", err)
			return
		}

		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("read: %v", err)
			return
		}
		if err = json.Unmarshal(buf[:n], &msg); err != nil {
			log.Printf("unmarshal: %v", err)
			return
		}
		msg.ClientRecv = time.Now()

		ch <- msg
	}
}

func main() {
	var c, n int
	var addr string
	var t1 time.Duration
	var t2 time.Duration

	flag.IntVar(&c, "c", 100, "c")
	flag.IntVar(&n, "n", 100, "n")
	flag.DurationVar(&t1, "t1", 0, "t1")
	flag.DurationVar(&t2, "t2", 0, "t2")
	flag.StringVar(&addr, "addr", "localhost:2000", "addr")
	flag.Parse()

	var wg sync.WaitGroup
	ch := make(chan proto.Message, c*n)
	for i := 0; i < c; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			RunTest(n, addr, ch)
		}()
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	var cnt int
	var sum1 time.Duration
	var sum2 time.Duration
	for msg := range ch {
		d1 := msg.ServerRecv.Sub(msg.ClientSend)
		d2 := msg.ClientRecv.Sub(msg.ServerSend)
		if (t1 > 0 && d1 >= t1) || (t2 > 0 && d2 >= t2) {
			log.Printf("%v: %v, %v", msg, d1, d2)
		}
		cnt++
		sum1 += d1
		sum2 += d2
	}

	log.Printf("ave: %v, %v", sum1/time.Duration(cnt), sum2/time.Duration(cnt))
}
