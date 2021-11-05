package main

import (
  "net"
  "encoding/json"
  "log"
  "fmt"
  "time"
  "context"
  "strconv"
  "math/rand"
  "google.golang.org/grpc"
  "github.com/streadway/amqp"
  pb "example.com/go-comm-grpc/comm"
)

const (
  port = ":9000"
  capacidadJugadores = 2
  address = "localhost:3000"  //REVISAR
)

//
var jugadoresActivos int32 = 0
var jugadoresListos int32 = 0

//
var juegoActivo bool = false
var comienzoEtapa bool = false
var comienzoRonda bool = false
var jugadasRecolectadas bool = false

//
var contadorJugadaJugador [capacidadJugadores]int32
var conteoEquipo1 int32
var conteoEquipo2 int32
var equipoGanador int32
var jugadaLider int32

type CommServer struct {
  pb.UnimplementedCommServer
}

func random(min, max int) int32 {
  return int32(rand.Intn(max-min+1) + min)
}

func elegirJugadorRandomActivo() int32{
  var posicion int32
  for {
    posicion = random(1, capacidadJugadores)
    if contadorJugadaJugador[posicion - 1] != -1 {
      return posicion
    }
  }
}

func elegirEquipos() {
  var cantidad int = int(jugadoresActivos/2)
  var count int = 0
  for i := 0; i < capacidadJugadores; i++ {
    if contadorJugadaJugador[i] != -1 {
      if count < cantidad {
        contadorJugadaJugador[i] = 1
      } else {
        contadorJugadaJugador[i] = 2
      }
      count ++
    }
  }
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
  var etapa int32 = in.GetEtapa()
  var numeroJugador = in.GetNumeroJugador()
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
    
    //reset de contadores de los jugadores activos a 0
    resetContadorJugadores()

    //Etapa 2. Verificar paridad de jugadores y asignar equipos.
    if etapa == 2 {
      //Verificar si hay jugadores impares
      if jugadoresActivos % 2 == 1 {
        //Eliminar jugador al azar
        jugadorRandom := elegirJugadorRandomActivo()
        contadorJugadaJugador[jugadorRandom - 1] = -1
        jugadoresActivos --
      }
      //Crear equipos
      elegirEquipos()
      conteoEquipo1 = 0
      conteoEquipo2 = 0
    }

    //Setear variables para ejecutar la siguiente etapa
    comienzoEtapa = true
    jugadoresListos = 0
  }

  for {

    if comienzoEtapa {
      //Manejo de Etapas
      if etapa == 1 { 
        //Etapa 1. No hay manipulación extra
        if jugadoresActivos == 1{
          return &pb.ResponseEtapa{Body: 1, TerminoJuego: true}, nil
        } else {
          return &pb.ResponseEtapa{Body: 1, TerminoJuego: false}, nil
        }
      } else if etapa == 2 { 
        //Etapa 2. Verificar paridad de jugadores y asignar equipos.
        var terminoJuego bool
        var body int32 = contadorJugadaJugador[numeroJugador - 1]
        if jugadoresActivos == 1 {
          terminoJuego = true
        } else {
          terminoJuego = false
        }
        return &pb.ResponseEtapa{Body: body, TerminoJuego: terminoJuego}, nil
      }
    }
  }
}

func (s *CommServer) TerminoRonda(ctx context.Context, in *pb.RequestRonda) (*pb.ReponseRonda, error){
  var input int
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

      //Si es la última ronda, hacer reset
      if in.GetRondaFinal() {
        resetContadorJugadores()
      }

      if jugadoresActivos == 1{
        return &pb.ReponseRonda{Body: 1, TerminoJuego: true}, nil
      } else {
        return &pb.ReponseRonda{Body: 1, TerminoJuego: false}, nil
      }
    }
  }
}

