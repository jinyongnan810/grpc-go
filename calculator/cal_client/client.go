package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/grpc/status"

	"github.com/jinyongnan810/grpc-go/calculator/calpb"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("hello grpc client")
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalln("Fail to dail.", err)
	}
	c := calpb.NewCalculatorServiceClient(conn)
	fmt.Println("Client created.", c)
	// running unary
	runUnary(c)
	// running server stream
	runServerStream(c)
	// running client stream
	runClientStream(c)
	// running bi stream
	runBiStream(c)

	// running square root
	runSquareRoot(c, 7)
	runSquareRoot(c, -1)
}

func runUnary(c calpb.CalculatorServiceClient) {
	print("-----running unary-----\n")
	req := &calpb.UnaryRequest{
		First:  7,
		Second: 3,
	}
	res, err := c.UnarySum(context.Background(), req)
	if err != nil {
		log.Fatalln("Fail to call UnarySum.", err)
	}
	println("Unary res is:", fmt.Sprint((res.GetResult())))
}

func runServerStream(c calpb.CalculatorServiceClient) {
	print("-----running server stream-----\n")
	req := &calpb.ServerStreamRequest{
		Number: int32(7),
	}
	stream, err := c.ServerStreamDecomposite(context.Background(), req)

	if err != nil {
		log.Fatalln("Fail to call ServerStreamDecomposite.", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			println("Server streaming done.")
			break
		}
		if err != nil {
			log.Fatalln("Fail to receive ServerStreamDecomposite response.", err)
			break
		}
		println("Server stream res is:", fmt.Sprint((res.GetDecomposition())))
	}

}
func runClientStream(c calpb.CalculatorServiceClient) {
	print("-----running client stream-----\n")
	reqs := []int32{
		5, 6, 7, 9,
	}
	stream, err := c.CLientStreamAverage(context.Background())
	if err != nil {
		log.Fatalln("Fail to call CLientStreamAverage.", err)
	}
	for _, req := range reqs {
		stream.Send(&calpb.ClientStreamRequest{
			Number: req,
		})
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalln("Fail to receive CLientStreamAverage response.", err)
	}
	println("Client stream res is:", fmt.Sprint((res.GetAverage())))

}

func runBiStream(c calpb.CalculatorServiceClient) {
	print("-----running bi stream-----\n")
	reqs := []int32{
		1, 5, 3, 6, 2, 20,
	}
	waitc := make(chan int, 2)

	stream, err := c.BiStreamFindMax(context.Background())
	if err != nil {
		log.Fatalln("Fail to call BiStreamFindMax.", err)
	}
	go func() {
		for _, num := range reqs {
			stream.Send(&calpb.BiStreamRequest{
				Number: num,
			})
			println("client sent:", fmt.Sprint(num))
			time.Sleep(time.Second * 1)
		}
		stream.CloseSend()
		waitc <- 1
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				println("server stream ended.")
				break
			}
			if err != nil {
				log.Fatalln("Fail to recv server stream.", err)
			}
			println("recv from server:", fmt.Sprint(res.GetMax()))
		}
		waitc <- 1

	}()
	for i := 0; i < 2; i++ {
		<-waitc
	}
	println("bi streaming done.")
}
func runSquareRoot(c calpb.CalculatorServiceClient, num int32) {
	fmt.Printf("-----running square root %v -----\n", num)
	req := &calpb.SquareRootRequest{
		Number: num,
	}
	res, err := c.SquareRoot(context.Background(), req)
	if err != nil {
		fmt.Println("Fail to call SquareRoot.", err)
		errDetail, ok := status.FromError(err)
		if ok {
			fmt.Println("User error message:", errDetail.Message())
			fmt.Printf("User error code:%v\n", errDetail.Code())
		} else {
			log.Fatalln("Fatal error", err)
		}
		return
	}
	println("Square root is:", fmt.Sprint((res.GetRoot())))
}
