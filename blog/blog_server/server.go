package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"google.golang.org/grpc/reflection"

	"go.mongodb.org/mongo-driver/bson"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/jinyongnan810/grpc-go/blog/blogpb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"google.golang.org/grpc"
)

var collection *mongo.Collection

type server struct{}

type blogItem struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Author  string             `bson:"author"`
	Title   string             `bson:"title"`
	Content string             `bson:"content"`
}

func (*server) CreateBlog(c context.Context, req *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {
	blog := req.GetBlog()
	data := blogItem{
		Author:  blog.GetAuthor(),
		Title:   blog.GetTitle(),
		Content: blog.GetContent(),
	}
	insertRes, insertErr := collection.InsertOne(context.Background(), data)
	if insertErr != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Fail to insert to mongo db.%v\n", insertErr))
	}
	insertID, ok := insertRes.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Fail to convert oid.%v\n", insertErr))
	}
	res := &blogpb.CreateBlogResponse{
		Blog: &blogpb.Blog{
			Id:      insertID.Hex(),
			Author:  blog.GetAuthor(),
			Title:   blog.GetTitle(),
			Content: blog.GetContent(),
		},
	}
	return res, nil
}
func (*server) GetBlog(c context.Context, req *blogpb.GetBlogRequest) (*blogpb.GetBlogResponse, error) {
	id := req.GetId()
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Fail to convert id from hex.%v", err))
	}
	data := &blogItem{}
	filter := bson.M{"_id": oid}
	dbres := collection.FindOne(context.Background(), filter)
	if err := dbres.Decode(data); err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Fail to decode data.%v", err))
	}
	return &blogpb.GetBlogResponse{
		Blog: &blogpb.Blog{
			Id:      id,
			Author:  data.Author,
			Title:   data.Title,
			Content: data.Content,
		},
	}, nil
}

func (*server) UpdateBlog(c context.Context, req *blogpb.UpdateBlogRequest) (*blogpb.UpdateBlogResponse, error) {
	id := req.GetBlog().GetId()
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Fail to convert id from hex.%v", err))
	}
	data := &blogItem{}
	filter := bson.M{"_id": oid}
	dbres := collection.FindOne(context.Background(), filter)
	if err := dbres.Decode(data); err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Fail to decode data.%v", err))
	}

	// replace
	data.Author = req.GetBlog().GetAuthor()
	data.Title = req.GetBlog().GetTitle()
	data.Content = req.GetBlog().GetContent()
	_, updateErr := collection.ReplaceOne(context.Background(), filter, data)
	if updateErr != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Fail to update blog.%v", updateErr))
	}

	return &blogpb.UpdateBlogResponse{
		Blog: &blogpb.Blog{
			Id:      id,
			Author:  data.Author,
			Title:   data.Title,
			Content: data.Content,
		},
	}, nil
}

func (*server) DeleteBlog(c context.Context, req *blogpb.DeleteBlogRequest) (*blogpb.DeleteBlogResponse, error) {
	id := req.GetId()
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Fail to convert id from hex.%v", err))
	}
	filter := bson.M{"_id": oid}
	dbres, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Fail to delete data.%v", err))
	}
	if dbres.DeletedCount == 0 {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Blog not exist.%v", err))
	}
	return &blogpb.DeleteBlogResponse{
		Id: req.GetId(),
	}, nil
}
func (*server) ListBlog(req *blogpb.ListRequest, stream blogpb.BlogService_ListBlogServer) error {
	cursor, err := collection.Find(context.Background(), primitive.D{{}})
	if err != nil {
		return status.Errorf(codes.Internal, fmt.Sprintf("Fail to list blog.%v", err))
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		data := &blogItem{}
		err := cursor.Decode(data)
		if err != nil {
			return status.Errorf(codes.Internal, fmt.Sprintf("Fail to parse blog.%v", err))
		}
		stream.Send(&blogpb.ListResponse{
			Blog: &blogpb.Blog{
				Id:      data.ID.Hex(),
				Author:  data.Author,
				Title:   data.Title,
				Content: data.Content,
			},
		})
	}
	if err := cursor.Err(); err != nil {
		return status.Errorf(codes.Internal, fmt.Sprintf("unknown cursor error.%v", err))
	}

	return nil
}
func main() {
	// to show error file and line
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	fmt.Println("hello go server.")

	println("Connecting to mongo db.")
	// connect to mongo db
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://jinyongnan:jinyongnan@cluster0-xk5om.gcp.mongodb.net/?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatalf("Cant get mongo client.%v\n", err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Cant connect to mongodb.%v\n", err)
		return
	}
	collection = client.Database("grpc-go").Collection("blog")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatal("Fail to listen.", err)
	}
	s := grpc.NewServer()
	blogpb.RegisterBlogServiceServer(s, &server{})
	// register grpc reflection for cli
	reflection.Register(s)
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
	println("Closing listener...")
	lis.Close()
	println("Closing mongodb connection")
	client.Disconnect(ctx)
	println("Server stopped.")

}
