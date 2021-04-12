package main

import (
	"context"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"

	"github.com/ironzhang/projects/testsrc/grpc/balancer"
	"github.com/ironzhang/projects/testsrc/grpc/protobuf/greet"
	_ "github.com/ironzhang/projects/testsrc/grpc/resolver"
)

func main() {
	const target = "registry://seeds/localhost:50051,localhost:50052"
	conn, err := grpc.Dial(target, grpc.WithInsecure(), grpc.WithBalancerName(balancer.Name))
	if err != nil {
		log.Fatalf("failed to dial to %s: %v", target, err)
	}
	defer conn.Close()
	c := greet.NewGreeterClient(conn)

	name := "world"
	if len(os.Args) > 1 {
		name = os.Args[1]
	}

	for i := 0; i < 500; i++ {
		SayHello(c, name)
		time.Sleep(time.Second)
	}
}

func SayHello(c greet.GreeterClient, name string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.SayHello(ctx, &greet.HelloRequest{Name: name})
	if err != nil {
		log.Printf("failed to say hello: %v", err)
		return
	}
	log.Printf("Greeting: %s", r.Mesg)
}
