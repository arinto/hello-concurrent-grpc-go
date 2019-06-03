package main

import (
	"context"
	"flag"
	pb "github.com/arinto/hello-concurrent-grpc-go/internal/helloworld"
	"google.golang.org/grpc"
	"log"
	"time"
)

func main() {
	address1 := flag.String("address1", "localhost:50051", "Address of the 1st gRPC server")
	address2 := flag.String("address2", "localhost:50052", "Address of the 2nd gRPC server")
	deadline1 := flag.Duration("deadline1", 200*time.Millisecond, "Deadline for 1st gRPC call")
	deadline2 := flag.Duration("deadline2", 200*time.Millisecond, "Deadline for 2nd gRPC call")

	name := flag.String("name", "hello-concurrent", "Name to be greet")
	flag.Parse()

	conn1, err := grpc.Dial(*address1, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn1.Close()
	c1 := pb.NewGreeterClient(conn1)
	// Contact the server and print out its response.

	conn2, err := grpc.Dial(*address2, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn2.Close()
	c2 := pb.NewGreeterClient(conn2)

	parentCtx := context.Background()

	// First grpc call
	ctx1, cancel := context.WithTimeout(parentCtx, *deadline1)
	defer cancel()
	r1, err := c1.SayHello(ctx1, &pb.HelloRequest{Name: *name})
	if err != nil {
		log.Fatalf("1: could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r1.Message)

	//Second grpc call
	ctx2, cancel := context.WithTimeout(parentCtx, *deadline2)
	defer cancel()
	r2, err := c2.SayHello(ctx2, &pb.HelloRequest{Name: *name})
	if err != nil {
		log.Fatalf("2: could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r2.Message)

}
