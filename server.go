package main

import (
	"log"
	"net"
	"context"
	pb "example.com/go-comm-grpc/comm"
	"google.golang.org/grpc"
)

const (
	port = ":9000"
)

type CommServer struct {
	pb.UnimplementedCommServer
}

func (s *CommServer) FunTest(ctx context.Context, in *pb.RequestTest) (*pb.ResponseTest, error){
	log.Printf("Request from client: %v", in.GetBody())
	return &pb.ResponseTest{Body: "Hello From Server!"}, nil
}

func main() {

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterCommServer(s, &CommServer{})

	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}

}