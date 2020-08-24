package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/jinyongnan810/grpc-go/blog/blogpb"

	"google.golang.org/grpc"
)

type server struct{}

func main() {
	// to show error file and line
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	fmt.Println("hello go server.")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatal("Fail to listen.", err)
	}
	s := grpc.NewServer()
	blogpb.RegisterBlogServiceServer(s, &server{})
	go func() {
		println("Server starting.")
		if err := s.Serve(lis); err != nil {
			log.Fatal("Fail to serve.", err)
		}
	}()

	// block until interrupted
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
	// stop server and listener
	println("Stopping server...")
	s.Stop()
	println("Stopping listener...")
	lis.Close()
	println("Server stopped.")

}
