package main

import (
	"log"
	"net"
	"context"
	pb "example.com/go-comm-grpc/comm"
	"google.golang.org/grpc"
)

func random(min, max int) int {
  //49152-65535
  return rand.Intn(max-min) + min
}

const (
  port = ":9000"
  var direcciones_dataNode [3]string = [3]{"localhost", "localhost", "localhost"}
)

type CommServer struct {
	pb.UnimplementedCommServer
}

func (s *CommServer) RegistrarJugadaJugador(ctx context.Context, in *pb.RequestRJJ) (*pb.ResponseRJJ, error){
	log.Printf("Numero jugador: %d", in.GetN_jugador())
  log.Printf("Numero ronda: %d", in.GetN_ronda())
  log.Printf("Jugada: %d", in.GetJugada())

  pos_aleatorio := random(0,2)

  conn, err := grpc.Dial(direcciones_dataNode[pos_aleatorio], grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("no se conecto: %v", err)
	}
	defer conn.Close()

	c := pb.NewCommClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := c.RequestRJDN(ctx, &pb.RequestRJDN{in.GetN_jugador(), in.GetN_ronda(), in.GetJugada()})
	if err != nil {
		log.Fatalf("Error when calling SayHello: %v", err)
	}
	log.Printf("Respuesta del servidor: %v", response.Body)

	return &pb.ResponseRJJ{Body: "Jugada recibida por NameNode"}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("fallo el escuchar: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterCommServer(s, &CommServer{})

	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
