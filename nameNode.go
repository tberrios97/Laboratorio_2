package main

import (
	"log"
	"net"
  "bufio"
  "strings"
  "time"
	"io"
  "os"
  "strconv"
  "math/rand"
	"context"
	pb "example.com/go-comm-grpc/comm"
	"google.golang.org/grpc"
)

func random(min, max int) int {
  //49152-65535
  return rand.Intn(max-min) + min
}

const (
  port = ":9100"
)

var direcciones_dataNode = []string{"dist57", "dist58", "dist59"}

type CommServer struct {
	pb.UnimplementedCommServer
}
func check(e error) {
    if e != nil {
        panic(e)
    }
}

func existeArchivo(archivo string) bool {
    info, err := os.Stat(archivo)
    if os.IsNotExist(err) {
        return false
    }
    return !info.IsDir()
}

func buscarEnArchivo(n_jugador int, n_ronda int) string{
  var retorno string
  archivo, err := os.Open("registro_jugadas.txt")
  if err != nil {
    log.Fatal(err)
  }
  defer archivo.Close()

  scanner := bufio.NewScanner(archivo)
  for scanner.Scan() {
    linea := strings.Split(scanner.Text()," ")
    if (linea[0][len(linea[0])-1:] == strconv.Itoa(n_jugador) && linea[1][len(linea[1])-1:] == strconv.Itoa(n_ronda)) {
      retorno = strings.Replace(linea[2], "\n","",-1)
    }
  }

  return retorno
}

func archivoJugada(n_jugador int, n_ronda int, direccion_dataNode string){
  //poner en el README esto
  nombre_archivo := "registro_jugadas.txt"
  if (existeArchivo(nombre_archivo)) {
    archivo, err := os.OpenFile(nombre_archivo, os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
      log.Fatalf("fallo la apertura del archivo: %s", err)
    }
    defer archivo.Close()
    linea := "Jugador_"+strconv.Itoa(n_jugador)+" Ronda_"+strconv.Itoa(n_ronda)+" "+direccion_dataNode+""
    _, err = archivo.WriteString(linea+"\n")
    if err != nil {
      log.Fatalf("fallo escritura en archivo: %s", err)
    }
  } else {
    archivo, err := os.Create(nombre_archivo)
    if err != nil {
      log.Fatalf("fallo escritura en archivo: %s", err)
    }
    defer archivo.Close()
    linea := "Jugador_"+strconv.Itoa(n_jugador)+" Ronda_"+strconv.Itoa(n_ronda)+" "+direccion_dataNode+""
    _, err = archivo.WriteString(linea+"\n")
    if err != nil {
      log.Fatalf("fallo escritura en archivo: %s", err)
    }
  }
  return
}

func resetNameNode(){
	nombre_archivo:="registro_jugadas.txt"
	if (existeArchivo(nombre_archivo)){
		err := os.Remove(nombre_archivo)
		check(err)
	}

	return
}

func resetDataNode(address string){
  coneccion, err := grpc.Dial(address, grpc.WithInsecure())
  if err != nil {
    log.Fatalf("did not connect: %v", err)
  }
  defer coneccion.Close()
  cliente := pb.NewCommClient(coneccion)
  ctx, cancel := context.WithTimeout(context.Background(), time.Second)
  defer cancel()
  _, err = cliente.ReiniciarPartida(ctx, &pb.RequestTest{Body: "hola jorge :D"})
  if err != nil {
      log.Fatalf("Error en la conexión con el servidor: %v", err)
    }
  return
}

func (s *CommServer) ReiniciarPartida(ctx context.Context, in *pb.RequestTest) (*pb.ResponseTest, error){
	nombre_archivo:="registro_jugadas.txt"
	if (existeArchivo(nombre_archivo)){
		err := os.Remove(nombre_archivo)
		check(err)
	}
	resetDataNode("dist57:9300")
	resetDataNode("dist58:9300")
	resetDataNode("dist59:9300")
	return &pb.ResponseTest{Body: "hola jorge :D"}, nil
}

func (s *CommServer) BuscarJugada(in *pb.RequestBJ, stream pb.Comm_BuscarJugadaServer) error{
	var direccion string

	log.Printf("Numero jugador: %d", in.GetNJugador())
  log.Printf("Numero ronda: %d", in.GetNRonda())

	direccion = buscarEnArchivo(int(in.GetNJugador()), int(in.GetNRonda()))

	conn, err := grpc.Dial(direccion+":9300", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("no se conecto: %v", err)
	}
	defer conn.Close()

	c := pb.NewCommClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(),10*time.Second)
	defer cancel()

	streaming, err := c.ObtenerJugada(ctx, &pb.RequestOJ{NJugador: in.GetNJugador(), NRonda: in.GetNRonda()})
	if err != nil {
		log.Fatalf("%v.ListFeatures(_) = _, %v", c, err)
	}
	for {
		jugada, err := streaming.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.ListFeatures(_) = _, %v", c, err)
		}
		log.Printf("Jugada recibida: %v", jugada.GetJugadas())
		if err := stream.Send(&pb.ResponseBJ{Jugadas: jugada.GetJugadas()}); err != nil {
			return err
		}
	}

	return nil
}

func (s *CommServer) RegistrarJugadaJugador(ctx context.Context, in *pb.RequestRJJ) (*pb.ResponseRJJ, error){
  var direccion string
  var pos_aleatorio int

	log.Printf("Numero jugador: %d", in.GetNJugador())
  log.Printf("Numero ronda: %d", in.GetNRonda())
  log.Printf("Jugada: %d", in.GetJugada())

  if(int(in.GetNRonda()) == 1 && existeArchivo("registro_jugadas.txt")){
    direccion = buscarEnArchivo(int(in.GetNJugador()), int(in.GetNRonda()))
  }else{
    pos_aleatorio = random(0,2)
    direccion = direcciones_dataNode[pos_aleatorio]
    archivoJugada(int(in.GetNJugador()), int(in.GetNRonda()), direccion)
  }

  conn, err := grpc.Dial(direccion+":9300", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("no se conecto: %v", err)
	}
	defer conn.Close()

	c := pb.NewCommClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := c.RegistrarJugadaDN(ctx, &pb.RequestRJDN{NJugador: in.GetNJugador(), NRonda: in.GetNRonda(), Jugada: in.GetJugada()})
	if err != nil {
		log.Fatalf("Error when calling SayHello: %v", err)
	}
	log.Printf("Respuesta del servidor: %v", response.Body)

	return &pb.ResponseRJJ{Body: "Jugada recibida por NameNode"}, nil
}

func main() {
	resetNameNode()

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
