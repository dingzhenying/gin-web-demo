package main

import (
	pb "gin-web-demo/server/grpc/helloworld"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"os"
)

//常量参数
const (
	address     = "localhost:50051"
	defaultName = "test name:world"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}

	res, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("####### get server Greeting response: %s", res.Message)
}
