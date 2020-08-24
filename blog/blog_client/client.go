package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

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

	// create blog
	createBlog(c, "kin", "title1", "content1.")
	// get blog
	getBlog(c, "5f434c81984e155761a11e1e")
	getBlog(c, "5f434c81984e155761a11eww") // invalid hex
	getBlog(c, "5f434c81984e155761a11e11") // not exist
	// update blog
	updateBlog(c, "5f434c81984e155761a11e1e", "kinx", "titlex", "contentx.")
	// delete blog
	deleteBlog(c, "5f4353255b9e202d9fd0b6bc")
	deleteBlog(c, "5f4353255b9e202d9fd0b6bc") //expect not found
	// list blogs
	listBlogs(c)

}

func createBlog(c blogpb.BlogServiceClient, author string, title string, content string) {
	println("Create blog.")
	// create context with deadline
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	// defer cancel
	defer cancel()
	res, err := c.CreateBlog(ctx, &blogpb.CreateBlogRequest{
		Blog: &blogpb.Blog{
			Author:  author,
			Title:   title,
			Content: content,
		},
	})
	if err != nil {
		log.Fatalf("Fail to envoke CreateBlog.%v\n", err)
	}
	fmt.Printf("Create blog successfully.blog:%v.\n", res)
}

func getBlog(c blogpb.BlogServiceClient, id string) {
	println("Get blog")
	// create context with deadline
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	// defer cancel
	defer cancel()
	res, err := c.GetBlog(ctx, &blogpb.GetBlogRequest{
		Id: id,
	})
	if err != nil {
		fmt.Printf("Fail to envoke GetBlog.%v\n", err)
		return
	}
	fmt.Printf("Get blog successfully.blog:%v.\n", res)
}

func updateBlog(c blogpb.BlogServiceClient, id string, author string, title string, content string) {
	println("Update blog.")
	// create context with deadline
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	// defer cancel
	defer cancel()
	res, err := c.UpdateBlog(ctx, &blogpb.UpdateBlogRequest{
		Blog: &blogpb.Blog{
			Id:      id,
			Author:  author,
			Title:   title,
			Content: content,
		},
	})
	if err != nil {
		log.Fatalf("Fail to envoke UpdateBlog.%v\n", err)
	}
	fmt.Printf("Update blog successfully.blog:%v.\n", res)
}

func deleteBlog(c blogpb.BlogServiceClient, id string) {
	println("Delete blog")
	// create context with deadline
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	// defer cancel
	defer cancel()
	res, err := c.DeleteBlog(ctx, &blogpb.DeleteBlogRequest{
		Id: id,
	})
	if err != nil {
		fmt.Printf("Fail to envoke DeleteBlog.%v\n", err)
		return
	}
	fmt.Printf("Delete blog successfully.blog:%v.\n", res)
}

func listBlogs(c blogpb.BlogServiceClient) {
	println("List blogs")
	// create context with deadline
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	// defer cancel
	defer cancel()
	stream, err := c.ListBlog(ctx, &blogpb.ListRequest{})
	if err != nil {
		fmt.Printf("Fail to envoke ListBlog.%v\n", err)
		return
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			println("Streaming finished.")
			return
		}
		if err != nil {
			fmt.Printf("Fail to recv ListBlog.%v\n", err)
			return
		}
		fmt.Printf("Get blog successfully.blog:%v.\n", res)
	}

}
