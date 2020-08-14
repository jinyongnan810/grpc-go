package main

import (
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
	defer conn.Close() // when done, close connection
	fmt.Printf("Client created.%f", c)
}
