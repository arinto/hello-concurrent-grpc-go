package main

import (
	"context"
	"flag"
	pb "github.com/arinto/hello-concurrent-grpc-go/internal/helloworld"
	"google.golang.org/grpc"
	"log"
	"os"
	"sync"
	"time"
)

func asyncHello(ctx context.Context, client pb.GreeterClient,  name *string, resChan chan *pb.HelloReply, errChan chan error) {
	defer close(resChan)
	defer close(errChan)

	r, err := client.SayHello(ctx, &pb.HelloRequest{Name: *name})
	if err != nil {
		log.Println(err)
		errChan <- err
	}
	resChan <- r
}

func waitHello(helloId string, wg *sync.WaitGroup,  ctx context.Context, resChan chan *pb.HelloReply, errChan chan error) {
	select {
	case r := <- resChan:
		log.Printf("Greeting: %s to %s", r.Message, helloId)
	case err := <- errChan:
		log.Printf("Error calling helloworld %s: %+v", helloId, err)
	case <-ctx.Done():
		log.Printf("Timeout calling helloworld %s", helloId)
	}
	wg.Done()
}

func anExperiment(expId int, timeout1 *time.Duration, timeout2 *time.Duration, c1 pb.GreeterClient, c2 pb.GreeterClient, address1 *string, address2 *string, name *string){
	ctx := context.Background()
	wg := sync.WaitGroup{}
	wg.Add(2)

	// First async grpc call
	ctx1, cancel1 := context.WithTimeout(ctx, *timeout1)
	defer cancel1()
	resChan1 := make(chan *pb.HelloReply, 1)
	errChan1 := make(chan error, 1)
	go asyncHello(ctx1, c1, name, resChan1, errChan1)
	go waitHello(*address1, &wg, ctx1, resChan1, errChan1)

	// Second async grpc call
	ctx2, cancel2 := context.WithTimeout(ctx, *timeout2)
	defer cancel2()
	resChan2 := make(chan *pb.HelloReply, 1)
	errChan2 := make(chan error, 1)
	go asyncHello(ctx2, c2, name, resChan2, errChan2)
	go waitHello(*address2, &wg, ctx2, resChan2, errChan2)
	wg.Wait()
	time.Sleep(100 * time.Millisecond)
	log.Printf("Done for experiment: %d\n", expId)
}

func main() {
	address1 := flag.String("address1", "localhost:50051", "Address of the 1st gRPC server")
	address2 := flag.String("address2", "localhost:50052", "Address of the 2nd gRPC server")
	timeout1 := flag.Duration("timeout1", 200*time.Millisecond, "Timeout for 1st gRPC call")
	timeout2 := flag.Duration("timeout2", 200*time.Millisecond, "Timeout for 2nd gRPC call")
	n := flag.Int("n", 10, "Number of experiment")

	name := flag.String("name", "concurrent", "Name to be greet")
	flag.Parse()

	log.SetOutput(os.Stdout)

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

	for i := 0; i < *n; i++ {
		anExperiment(i, timeout1, timeout2, c1, c2, address1, address2, name)
	}

}
