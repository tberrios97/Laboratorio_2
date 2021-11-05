package main

import (
  "io"
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
  address = "dist59:3000"  //REVISAR
)

//
var jugadoresActivos int32 = 0
var jugadoresListos int32 = 0

//
var juegoActivo bool = false
var comienzoEtapa bool = false
var comienzoRonda bool = false
var jugadasRecolectadas bool = false
var montoFinal int32 = -1
var jugadoresFinales int32 = -1

//
var contadorJugadaJugador [capacidadJugadores]int32
var conteoEquipo1 int32 = 0
var conteoEquipo2 int32 = 0
var equipoGanador int32 = 0
var jugadaLider int32 = 0

type CommServer struct {
  pb.UnimplementedCommServer
}

func mostrarGanadores(){
  log.Printf("[*] Ganador(es) del Juego del Calamar:")
  for i := 0; i < capacidadJugadores; i++ {
    if contadorJugadaJugador[i] != -1 {
      log.Printf("\t[*] Jugador %d.", i + 1)
    }
  }
}

func mostrarJugadoresVivos(){
  log.Printf("[*] Lista de Jugadores vivos terminada la etapa:")
  for i := 0; i < capacidadJugadores; i++ {
    if contadorJugadaJugador[i] != -1 {
      log.Printf("\t[*] Jugador %d.", i + 1)
    }
  }
}