func (s *CommServer) JugadaPrimeraEtapa(ctx context.Context, in *pb.RequestPrimeraEtapa) (*pb.ResponsePrimeraEtapa, error) {
  jugadasRecolectadas = false

  jugada := in.GetJugada()
  ronda := in.GetRonda()
  jugador := in.GetJugador()

  //Suma de la jugada actual del jugador
  contadorJugadaJugador[jugador - 1] = contadorJugadaJugador[jugador - 1] + jugada

  //Condición cuando el último jugador haya realizado su jugada
  jugadoresListos ++
  if jugadoresListos == jugadoresActivos && juegoActivo{
    jugadoresListos = 0
    jugadaLider = random(6, 10)
    jugadasRecolectadas = true
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

func (s *CommServer) JugadaSegundaEtapa(ctx context.Context, in *pb.RequestSegundaEtapa) (*pb.ResponseSegundaEtapa, error) {
  jugadasRecolectadas = false
  
  jugada := in.GetJugada()
  jugador := in.GetJugador()

  //Contar jugada segun el equipo del jugador
  if contadorJugadaJugador[jugador - 1] == 1 {
    conteoEquipo1 = conteoEquipo1 + jugada
  } else {
    conteoEquipo2 = conteoEquipo2 + jugada
  }

  //Condición cuando el último jugador haya realizado su jugada
  jugadoresListos ++
  if jugadoresListos == jugadoresActivos && juegoActivo{
    jugadoresListos = 0
    jugadaLider = random(1, 4)
    
    //Verificar que equipo gano, en el caso
    if conteoEquipo1%2 == jugadaLider%2 && conteoEquipo2%2 == jugadaLider%2 {
      //Caso en que ambos equipos hayan ganado
      equipoGanador = -1
    } else if conteoEquipo1%2 == jugadaLider%2 {
      //Caso en que el equipo 1 haya ganado
      equipoGanador = 1
      jugadoresActivos = jugadoresActivos/2
    } else if conteoEquipo2%2 == jugadaLider%2{
      //Caso en que el equipo 2 haya ganado
      equipoGanador = 2
      jugadoresActivos = jugadoresActivos/2
    } else {
      //Caso en que ningún equipo haya ganado
      //Eliminar un equipo al azar
      var equipoEliminado int32 = random(1, 2)
      if equipoEliminado == 1 {
        equipoGanador = 2
      } else {
        equipoGanador = 1
      }
      jugadoresActivos = jugadoresActivos/2
    }

    jugadasRecolectadas = true
    log.Printf("Jugada del lider: %d", jugadaLider)
  }

  for {
    //Cuando todos los jugadores hayan realizado sus jugadas
    if jugadasRecolectadas {
      //Caso en que ambos equipos ganaron
      if equipoGanador == -1 {
        return &pb.ResponseSegundaEtapa{Estado: true}, nil
      } else {
        //Verificar si el jugador sobrevivio o no
        if contadorJugadaJugador[jugador - 1] == equipoGanador {
          return &pb.ResponseSegundaEtapa{Estado: true}, nil
        } else {
          return &pb.ResponseSegundaEtapa{Estado: false}, nil
        }
      }
    }
  }

}

func failOnError(err error, msg string) {
  if err != nil {
    log.Fatalf("%s: %s", msg, err)
  }
}

func SolicitarMonto() int32{
  coneccion, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
  if err != nil {
    log.Fatalf("did not connect: %v", err)
  }
  defer coneccion.Close()
  cliente := pb.NewCommClient(coneccion)
  ctx, cancel := context.WithTimeout(context.Background(), time.Second)
  defer cancel()
  response, err := cliente.SolicitarMonto(ctx, &pb.RequestMonto{MontoAcumulado: int32(1)}) //Revisar
  if err != nil {
      log.Fatalf("Error en la conexión con el servidor: %v", err)
    }
  return response.GetMontoAcumulado()
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

func menu_prints(opcion int, adicional int)(){ //adicional es para agregar el numero del jugador eliminado | el numero de ganadores | numero de etapa
  if opcion == 0{ //Inicio
    log.Printf("[*] Bienvenidos al juego del Calamar\n")
    log.Printf("[*] Esperando jugadores...\n")

  }else if opcion == 1{ //Jugador eliminado
    log.Printf("[*] Jugador"+strconv.Itoa(adicional)+"eliminado\n")
    log.Printf("[*] %d jugadores restantes\n", jugadoresActivos)

  }else if opcion == 2{ //Pedir monto acumulado
    log.Printf("[*] Solicitando el monto total acumulado...\n")
    SolicitarMonto()

  }else if opcion == 3{ //Ganadores
    if adicional > 1{
      fmt.Println("[*] El juego ha finalizado, tenemos %d ganadores!\n",adicional)
      fmt.Println("[*] Felicidades! han ganado el Juego del Calamar\n") 
    }else{
      fmt.Println("[*] El juego ha finalizado, tenemos un ganador!\n")
      fmt.Println("[*] Felicidades! has ganado el Juego del Calamar\n")
    }
  }else if opcion ==4{ //Inicio etapa
    if adicional == 1{
      fmt.Println("[*] Comienza la Etapa 1: Luz Roja, Luz Verde\n") 
    }else if adicional ==2{
      fmt.Println("[*] Comienza la Etapa 2: Tirar la cuerda\n") 
    }else{
      fmt.Println("[*] Comienza la Etapa 3: Todo o nada\n") 
    }
  }else if opcion ==5{ //Pedir jugadas
    log.Printf("[*] Solicitando jugadas del jugador"+strconv.Itoa(adicional)+"...\n")
  }
  return
}

func main(){
  //Seteo de semilla aleatoria para que funcione mejor el random
  rand.Seed(time.Now().UnixNano())

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
