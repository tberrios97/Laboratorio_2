package main

import (
  "net"
  "log"
  "time"
  "strconv"
  "context"
  //"strings"
  "math/rand"
  "google.golang.org/grpc"
  "github.com/streadway/amqp"
  pb "example.com/go-comm-grpc/comm"
  
)

var jugadoresActivos int32 = 0
var juegoActivo bool = false

const (
  port = ":9000"
)

type CommServer struct {
  pb.UnimplementedCommServer
}

func (s *CommServer) UnirseJuegoCalamar(stream pb.Comm_UnirseJuegoCalamarClient) error {
  log.Printf("[*] Petición de unirse al juego del calamar recibida.")
  for {
    requestUnirse, err := stream.Recv()
    if juegoActivo {
      log.Printf("[*] Juego en transcurso. Petición denegada")
      return stream.SendAndClose(&pb.ResponseUnirse{NumeroJugador: 0})
    } else {
      log.Printf("[*] Espacio disponible. Petición aceptada")
      jugadoresActivos = jugadoresActivos + 1
      if jugadoresActivos == 2{
        juegoActivo = true
      }
      log.Printf("[*] Jugadores activos: %d/16", jugadoresActivos)
      //return stream.SendAndClose(&pb.ResponseUnirse{NumeroJugador: jugadoresActivos})
    }
  }
}

func random(min, max int) int {
  //49152-65535
  return rand.Intn(max-min) + min
}

func failOnError(err error, msg string) {
  if err != nil {
    log.Fatalf("%s: %s", msg, err)
  }
}

func registrar_jugada_nameNode(id_jugador int, num_ronda int, jugada int, direccion_nameNode string){
  conn, err := grpc.Dial(direccion_nameNode, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("No se pudo lograr conexión: %v", err)
	}
	defer conn.Close()

	c := pb.NewCommClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := c.RegistrarJugadaJugador(ctx, &pb.RequestRJJ{N_jugador: id_jugador, N_ronda: num_ronda, Jugada: jugada})
	if err != nil {
		log.Fatalf("Error al hacer request a servidor: %v", err)
	}
	log.Printf("Response desde Data Node: %v", response.Body)
}

func informar_jugador_eliminado(id_jugador int){
  //Se crea ña conexión y se abre el canal para el paso de mensajes:
  conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
  failOnError(err, "Failed to connect to RabbitMQ")
  defer conn.Close()

  ch, err := conn.Channel()
  failOnError(err, "Failed to open a channel")
  defer ch.Close()

  q, err := ch.QueueDeclare(
    "cola RabbitMQ", // name
    false,   // durable
    false,   // delete when unused
    false,   // exclusive
    false,   // no-wait
    nil,     // arguments
  )
  failOnError(err, "Failed to declare a queue")
  rand.Seed(time.Now().UnixNano())
  jugador := strconv.Itoa(id_jugador)
  body := "Jugador Eliminado "+jugador
  err = ch.Publish(
    "",     // exchange
    q.Name, // routing key
    false,  // mandatory
    false,  // immediate
    amqp.Publishing{
      ContentType: "text/plain",
      Body:        []byte(body),
    })
  failOnError(err, "Failed to publish a message")
  log.Printf(" [*] Mensaje enviado al Pozo: %s", body)
  return &body
}

func main(){
  //var input int
  response := registrar_jugada_nameNode(6 ,1 ,4, "localhost")

  lis, err := net.Listen("tcp", port)
  if err != nil {
    log.Fatalf("failed to listen: %v", err)
  }

  s := grpc.NewServer()

  pb.RegisterCommServer(s, &CommServer{})

  log.Printf("Servidor escuchando en %v", lis.Addr())

  if err := s.Serve(lis); err != nil {
    log.Fatalf("failed to serve: %s", err)
  }


  return
}
