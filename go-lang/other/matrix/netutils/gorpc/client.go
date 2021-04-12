package gorpc

import (
	"errors"
	"io"
	"net/rpc"
	"sync"
)

func Dial(net, addr string) (*Client, error) {
	c := &Client{calls: make(chan *rpc.Call, 20)}
	go calling(net, addr, c.calls)
	return c, nil
}

type Client struct {
	mu    sync.RWMutex
	calls chan *rpc.Call
}

func (c *Client) Close() error {
	c.mu.Lock()
	if c.calls != nil {
		close(c.calls)
		c.calls = nil
	}
	c.mu.Unlock()
	return nil
}

func (c *Client) Call(serviceMethod string, args interface{}, reply interface{}) error {
	call := <-c.Go(serviceMethod, args, reply, make(chan *rpc.Call, 1)).Done
	return call.Error
}

func (c *Client) Go(serviceMethod string, args interface{}, reply interface{}, done chan *rpc.Call) *rpc.Call {
	if done == nil {
		done = make(chan *rpc.Call, 1)
	}
	call := &rpc.Call{
		ServiceMethod: serviceMethod,
		Args:          args,
		Reply:         reply,
		Done:          done,
	}

	c.mu.RLock()
	if c.calls != nil {
		c.calls <- call
	} else {
		go func() {
			call.Error = errors.New("client is closed")
			call.Done <- call
		}()
	}
	c.mu.RUnlock()

	return call
}

func calling(net, addr string, calls <-chan *rpc.Call) {
	var err error
	var client *rpc.Client
	for call := range calls {
		if client == nil {
			if client, err = rpc.Dial(net, addr); err != nil {
				call.Error = err
				call.Done <- call
				continue
			}
		}
		if err = client.Call(call.ServiceMethod, call.Args, call.Reply); err == rpc.ErrShutdown || err == io.ErrUnexpectedEOF {
			client = nil
		}
		call.Error = err
		call.Done <- call
	}
	if client != nil {
		client.Close()
	}
}
