package arith

import (
	"context"
	"errors"
)

type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

type Arith int

func (t Arith) Add(ctx context.Context, args Args, reply *int) error {
	*reply = args.A + args.B
	return nil
}

func (t Arith) Sub(ctx context.Context, args Args, reply *int) error {
	*reply = args.A - args.B
	return nil
}

func (t Arith) Multiply(ctx context.Context, args Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func (t Arith) Divide(ctx context.Context, args Args, quo *Quotient) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	quo.Quo = args.A / args.B
	quo.Rem = args.A % args.B
	return nil
}
