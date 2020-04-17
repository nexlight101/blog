package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/nexlight101/blog/blog_server/blogpb"
	"google.golang.org/grpc"
)

func main() {

	fmt.Println("Blog client")
	// Create a connection
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	// Close connection when done
	defer conn.Close()

	// Create a new client
	c := blogpb.NewBlogServiceClient(conn)
	//Create Blog
	fmt.Println("Creating the blog")
	blog := &blogpb.Blog{
		AuthorId: "Hendrik",
		Title:    "My First Blog",
		Content:  "Content of the first blog",
	}
	createBlogRes, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}
	fmt.Printf("Blog has been created: %v", createBlogRes)
	// extract the id of the created blog for testing the readblog
	blogID := createBlogRes.GetBlog().GetId()

	// read Blog
	fmt.Println("Reading the blog")

	_, err2 := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{
		BlogId: "Hello I'm testing for crap",
	})
	if err2 != nil {
		fmt.Printf("Error happened while reading: %v\n", err2)
	}
	readBlogReq := &blogpb.ReadBlogRequest{BlogId: blogID}
	readBlogRes, readBlogErr := c.ReadBlog(context.Background(), readBlogReq)
	if err2 != nil {
		fmt.Printf("Error happened while reading: %v\n", readBlogErr)
	}

	fmt.Printf("Blog was read: %v\n", readBlogRes)
	// Update Blog
	newBlog := &blogpb.Blog{
		Id:       blogID,
		AuthorId: "Changed Author",
		Title:    "My First Blog (edited)",
		Content:  "Content of the first blog, with some awesome additions",
	}
	updateRes, upddateErr := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{Blog: newBlog})
	if upddateErr != nil {
		fmt.Printf("Error happened while updating: %v\n", upddateErr)
	}
	fmt.Printf("Blog was updated: %v\n", updateRes)

	// Delete blog
	fmt.Println("Deleting the blog")

	DeleteBlogReq := &blogpb.DeleteBlogRequest{BlogId: blogID}
	DeleteBlogRes, DeleteBlogErr := c.DeleteBlog(context.Background(), DeleteBlogReq)
	if DeleteBlogErr != nil {
		fmt.Printf("Error happened while Deleteing: %v\n", DeleteBlogErr)
	}

	fmt.Printf("Blog was Delete: %v\n", DeleteBlogRes)

	// list blogs
	fmt.Println("Starting to list blogs...")

	resStream, streamErr := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})
	if err != nil {
		log.Fatalf("error while calling ListBlog RPC: %v\n ", streamErr)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// we've reached the end of the stream
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v\n", err)
		}
		fmt.Println(msg.GetBlog())
	}
}
