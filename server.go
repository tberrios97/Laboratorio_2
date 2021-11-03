package main

import (
	"fmt"
	"log"
	"net"
	"context"
	pb "github.com/tberrios97/Laboratorio_2/comm"
	"google.golang.org/grpc"
)

const (
	port = ":9000"
)

type CommServer struct {
	pb.UnimplementedCommServer
}

fun (s *CommServer) FunTest(ctx context.Context, in *pb.RequestTest) (*pb.ResponseTest, error){
	log.Printf("Request from client: %v", in.GetBody())
	return &pb.RequestTest{body: "Hello From Server!"}, nil
}

func main() {

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterCommServer(s, &CommServer{})

	log.Printf("server listening at %v", lis.Addr())

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}

}