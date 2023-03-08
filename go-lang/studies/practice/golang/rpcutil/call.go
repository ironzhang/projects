package rpcutil

import (
	"errors"
	"net/rpc"
	"time"
)

func Call(client *rpc.Client, method string, args, reply interface{}, timeout time.Duration) (err error) {
	if client == nil {
		return errors.New("invalid client")
	}

	ch := make(chan bool)
	go func() {
		err = client.Call(method, args, reply)
		ch <- true
	}()

	select {
	case <-time.After(timeout):
		err = errors.New("call timeout")
	case <-ch:
	}

	return err
}
