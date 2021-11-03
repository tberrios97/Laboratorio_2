package main

import (
	"log"
	"context"
	"time"
	pb "example.com/go-comm-grpc/comm"
	"google.golang.org/grpc"
)

const (
	address = "localhost:9000"
)

func main() {
 
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewCommClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := c.FunTest(ctx, &pb.RequestTest{Body: "Hello From Client!"})
	if err != nil {
		log.Fatalf("Error when calling SayHello: %v", err)
	}
	log.Printf("Response from server: %v", response.Body)

}