package main

import (
	"fmt"
	"log"

	"github.com/jinyongnan810/grpc-go/blog/blogpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("hello grpc client")
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalln("Fail to dail.", err)
	}
	c := blogpb.NewBlogServiceClient(conn)
	fmt.Println("Client created.", c)
}
