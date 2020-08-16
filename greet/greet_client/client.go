package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/jinyongnan810/grpc-go/greet/greetpb"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("hello grpc client")
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalln("Fail to dail.", err)
	}
	c := greetpb.NewGreetServiceClient(conn)
	fmt.Println("Client created.", c)
	// do unary
	doUnary(c)
	// do streaming server
	doStreamingServer(c)
	// do streaming client
	doStreamingClient(c)
	// do bi directional streaming
	doBiStreamingClient(c)

	defer conn.Close() // when done, close connection

}

// unary req/res
func doUnary(c greetpb.GreetServiceClient) {
	print("-----running unary-----\n")
	// send request
	greeting := greetpb.Greeting{
		FirstName: "test first name",
		LastName:  "test last name",
	}
	req := greetpb.GreetRequest{
		Greeting: &greeting,
	}
	res, err := c.Greet(context.Background(), &req)
	if err != nil {
		log.Fatalln("Fail to envoke Greet.", err)
	}
	fmt.Println("Response is :", res.GetResult())
}

// client of streaming server
func doStreamingServer(c greetpb.GreetServiceClient) {
	print("-----running server streaming-----\n")
	greeting := greetpb.Greeting{
		FirstName: "test first name",
		LastName:  "test last name",
	}
	req := greetpb.GreetManyTimesRequest{
		Greeting: &greeting,
	}
	client, err := c.GreetedManyTimes(context.Background(), &req)
	if err != nil {
		log.Fatalln("Fail to envoke GreetedManyTimes.", err)
	}
	for {
		res, err := client.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalln("Fail to receive stream.", err)
		} else {
			print(res.GetResult())
		}
	}

}

// streaming client
func doStreamingClient(c greetpb.GreetServiceClient) {
	print("-----running client streaming-----\n")
	reqDatas := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Jack",
				LastName:  "test last name",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Tom",
				LastName:  "test last name",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Lili",
				LastName:  "test last name",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Jill",
				LastName:  "test last name",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Mark",
				LastName:  "test last name",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Anni",
				LastName:  "test last name",
			},
		},
	}
	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalln("Fail to call LongGreet.", err)
	}
	for _, reqData := range reqDatas {
		err := stream.Send(reqData)
		if err != nil {
			log.Fatalln("Fail to send client stream.", err)
		}
	}
	print("sending client stream done.")
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalln("Fail to receive server response.", err)
	}
	print(res.GetResult())

}

// bi direction streaming client
func doBiStreamingClient(c greetpb.GreetServiceClient) {
	println("-----running bi directional streaming client-----")
	reqDatas := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Jack",
				LastName:  "test last name",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Tom",
				LastName:  "test last name",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Lili",
				LastName:  "test last name",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Jill",
				LastName:  "test last name",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Mark",
				LastName:  "test last name",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Anni",
				LastName:  "test last name",
			},
		},
	}
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalln("Fail to call GreetEveryone.", err)
	}
	waitc := make(chan int, 2)
	// recv
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal("Fail to receive server stream.", err)
				break
			}
			println("client received server stream:", res.GetResult())
		}
		waitc <- 1
	}()
	// send
	go func() {
		for _, reqData := range reqDatas {
			err := stream.Send(reqData)
			if err != nil {
				log.Fatalln("Fail to send client stream.", err)
				break
			}
			println("client sent:", reqData.GetGreeting().GetFirstName())
			time.Sleep(2 * time.Second)
		}
		stream.CloseSend()
		waitc <- 2
	}()
	// block
	for i := 0; i < 2; i++ {
		<-waitc
	}
	close(waitc)
	print("done bi streaming.")
}
