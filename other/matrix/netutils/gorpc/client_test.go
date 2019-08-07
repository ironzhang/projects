package gorpc

import (
	"net"
	"net/rpc"
	"testing"
)

type Args struct {
	A, B int
}

type Arith struct{}

func (p *Arith) Add(args *Args, reply *int) error {
	*reply = args.A + args.B
	return nil
}

func NewArithServer() (*rpc.Server, error) {
	svr := rpc.NewServer()
	if err := svr.Register(&Arith{}); err != nil {
		return nil, err
	}
	return svr, nil
}

func RunArithServer(t *testing.T, network, address string) {
	svr, err := NewArithServer()
	if err != nil {
		t.Fatal(err)
	}
	ln, err := net.Listen(network, address)
	if err != nil {
		t.Fatal(err)
	}
	go svr.Accept(ln)
}

func ArithAdd(c *Client, a, b int) (int, error) {
	var result int
	if err := c.Call("Arith.Add", &Args{A: a, B: b}, &result); err != nil {
		return 0, err
	}
	return result, nil
}

func TestClientWithServer(t *testing.T) {
	network, address := "tcp", "localhost:2000"
	RunArithServer(t, network, address)

	c, err := Dial(network, address)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	a, b := 1, 2
	got, err := ArithAdd(c, a, b)
	if err != nil {
		t.Fatal(err)
	}
	if got != a+b {
		t.Errorf("%d != ArithAdd(%d, %d)", got, a, b)
	} else {
		t.Logf("%d == ArithAdd(%d, %d)", got, a, b)
	}
}

func TestClientWithoutServer(t *testing.T) {
	network, address := "tcp", "localhost:2001"

	c, err := Dial(network, address)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	_, err = ArithAdd(c, 0, 0)
	if err == nil {
		t.Fatal("ArithAdd success without server")
	} else {
		t.Log(err)
	}

	RunArithServer(t, network, address)

	a, b := 1, 2
	got, err := ArithAdd(c, a, b)
	if err != nil {
		t.Fatal(err)
	}
	if got != a+b {
		t.Errorf("%d != ArithAdd(%d, %d)", got, a, b)
	} else {
		t.Logf("%d == ArithAdd(%d, %d)", got, a, b)
	}
}
