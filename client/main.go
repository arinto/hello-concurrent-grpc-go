package main

import (
	"context"
	"flag"
	pb "github.com/arinto/hello-concurrent-grpc-go/internal/helloworld"
	"google.golang.org/grpc"
	"log"
	"time"
)

func hello(ctx context.Context, client pb.GreeterClient,  name *string, resChan chan *pb.HelloReply, errChan chan error) {
	defer close(resChan)
	defer close(errChan)

	r, err := client.SayHello(ctx, &pb.HelloRequest{Name: *name})
	if err != nil {
		log.Println(err)
		log.Printf("Sending to err channel: %+v\n", r)
		errChan <- err
	}
	log.Printf("Sending to result channel: %+v\n", r)
	resChan <- r
}

func main() {
	address1 := flag.String("address1", "localhost:50051", "Address of the 1st gRPC server")
	address2 := flag.String("address2", "localhost:50052", "Address of the 2nd gRPC server")
	timeout1 := flag.Duration("timeout1", 200*time.Millisecond, "Timeout for 1st gRPC call")
	timeout2 := flag.Duration("timeout2", 200*time.Millisecond, "Timeout for 2nd gRPC call")

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

	ctx := context.Background()

	// First async grpc call
	ctx1, cancel1 := context.WithTimeout(ctx, *timeout1)
	defer cancel1()
	resChan1 := make(chan *pb.HelloReply, 1)
	errChan1 := make(chan error, 1)
	go hello(ctx1, c1, name, resChan1, errChan1)

	// Second async grpc call
	ctx2, cancel2 := context.WithTimeout(ctx, *timeout2)
	defer cancel2()
	resChan2 := make(chan *pb.HelloReply, 1)
	errChan2 := make(chan error, 1)
	go hello(ctx2, c2, name, resChan2, errChan2)


	select {
	case <-ctx1.Done():
		log.Println("Timeout calling helloworld1")
	case error1 := <- errChan1:
		log.Printf("Error calling helloworld1: %+v", error1)
	case r1 := <- resChan1:
		log.Printf("Greeting: %s", r1.Message)
	}

	select {
	case <-ctx2.Done():
		log.Println("Timeout calling helloworld2")
	case error2 := <- errChan2:
		log.Printf("Error calling helloworld2: %+v", error2)
	case r2 := <- resChan2:
		log.Printf("Greeting: %s", r2.Message)
	}

	time.Sleep(1*time.Second)

}
