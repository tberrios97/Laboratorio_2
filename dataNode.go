package main

import (
	"log"
	"net"
	"bufio"
	"strings"
  "os"
  "strconv"
	"context"
	pb "example.com/go-comm-grpc/comm"
	"google.golang.org/grpc"
)

const (
	port = ":9300"
)

type CommServer struct {
	pb.UnimplementedCommServer
}

func existeArchivo(archivo string) bool {
    info, err := os.Stat(archivo)
    if os.IsNotExist(err) {
        return false
    }
    return !info.IsDir()
}

func obtenerJugadasTexto(n_jugador int, n_ronda int) [4]int {
	var jugadas [4]int

	nombre_archivo := "jugador_"+strconv.Itoa(n_jugador)+"__ronda_"+strconv.Itoa(n_ronda)+".txt"
	log.Printf("nombre del archivo: %v", nombre_archivo)

	archivo, err := os.Open(nombre_archivo)
  if err != nil {
    log.Fatal(err)
  }
  defer archivo.Close()

  scanner := bufio.NewScanner(archivo)

	posicion := 0
  for scanner.Scan() {
    linea := strings.Replace(scanner.Text(), "\n","",-1)
		jugadas[posicion], _ = strconv.Atoi(linea)
		log.Printf("archivo jugadas en posicion %d es %v", posicion, jugadas[posicion])
		posicion++
  }

	log.Printf("arreglo jugadas %v", jugadas)
	return jugadas
}

func archivoJugada(n_jugador int, n_ronda int, jugada int){
  nombre_archivo := "jugador_"+strconv.Itoa(n_jugador)+"__ronda_"+strconv.Itoa(n_ronda)+".txt"
  if (existeArchivo(nombre_archivo)) {
    archivo, err := os.OpenFile(nombre_archivo, os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
      log.Fatalf("fallo la apertura del archivo: %s", err)
    }
    defer archivo.Close()
    linea := strconv.Itoa(jugada)
    _, err = archivo.WriteString(linea+"\n")
    if err != nil {
      // log.Printf("aqui 1")
      log.Fatalf("fallo escritura en archivo: %s", err)
    }
  } else {
    archivo, err := os.Create(nombre_archivo)
    if err != nil {
      // log.Printf("aqui 2")
      log.Fatalf("fallo escritura en archivo: %s", err)
    }
    defer archivo.Close()
    linea := strconv.Itoa(jugada)
    _, err = archivo.WriteString(linea+"\n")
    if err != nil {
      // log.Printf("aqui 3")
      log.Fatalf("fallo escritura en archivo: %s", err)
    }
  }

  return
}

func (s *CommServer) ObtenerJugada(in *pb.RequestOJ, stream pb.Comm_ObtenerJugadaServer) error{
	var jugadas [4]int

	log.Printf("Numero jugador: %d", in.GetNJugador())
  log.Printf("Numero ronda: %d", in.GetNRonda())

	jugadas = obtenerJugadasTexto(int(in.GetNJugador()), int(in.GetNRonda()))

	for i := range jugadas {
		if (jugadas[i] == 0){
			continue
		}
		if err := stream.Send(&pb.ResponseOJ{Jugadas: int32(jugadas[i])}); err != nil {
				return err
			}
	}

	return nil
}

func (s *CommServer) RegistrarJugadaDN(ctx context.Context, in *pb.RequestRJDN) (*pb.ResponseRJDN, error){
  log.Printf("Numero jugador: %d", in.GetNJugador())
  log.Printf("Numero ronda: %d", in.GetNRonda())
  log.Printf("Jugada: %d", in.GetJugada())

  archivoJugada(int(in.GetNJugador()), int(in.GetNRonda()), int(in.GetJugada()))

	return &pb.ResponseRJDN{Body: "Datos recibidos por el Data Node!"}, nil
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
