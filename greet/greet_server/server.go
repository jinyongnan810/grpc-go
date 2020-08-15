package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/jinyongnan810/grpc-go/greet/greetpb"

	"google.golang.org/grpc"
)

type server struct{}

// unary
func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	firstName := req.GetGreeting().GetFirstName()
	result := "hello, " + firstName
	res := greetpb.GreetResponse{
		Result: result,
	}
	return &res, nil
}

// streaming server
func (*server) GreetedManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetedManyTimesServer) error {
	first := req.GetGreeting().GetFirstName()
	for i := 0; i < 10; i++ {
		result := "hello, " + first + ". No." + strconv.Itoa(i) + "\n"
		res := greetpb.GreetManyTimesResponse{
			Result: result,
		}
		err := stream.Send(&res)
		if err != nil {
			log.Fatal("Fail to send stream.", err)
		}
		time.Sleep(time.Millisecond * 200)
	}
	return nil
}

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
