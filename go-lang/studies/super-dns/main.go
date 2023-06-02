package main

import (
	"errors"
	"net"

	"github.com/ironzhang/tlog"
	"golang.org/x/net/dns/dnsmessage"
)

func HandleMessage(msg dnsmessage.Message) (res dnsmessage.Message, err error) {
	if msg.Response {
		return dnsmessage.Message{}, errors.New("can not process response message")
	}

	if len(msg.Questions) <= 0 {
		return dnsmessage.Message{}, errors.New("can not process a none question message")
	}

	res.ID = msg.ID
	res.Response = true
	//res.Authoritative = true
	res.RCode = dnsmessage.RCodeSuccess

	for _, q := range msg.Questions {
		switch q.Type {
		case dnsmessage.TypeA:
			domain := q.Name.String()
			if domain == "www.baidu.com." {
				tlog.Infof("process domain %s %s", domain, q.Type.String())
				ans := dnsmessage.Resource{
					Header: dnsmessage.ResourceHeader{
						Name:  q.Name,
						Type:  dnsmessage.TypeA,
						Class: dnsmessage.ClassINET,
						TTL:   600,
					},
					Body: &dnsmessage.AResource{
						A: [4]byte{128, 0, 0, 1},
					},
				}
				res.Answers = append(res.Answers, ans)
			} else {
				res.RCode = dnsmessage.RCodeNameError
			}
			res.Questions = append(res.Questions, q)

		case dnsmessage.TypeAAAA:
			res.RCode = dnsmessage.RCodeNameError
			res.Questions = append(res.Questions, q)

		default:
			res.RCode = dnsmessage.RCodeNameError
			res.Questions = append(res.Questions, q)
		}
	}

	return res, nil
}

func process(ln net.PacketConn, peer net.Addr, buf []byte) {
	var err error
	var msg dnsmessage.Message
	if err = msg.Unpack(buf); err != nil {
		tlog.Errorw("unpack", "error", err)
		return
	}

	res, err := HandleMessage(msg)
	if err != nil {
		tlog.Errorw("handle message", "error", err)
		return
	}

	nbuf, err := res.Pack()
	if err != nil {
		tlog.Errorw("pack", "error", err)
		return
	}

	n, err := ln.WriteTo(nbuf, peer)
	if err != nil {
		tlog.Errorw("write to", "error", err)
		return
	}
	_ = n
	//tlog.Infof("write %v, n=%d", res.GoString(), n)
}

func main() {
	addr, err := net.ResolveUDPAddr("udp", ":53")
	if err != nil {
		tlog.Errorw("net resolve udp addr", "error", err)
		return
	}
	ln, err := net.ListenUDP("udp", addr)
	if err != nil {
		tlog.Errorw("net listen udp", "error", err)
		return
	}
	defer ln.Close()
	tlog.Infof("serve on %s", ":53")

	buf := make([]byte, 1500)
	for {
		n, peer, err := ln.ReadFrom(buf)
		if err != nil {
			tlog.Errorw("read from", "error", err)
			continue
		}
		go process(ln, peer, buf[:n])
	}
}
