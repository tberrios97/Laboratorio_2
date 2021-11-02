package main

import (
	"log"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"github.com/tberrios97/Laboratorio_2/comm"
)

func main() {

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := comm.NewCommClient(conn)

	response, err := c.funTest(context.Background(), &comm.requestTest{Body: "Hello From Client!"})
	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}
	log.Printf("Response from server: %s", response.Body)

}