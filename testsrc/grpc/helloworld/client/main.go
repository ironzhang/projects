package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/ironzhang/projects/testsrc/grpc/protobuf/greet"
	"google.golang.org/grpc"
)

func main() {
	const addr = "localhost:50051"
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to dial to %s: %v", addr, err)
	}
	defer conn.Close()
	c := greet.NewGreeterClient(conn)

	name := "world"
	if len(os.Args) > 1 {
		name = os.Args[1]
	}

	for i := 0; i < 5; i++ {
		SayHello(c, name)
	}
}

func SayHello(c greet.GreeterClient, name string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.SayHello(ctx, &greet.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("failed to say hello: %v", err)
	}
	log.Printf("Greeting: %s", r.Mesg)
}
