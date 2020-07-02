package main

import (
	"context"
	"flag"
	"log"
	"net"

	"github.com/ironzhang/projects/testsrc/grpc/protobuf/greet"
	"google.golang.org/grpc"
)

type Greeter struct{}

func (p *Greeter) SayHello(ctx context.Context, req *greet.HelloRequest) (*greet.HelloReply, error) {
	log.Printf("Received: %v", req.Name)
	return &greet.HelloReply{Mesg: "Hello " + req.Name}, nil
}

func main() {
	addr := flag.String("addr", ":50051", "listen address")
	flag.Parse()

	ln, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	greet.RegisterGreeterServer(s, &Greeter{})

	log.Printf("grpc serve on %v", *addr)
	if err = s.Serve(ln); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
