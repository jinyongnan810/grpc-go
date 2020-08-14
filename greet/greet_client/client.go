package main

import (
	"context"
	"fmt"
	"log"

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

	defer conn.Close() // when done, close connection

}

// unary req/res
func doUnary(c greetpb.GreetServiceClient) {
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
