package main

import (
	"context"
	"flag"
	pb "github.com/arinto/hello-concurrent-grpc-go/internal/helloworld"
	otgrpc "github.com/opentracing-contrib/go-grpc"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"net"
	"time"
)

type server struct {
	id        string
	sleepms   int
	randomize bool
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	sleepms := s.sleepms
	if s.randomize {
		sleepms = rand.Intn(s.sleepms)
	}
	log.Printf("Received %+v, sleepms for %+v ms\n", in.Name, sleepms)
	time.Sleep(time.Duration(sleepms) * time.Millisecond)
	log.Println("Done!")
	return &pb.HelloReply{Message: s.id + ": Hello " + in.Name}, nil
}

func main() {
	address := flag.String("address", ":50051", "Address of the gRPC server")
	sleepms := flag.Int("sleepms", 100, "Duration of sleepms in millis")
	randomize := flag.Bool("randomize", false, "Set to true to randomize sleepms")

	flag.Parse()
	rand.Seed(time.Now().Unix())

	lis, err := net.Listen("tcp", *address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(grpc.UnaryInterceptor(
		otgrpc.OpenTracingServerInterceptor(opentracing.GlobalTracer(), otgrpc.LogPayloads())))
	pb.RegisterGreeterServer(s, &server{id: *address, sleepms: *sleepms, randomize: *randomize})
	log.Println("Start to serve, listening at " + *address)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
