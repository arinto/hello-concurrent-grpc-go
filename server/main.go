package main

import (
	"context"
	"flag"
	pb "github.com/arinto/hello-concurrent-grpc-go/internal/helloworld"
	otgrpc "github.com/opentracing-contrib/go-grpc"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

type server struct{
	id string
	sleep time.Duration
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received %+v, sleep for %+v\n", in.Name, s.sleep)
	time.Sleep(s.sleep)
	return &pb.HelloReply{Message: s.id + ": Hello " + in.Name}, nil
}

func main() {
	address := flag.String("address", ":50051", "Address of the gRPC server")
	sleep := flag.Duration("sleepms", 100*time.Millisecond, "Duration of sleep in millis")

	flag.Parse()

	lis, err := net.Listen("tcp", *address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(grpc.UnaryInterceptor(
		otgrpc.OpenTracingServerInterceptor(opentracing.GlobalTracer(), otgrpc.LogPayloads())))
	pb.RegisterGreeterServer(s, &server{id: *address, sleep: *sleep})
	log.Println("Start to serve, listening at " + *address)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
