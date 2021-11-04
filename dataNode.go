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

func (s *CommServer) RequestRJDN(ctx context.Context, in *pb.RequestRJDN) (*pb.ResponseRJDN, error){
	log.Printf("Request desde Name Node: %v", in.GetBody())
	return &pb.ResponseTest{Body: "Datos recibidos por el Data Node!"}, nil
}

func main() {

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("fallo el escuchar: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterCommServer(s, &CommServer{})

	log.Printf("Data Node escuchando en %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("fallo el servir: %s", err)
	}

}
