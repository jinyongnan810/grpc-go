package main

import (
	"context"
	"fmt"
	"io"
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

// streaming client
func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	result := ""
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		} else if err != nil {
			log.Fatal("Fail to receive client stream.", err)
			break
		} else {
			result += "Hello, " + req.GetGreeting().GetFirstName() + "!"
		}
	}
	return nil
}

// bi direction streaming server
func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	waitc := make(chan int, 2)
	// recv
	go func() {
		for {
			req, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal("Fail to receive client stream.", err)
				break
			}
			println("server received client stream:", req.GetGreeting().GetFirstName())
		}
		waitc <- 1
	}()
	// send
	go func() {
		for i := 0; i < 5; i++ {
			err := stream.Send(&greetpb.GreetEveryoneResponse{
				Result: "server sends hello, No." + strconv.Itoa(i),
			})
			if err != nil {
				log.Fatal("Fail to send server stream.", err)
				break
			}
			time.Sleep(1 * time.Second)
		}
		waitc <- 2

	}()

	// block
	for i := 0; i < 2; i++ {
		<-waitc
	}
	close(waitc)
	println("done bi streaming.")
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
