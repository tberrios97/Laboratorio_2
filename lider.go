package main

import (
  "encoding/json"
  "log"
  "time"
  "context"
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

func (s *CommServer) UnirseJuegoCalamar(ctx context.Context, in *pb.RequestUnirse) (*pb.ResponseUnirse, error){
  log.Printf("[*] Petición de unirse al juego del calamar recibida.")
  if juegoActivo {
    log.Printf("[*] Juego en transcurso. Petición denegada")
    return &pb.ResponseUnirse{NumeroJugador: 0}, nil
  } else {
    log.Printf("[*] Espacio disponible. Petición aceptada")
    jugadoresActivos = jugadoresActivos + 1
    if jugadoresActivos == 16{
      juegoActivo = true
    }
    log.Printf("[*] Jugadores activos: %d/16", jugadoresActivos)
    return &pb.ResponseUnirse{NumeroJugador: jugadoresActivos}, nil
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

func registrar_jugada_nameNode(id_jugador int32, num_ronda int32, jugada int32, direccion_nameNode string) string{
  conn, err := grpc.Dial(direccion_nameNode, grpc.WithInsecure(), grpc.WithBlock())
  if err != nil {
    log.Fatalf("No se pudo lograr conexión: %v", err)
  }
  defer conn.Close()

  c := pb.NewCommClient(conn)

  ctx, cancel := context.WithTimeout(context.Background(), time.Second)
  defer cancel()

  response, err := c.RegistrarJugadaJugador(ctx, &pb.RequestRJJ{NJugador: id_jugador, NRonda: num_ronda, Jugada: jugada})
  if err != nil {
    log.Fatalf("Error al hacer request a servidor: %v", err)
  }
  log.Printf("Response desde Data Node: %v", response.Body)
  return response.Body
}

type jugador_eliminado struct {
  Jugador int `json: jugador`
  Ronda int `json: ronda`
}
func informar_jugador_eliminado(id_jugador int, ronda int){
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

  data := jugador_eliminado{
    Jugador: id_jugador,
    Ronda: ronda,
  }
  dataBytes,err := json.Marshal(data)
  err = ch.Publish(
    "",     // exchange
    q.Name, // routing key
    false,  // mandatory
    false,  // immediate
    amqp.Publishing{
      ContentType: "text/plain",
      Body:        dataBytes,
    })
  failOnError(err, "Failed to publish a message")
  log.Printf("[*] Jugador eliminado informado al pozo")
}
func fin_ronda(jugadores [16]int,ronda int){
  for i := 0; i < 16; i++ {
      if jugadores[i] == -1{
        informar_jugador_eliminado(i, ronda)
      }
    }
}
func main(){
  //var input int
  response := registrar_jugada_nameNode(6 ,1 ,4, "localhost:9100")
  log.Printf("Response : %v", response)
  /*
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
  */
  //Hay que decidir como los tendremos, pensaba un numero para identificarlos y -1 si muere
  var jugadores[16] int
  for i := 0; i < 16; i++ {
    jugadores[i] = i
  }
  jugadores[4] = -1 //Mato al jugador de la pos 4
  ronda := 3
  fin_ronda(jugadores,ronda) //al terminar una ronda se revisa el estado de los jugadores
  jugadores[7] = -1
  fin_ronda(jugadores,ronda)
  return
}