func resetPozo(){
  coneccion, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
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

func resetNameNode(){
  coneccion, err := grpc.Dial("dist60:9100", grpc.WithInsecure(), grpc.WithBlock())
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

func resetPartida() {
  jugadoresActivos = 0
  jugadoresListos = 0
  juegoActivo = false
  comienzoEtapa = false
  comienzoRonda = false
  jugadasRecolectadas = false
  conteoEquipo1 = 0
  conteoEquipo2 = 0
  equipoGanador = 0
  jugadaLider = 0
  for i := 0; i < capacidadJugadores; i++ {
    contadorJugadaJugador[i] = 0
  }
  resetPozo()
  resetNameNode()
}

func abs(number int32) int32 {
  if number < 0 {
    return -number
  } else {
    return number
  }
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

func elegirParejas() {
  a := -1
  b := -1
  for i := 0; i < capacidadJugadores; i++ {
    if contadorJugadaJugador[i] != -1 {
      if a == -1 {
        a = i
      } else if b == -1 {
        b = i
        contadorJugadaJugador[b] = int32(a + 1)
        contadorJugadaJugador[a] = int32(b + 1)
        a = -1
        b = -1
      }
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

type jugador_eliminado struct {
  Jugador int `json: jugador`
  Ronda int `json: ronda`
}

func failOnError(err error, msg string) {
  if err != nil {
    log.Fatalf("%s: %s", msg, err)
  }
}

func informar_jugador_eliminado(id_jugador int, ronda int){
  conn, err := amqp.Dial("amqp://guest:guest@dist59:4369/")
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
  //log.Printf("[*] Jugador eliminado informado al pozo")
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

//adicional es para agregar el numero del jugador eliminado | el numero de ganadores | numero de etapa
func menu_prints(etapa int32, ronda int32, esEtapa bool)(){
  var opcion int
  var jugadas [4]int
  if esEtapa {
    for {
      log.Printf("[*] ¿Qué acción desea realizar?")
      log.Printf("[*] (1) Pedir monto acumulado.\t (2) Continuar con la Etapa %d.", etapa)
      fmt.Scan(&opcion)
      if opcion == 1 {
        log.Printf("[*] Solicitando el monto total acumulado...")
        //Solicitar monto acumulado
        var monto int32 = SolicitarMonto()
        log.Printf("[*] Monto acumulado %d KRW", monto)
      } else {
        if etapa == 1{
          log.Printf("[*] Comienza la Etapa 1: Luz Roja, Luz Verde")
        }else if etapa ==2{
          log.Printf("[*] Comienza la Etapa 2: Tirar la cuerda")
        }else{
          log.Printf("[*] Comienza la Etapa 3: Todo o nada")
        }
        log.Printf("[*] Esperando jugadas...")
        break
      }
    }
  } else {
    for {
      log.Printf("[*] ¿Qué acción desea realizar?")
      log.Printf("[*] (1) Ver jugada.\t (2) Pedir monto acumulado.\t (3) Terminar la ronda.")
      fmt.Scan(&opcion)
      if opcion == 1 {
        log.Printf("[*] Ingrese número de jugador:")
        fmt.Scan(&opcion)
        //Pedir jugada
        log.Printf("[*] Solicitando jugadas del jugador "+strconv.Itoa(opcion)+"...")
        jugadas = buscar_jugada_nameNode(int32(opcion), etapa, "dist60:9100")
        posicion := 0
        for {
          if (posicion == 4) {
            break
          }
          if (jugadas[posicion] == -1) {
            log.Printf("[*] El jugador seleccionado no realizó jugadas en esta etapa")
            break
          }
          if (jugadas[posicion] == 0){
            break
          }
          log.Printf("[*] Las jugada número %d realizada por el Jugador %d fue: %v", posicion+1, opcion, jugadas[posicion])
          posicion++
        }
      } else if opcion == 2 {
        log.Printf("[*] Solicitando el monto total acumulado...")
        //Solicitar monto acumulado
        var monto int32 = SolicitarMonto()
        log.Printf("[*] Monto acumulado %d KRW", monto)
      } else {
        log.Printf("[*] Terminando ronda %d", ronda)
        log.Printf("[*] Esperando jugadas...")
        break
      }
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
    log.Printf("[*] Jugadores activos: %d/%d", jugadoresActivos, capacidadJugadores)
    if jugadoresActivos == capacidadJugadores {
      log.Printf("")
      log.Printf("[*][*][*][*][*][*][*][*][*][*][*][*][*][*][*][*][*][*][*][*][*][*][*][*][*][*]")
      log.Printf("")
      log.Printf("[*] Comienzo del Juego del Calamar.")
      juegoActivo = true
    }

    return &pb.ResponseUnirse{NumeroJugador: jugadoresActivos}, nil
  }
}

func (s *CommServer) InicioEtapa(ctx context.Context, in *pb.RequestEtapa) (*pb.ResponseEtapa, error){
  var etapa int32 = in.GetEtapa()
  var numeroJugador = in.GetNumeroJugador()
  comienzoEtapa = false
  jugadoresListos ++

  if jugadoresListos == jugadoresActivos && juegoActivo{

    //Consola del Líder al inicio de cada Etapa

    menu_prints(etapa, 1, true)

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
    } else if etapa == 3 { //Etapa 3. Verificar paridad de jugadores y asignar contrincante.
      //Verificar si hay jugadores impares
      if jugadoresActivos % 2 == 1 {
        //Eliminar jugador al azar
        jugadorRandom := elegirJugadorRandomActivo()
        contadorJugadaJugador[jugadorRandom - 1] = -1
        jugadoresActivos --
      }
      //Elegir parejas
      elegirParejas()
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
        //Etapa 2. Devolver equipo del jugador
        var terminoJuego bool
        var body int32 = contadorJugadaJugador[numeroJugador - 1]
        if jugadoresActivos == 1 {
          terminoJuego = true
        } else {
          terminoJuego = false
        }
        return &pb.ResponseEtapa{Body: body, TerminoJuego: terminoJuego}, nil
      } else {
        //Etapa 3. Devolver la pareja del jugador
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
  comienzoRonda = false
  jugadoresListos ++
  etapa := in.GetEtapa()
  ronda := in.GetRonda()

  if jugadoresListos == jugadoresActivos && juegoActivo{
    //Consola del Líder despues de cada ronda

    menu_prints(etapa, ronda, false)

    //Mostrar información terminado el juego o una etapa
    if in.GetTerminoJuego() || jugadoresActivos <= 1 {
      //Reset de variables para iniciar otra partida
      montoFinal = SolicitarMonto()
      jugadoresFinales = jugadoresActivos
      mostrarGanadores()
      resetPartida()
      log.Printf("[*] Partida Finalizada.")
      log.Printf("")
      log.Printf("[*][*][*][*][*][*][*][*][*][*][*][*][*][*][*][*][*][*][*][*][*][*][*][*][*][*]")
      log.Printf("")

      //Dar comienzo a la siguiente ronda
      comienzoRonda = true
      jugadoresListos = 0

      return &pb.ReponseRonda{Body: 1, TerminoJuego: true, MontoAcumulado: montoFinal, Jugadores: jugadoresFinales}, nil
    } else {
      if in.GetRondaFinal() {
        log.Printf("[*] Termino de la Etapa")
        mostrarJugadoresVivos()
        resetContadorJugadores()
      }

      //Dar comienzo a la siguiente ronda
      comienzoRonda = true
      jugadoresListos = 0
    }

  }

  for {
    if comienzoRonda {

      //Si se termino el juego
      if in.GetTerminoJuego(){
        return &pb.ReponseRonda{Body: 1, TerminoJuego: true, MontoAcumulado: montoFinal, Jugadores: jugadoresFinales}, nil
      } else {
        if jugadoresActivos == 1{
          return &pb.ReponseRonda{Body: 1, TerminoJuego: true, MontoAcumulado: montoFinal, Jugadores: jugadoresFinales}, nil
        } else {
          return &pb.ReponseRonda{Body: 1, TerminoJuego: false, MontoAcumulado: montoFinal, Jugadores: jugadoresFinales}, nil
        }
      }

    }
  }
}

func (s *CommServer) JugadaPrimeraEtapa(ctx context.Context, in *pb.RequestPrimeraEtapa) (*pb.ResponsePrimeraEtapa, error) {
  jugadasRecolectadas = false

  jugada := in.GetJugada()
  ronda := in.GetRonda()
  jugador := in.GetJugador()

  registrar_jugada_nameNode(jugador ,1 ,jugada, "dist60:9100")

  //Suma de la jugada actual del jugador
  contadorJugadaJugador[jugador - 1] = contadorJugadaJugador[jugador - 1] + jugada

  //Condición cuando el último jugador haya realizado su jugada
  jugadoresListos ++
  if jugadoresListos == jugadoresActivos && juegoActivo{
    jugadoresListos = 0
    jugadaLider = random(6, 10)
    jugadasRecolectadas = true
    log.Printf("[*] Etapa 1 Ronda %d Jugada del lider: %d ", ronda, jugadaLider)
  }

  for {
    //Cuando todos los jugadores hayan realizado sus jugadas
    if jugadasRecolectadas {
      //log.Printf("Ronda %d Jugada del jugador %d: %d", ronda, jugador, jugada)

      //Si un jugador ingresa un número mayor o igual al del lider.
      //Se debe informar que un jugador ha sido eliminado.
      if jugada >= jugadaLider {
        jugadoresActivos --
        contadorJugadaJugador[jugador - 1] = -1
        log.Printf("[*] Jugador %d eliminado.", jugador)
        //Mandar al pozo la eliminación del jugador
        informar_jugador_eliminado(int(jugador), 1)
        return &pb.ResponsePrimeraEtapa{Estado: false, Ganador: false}, nil
      }

      //Si un jugador logra llegar a los 21, gana
      if contadorJugadaJugador[jugador - 1] >= 21{
        return &pb.ResponsePrimeraEtapa{Estado: true, Ganador: true}, nil
      }

      //Ronda final
      if ronda == 4 {
        //log.Printf("Conteo del jugador %d: %d", jugador, contadorJugadaJugador[jugador - 1])
        //Si el jugador llego a la meta.
        //Caso contrario, se debe informar que se ha eliminado.
        if contadorJugadaJugador[jugador - 1] >= 21 {
          return &pb.ResponsePrimeraEtapa{Estado: true, Ganador: true}, nil
        } else {
          jugadoresActivos --
          contadorJugadaJugador[jugador - 1] = -1
          log.Printf("[*] Jugador %d eliminado.", jugador)
          //Mandar al pozo la eliminación del jugador
          informar_jugador_eliminado(int(jugador), 1)
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

  registrar_jugada_nameNode(jugador ,2 ,jugada ,"dist60:9100")

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
    log.Printf("[*] Etapa 2 Ronda 1 Jugada del lider: %d", jugadaLider)
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
          contadorJugadaJugador[jugador - 1] = -1
          log.Printf("[*] Jugador %d eliminado.", jugador)
          //Mandar al pozo la eliminación del jugador
          informar_jugador_eliminado(int(jugador), 2)
          return &pb.ResponseSegundaEtapa{Estado: false}, nil
        }
      }
    }
  }

}

func (s *CommServer) JugadaTerceraEtapa(ctx context.Context, in *pb.RequestTerceraEtapa) (*pb.ResponseTerceraEtapa, error) {
  jugadasRecolectadas = false

  jugada := in.GetJugada()
  jugador := in.GetJugador()

  registrar_jugada_nameNode(jugador ,3 ,jugada, "dist60:9100")

  //Obtener oponente y registrar jugada
  var oponente int32 = contadorJugadaJugador[jugador - 1]
  contadorJugadaJugador[jugador - 1] = jugada

  //Condición cuando el último jugador haya realizado su jugada
  jugadoresListos ++
  if jugadoresListos == jugadoresActivos && juegoActivo{
    jugadoresListos = 0
    //Realizar jugada del Lider
    jugadaLider = random(1, 10)
    jugadasRecolectadas = true
    log.Printf("[*] Jugada del lider: %d", jugadaLider)
  }

  for {
    if jugadasRecolectadas {
      //Calcular diferencia en absoluto del oponente y jugador
      var jugadaOponente int32 = contadorJugadaJugador[oponente - 1]
      var absolutoOponente int32 = abs(jugadaOponente - jugadaLider)
      var absolutoJugador int32 = abs(jugada - jugadaLider)
      //Comparar quien tiene un valor absoluto menor
      if absolutoJugador == absolutoOponente {
        //Caso en que ambos jugadores ganar
        return &pb.ResponseTerceraEtapa{Estado: true}, nil
      } else if absolutoJugador < absolutoOponente {
        //Caso en que el jugador ganar
        return &pb.ResponseTerceraEtapa{Estado: true}, nil
      } else {
        //Caso en que el oponente gana
        jugadoresActivos --
        contadorJugadaJugador[jugador - 1] = -1
        log.Printf("[*] Jugador %d eliminado.", jugador)
        //Mandar al pozo la eliminación del jugador
        informar_jugador_eliminado(int(jugador), 3)
        return &pb.ResponseTerceraEtapa{Estado: false}, nil
      }
    }
  }

}

func (s *CommServer) PedirMonto(ctx context.Context, in *pb.RequestPedirMonto) (*pb.ResponsePedirMonto, error) {
  var monto int32 = SolicitarMonto()
  return &pb.ResponsePedirMonto{Monto: monto}, nil
}

func buscar_jugada_nameNode(id_jugador int32, num_ronda int32, direccion_nameNode string) [4]int{
  var jugadas [4]int
  conn, err := grpc.Dial(direccion_nameNode, grpc.WithInsecure(), grpc.WithBlock())
  if err != nil {
      log.Fatalf("No se pudo lograr conexión: %v", err)
  }
  defer conn.Close()

  c := pb.NewCommClient(conn)

  ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
  defer cancel()

  response, err := c.BuscarJugada(ctx, &pb.RequestBJ{NJugador: id_jugador, NRonda: num_ronda})
  if err != nil {
        log.Fatalf("%v.ListFeatures(_) = _, %v", c, err)
    }

  posicion := 0
  for {
    jugada, err := response.Recv()
    if err == io.EOF || posicion == 4 {
        break
    }
    if err != nil {
        log.Fatalf("%v.ListFeatures(_) = _, %v", c, err)
    }
    //log.Printf("Jugada recibida desde dataNode: %v [%d]", jugada.GetJugadas())
    jugadas[posicion] = int(jugada.GetJugadas())
    posicion++
  }

  return jugadas
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
  return response.Body
}

func main(){
  //Seteo de semilla aleatoria para que funcione mejor el random
  rand.Seed(time.Now().UnixNano())

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
