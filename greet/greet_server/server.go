package main

import (
	"fmt"
	"log"
	"net"

	"github.com/jinyongnan810/grpc-go/greet/greetpb"

	"google.golang.org/grpc"
)

type server struct{}

func main() {
	fmt.Println("hello go server.")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatal("Fail to listen.", err)
	}
	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatal("Fail to serve.", err)
	}
}
