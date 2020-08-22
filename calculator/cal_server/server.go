package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/jinyongnan810/grpc-go/calculator/calpb"

	"google.golang.org/grpc"
)

type server struct{}

func (*server) UnarySum(c context.Context, req *calpb.UnaryRequest) (*calpb.UnaryResponse, error) {
	first := req.GetFirst()
	second := req.GetSecond()
	res := calpb.UnaryResponse{
		Result: first + second,
	}
	return &res, nil
}
func (*server) ServerStreamDecomposite(req *calpb.ServerStreamRequest, stream calpb.CalculatorService_ServerStreamDecompositeServer) error {
	original := req.GetNumber()
	num := original
	var i int32
	i = 1
	sent := []int32{}
	for {
		if num%i == 0 && i != 1 {
			err := stream.Send(&calpb.ServerStreamResponse{
				Decomposition: i,
			})
			if err != nil {
				log.Fatal("Fail to send server stream.", err)
				return err
			}
			sent = append(sent, i)
			num = num / i
		} else {
			i++
		}
		if i >= num {
			if len(sent) == 0 {
				err := stream.Send(&calpb.ServerStreamResponse{
					Decomposition: 1,
				})
				if err != nil {
					log.Fatal("Fail to send server stream.", err)
					return err
				}
			}
			err2 := stream.Send(&calpb.ServerStreamResponse{
				Decomposition: num,
			})
			if err2 != nil {
				log.Fatal("Fail to send server stream.", err2)
				return err2
			}
			break
		}

	}
	return nil
}

func (*server) CLientStreamAverage(stream calpb.CalculatorService_CLientStreamAverageServer) error {
	var sum, count float32
	sum = 0
	count = 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			println("done receiving client stream")
			break
		}
		if err != nil {
			log.Fatal("Fail to recv client stream.", err)
			break
		}
		sum += float32(req.GetNumber())
		count++
	}
	stream.SendAndClose(&calpb.ClientStreamResponse{
		Average: sum / count,
	})
	return nil
}

func (*server) BiStreamFindMax(stream calpb.CalculatorService_BiStreamFindMaxServer) error {
	var max int32
	max = 0
	waitc := make(chan int, 2)

	// recv
	go func() {
		for {
			req, err := stream.Recv()
			if err == io.EOF {
				println("done receive client bi stream.")
				break
			}
			if err != nil {
				log.Fatal("Fail to recv client bi stream.", err)
				break
			}
			if max < req.GetNumber() {
				max = req.GetNumber()
				stream.Send(&calpb.BiStreamResponse{
					Max: max,
				})
			}
		}
		waitc <- 1
	}()
	<-waitc
	return nil
}

func main() {
	fmt.Println("hello go server.")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatal("Fail to listen.", err)
	}
	s := grpc.NewServer()
	calpb.RegisterCalculatorServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatal("Fail to serve.", err)
	}
}
