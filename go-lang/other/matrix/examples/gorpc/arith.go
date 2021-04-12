package main

type Args struct {
	A, B int
}

type Arith struct{}

func (p *Arith) Add(args *Args, reply *int) error {
	*reply = args.A + args.B
	return nil
}
