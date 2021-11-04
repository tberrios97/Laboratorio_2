package main

import (
  "net"
  "encoding/json"
  "log"
  "fmt"
  "time"
  "context"
  "math/rand"
  "google.golang.org/grpc"
  "github.com/streadway/amqp"
  pb "example.com/go-comm-grpc/comm"
)

const (
  port = ":9000"
  capacidadJugadores = 2
)

//
var jugadoresActivos int32 = 0
var jugadoresListos int32 = 0

//
var juegoActivo bool = false
var comienzoEtapa bool = false
var comienzoRonda bool = false
var jugadasRecolectadas bool = false
var bloqueo bool = false

//
var contadorJugadaJugador [capacidadJugadores]int32
var jugadaLider int32

type CommServer struct {
  pb.UnimplementedCommServer
}

func random(min, max int) int32 {
  return int32(rand.Intn(max-min+1) + min)
}

func resetContadorJugadores() (){
  for i := 0; i < capacidadJugadores; i++ {
    if contadorJugadaJugador[i] != -1 {
      contadorJugadaJugador[i] = 0
    }
  }
  return
}

func (s *CommServer) UnirseJuegoCalamar(ctx context.Context, in *pb.RequestUnirse) (*pb.ResponseUnirse, error){
  log.Printf("[*] Petición de unirse al juego del calamar recibida.")
  if juegoActivo || jugadoresActivos >= capacidadJugadores{
    log.Printf("[*] Juego en transcurso. Petición denegada")
    return &pb.ResponseUnirse{NumeroJugador: 0}, nil
  } else {
    log.Printf("[*] Espacio disponible. Petición aceptada")
    jugadoresActivos ++
    if jugadoresActivos == capacidadJugadores {
      juegoActivo = true
    }
    log.Printf("[*] Jugadores activos: %d/%d", jugadoresActivos, capacidadJugadores)
    return &pb.ResponseUnirse{NumeroJugador: jugadoresActivos}, nil
  }
}

func (s *CommServer) InicioEtapa(ctx context.Context, in *pb.RequestEtapa) (*pb.ResponseEtapa, error){
  var input int
  bloqueo = false
  comienzoEtapa = false
  jugadoresListos ++

  if jugadoresListos == jugadoresActivos && juegoActivo{
    log.Printf("[*] ¿Listos para comenzar?\n[*] (1) Si\t(2)No")
    fmt.Scan(&input)
    if input == 1{
      log.Printf("[*] Si")
    }else {
      log.Printf("[*] No")
    }    
    comienzoEtapa = true
    jugadoresListos = 0
    resetContadorJugadores()
  }

  for {
    if comienzoEtapa {
      if jugadoresActivos == 1{
        return &pb.ResponseEtapa{Body: 1, TerminoJuego: true}, nil
      } else {
        return &pb.ResponseEtapa{Body: 1, TerminoJuego: false}, nil
      }
    }
  }
}

func (s *CommServer) InicioRonda(ctx context.Context, in *pb.RequestRonda) (*pb.ReponseRonda, error){
  var input int
  bloqueo = false
  comienzoRonda = false
  jugadoresListos ++

  if jugadoresListos == jugadoresActivos && juegoActivo{
    log.Printf("[*] ¿Listos para comenzar?\n[*] (1) Si\t(2)No")
    fmt.Scan(&input)
    if input == 1{
      log.Printf("[*] Si")
    }else {
      log.Printf("[*] No")
    }    
    comienzoRonda = true
    jugadoresListos = 0
  }

  for {
    if comienzoRonda {
      if jugadoresActivos == 1{
        return &pb.ReponseRonda{Body: 1, TerminoJuego: true}, nil
      } else {
        return &pb.ReponseRonda{Body: 1, TerminoJuego: false}, nil
      }
    }
  }
}

func (s *CommServer) JugadaPrimeraEtapa(ctx context.Context, in *pb.RequestPrimeraEtapa) (*pb.ResponsePrimeraEtapa, error) {
  bloqueo = false
  jugadasRecolectadas = false
  jugadoresListos ++

  jugada := in.GetJugada()
  ronda := in.GetRonda()
  jugador := in.GetJugador()

  //Suma de la jugada actual del jugador
  contadorJugadaJugador[jugador - 1] = contadorJugadaJugador[jugador - 1] + jugada

  //Condición cuando el último jugador haya realizado su jugada
  if jugadoresListos == jugadoresActivos && juegoActivo{
    jugadoresListos = 0
    jugadasRecolectadas = true
    jugadaLider = random(6, 10)
    log.Printf("Jugada del lider: %d", jugadaLider)
  }

  for {
    //Cuando todos los jugadores hayan realizado sus jugadas
    if jugadasRecolectadas {
      log.Printf("Ronda %d Jugada del jugador %d: %d", ronda, jugador, jugada)

      //Si un jugador ingresa un número mayor o igual al del lider.
      //Se debe informar que un jugador ha sido eliminado.
      if jugada >= jugadaLider {
        jugadoresActivos --
        contadorJugadaJugador[jugador - 1] = -1
        return &pb.ResponsePrimeraEtapa{Estado: false, Ganador: false}, nil
      }

      //Si un jugador logra llegar a los 21, gana
      if contadorJugadaJugador[jugador - 1] >= 21{
        return &pb.ResponsePrimeraEtapa{Estado: true, Ganador: true}, nil
      } 

      //Ronda final
      if ronda == 4 {
        log.Printf("Conteo del jugador %d: %d", jugador, contadorJugadaJugador[jugador - 1])
        //Si el jugador llego a la meta.
        //Caso contrario, se debe informar que se ha eliminado.
        if contadorJugadaJugador[jugador - 1] >= 21 {
          return &pb.ResponsePrimeraEtapa{Estado: true, Ganador: true}, nil
        } else {
          jugadoresActivos --
          contadorJugadaJugador[jugador - 1] = -1
          return &pb.ResponsePrimeraEtapa{Estado: false, Ganador: false}, nil
        }
      } else {
        //Caso en que debe seguir jugando pero no ha ganado la etapa aún
        return &pb.ResponsePrimeraEtapa{Estado: true, Ganador: false}, nil
      }
      
    }
  } 
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
  //Seteo de semilla aleatoria para que funcione mejor el random
  rand.Seed(time.Now().UnixNano())
  /*
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
  */
  //var input int
  //response := registrar_jugada_nameNode(6 ,1 ,4, "localhost:9100")
  //log.Printf("Response : %v", response)

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