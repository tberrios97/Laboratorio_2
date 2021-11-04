package main

import (
  "fmt"
  "log"
  "time"
  "context"
  //"strings"
  "math/rand"
  "google.golang.org/grpc"
  pb "example.com/go-comm-grpc/comm"
)

const (
  address = "localhost:9000"
)

func random(min, max int) int {
  //49152-65535
  return rand.Intn(max-min) + min
}

func intermedio(etapa string){
  var pozo int
  var input int
  fmt.Println("[*] Felicitaciones! has sobrevivido a la", etapa, "etapa.")
  for{
    fmt.Println("[*] ¿Qué deseas hacer?.\n[*] (1) Seguir a la siguiente etapa.\t(2) Ver pozo.")
    fmt.Scan(&input)

    if input == 2 {
      //Pedir ver el pozo
      pozo = 100000
      fmt.Println("[*] El pozo es de:", pozo, "KRW")
    }else {
      break
    }
  }
}

func juegoEtapa1(cliente pb.CommClient, ctx context.Context, numeroJugador int32) bool{
  var jugada int32
  var estado bool
  var ganador bool
  var terminoJuego bool
  var ronda int32 = 1
  fmt.Println("[*] Comenzando primera etapa.")
  fmt.Println("[*] Juego Luz Roja, Luz Verde.")
  fmt.Println("[*] Para jugar debes elegir un número entre 1 y 10.")
  for ronda <= 4{
    //Lectura de la jugada en cada ronda, hasta un máximo de 4
    fmt.Print("[*] Ronda", ronda, ".\n[*] Realice su jugada: ")
    fmt.Scan(&jugada)
    fmt.Println("[*] Jugada:", jugada)
    
    //Enviar al Líder la jugada y recibir respuesta
    respuesta, err := cliente.JugadaPrimeraEtapa(ctx, &pb.RequestPrimeraEtapa{Jugada: jugada, Ronda:ronda, Jugador:numeroJugador})
    if err != nil {
      log.Fatalf("Error: %v", err)
    }
    estado = respuesta.GetEstado()
    ganador = respuesta.GetGanador()

    //Comprobar si sigue ha sido eliminado el jugador
    if(!estado){
      fmt.Println("[*] Has sido eliminado.\n[*] Gracias por jugar.")
      break
    }
    //comprobar si ya ha ganado la etapa
    if ganador {
      break
    }

    //Espera del inicio de la siguiente ronda
    respuesta, err = cliente.InicioRonda(ctx, &pb.RequestRonda{Body: 1})
    if err != nil {
      log.Fatalf("Error en la conexión con el servidor: %v", err)
    }
    terminoJuego = respuesta.GetTerminoJuego()

    ronda = ronda + 1
    fmt.Println("[*]\n[*]\n[*]")
  }
  
  return ganador
}

func juego_etapa_2() bool{
  var equipo int
  var jugada int
  var respuesta int
  jugando := true
  fmt.Println("[*] Comenzado segunda etapa.")
  fmt.Println("[*] Juego Tirar la Cuerda.")

  //Realizar petición para entrar en la etapa
  //Mostrar en que equipo se encuentra o si fue eliminado al azar
  equipo = 1
  if equipo == 3{
    jugando = false
    fmt.Println("[*] Lo sentimos, has sido eliminado al azar.\n[*] Gracias por jugar.") 
    return jugando
  }

  fmt.Println("[*] Perteneces al equipo:", equipo)
  fmt.Println("[*] Elija un número entre 1 y 4.")
  fmt.Scan(&jugada)

  //Enviar al Líder la jugada y recibir respuesta
  respuesta = 1
  if respuesta == 1 {
    jugando = true
  }else {
    jugando = false
    fmt.Println("[*] Has sido eliminado.\n[*] Gracias por jugar.")
  }

  return jugando
}

func juego_etapa_3() bool{
  var jugada int
  var respuesta int
  var contrincante int
  jugando := true
  fmt.Println("[*] Comenzado última etapa.")
  fmt.Println("[*] Juego Todo o Nada.")
  //Realizar petición para entrar en la etapa.
  //Mostrar contrincante a vencer o si fue eliminado al azar.
  contrincante = 12

  if contrincante == 17 {
    jugando = false
    fmt.Println("[*] Lo sentimos, has sido eliminado al azar.\n[*] Gracias por jugar.") 
    return jugando
  }

  fmt.Println("[*] Tu contrincante es el Jugador", contrincante)
  fmt.Println("[*] Elija un número entre el 1 y 10.")
  fmt.Scan(&jugada)

  //Enviar al Líder la jugada y recibir respuesta
  respuesta = 1
  if respuesta == 1 {
    jugando = true
  }else {
    jugando = false
    fmt.Println("[*] Has sido eliminado.\n[*] Gracias por jugar.")
  }

  return jugando
}

func main(){
  var pozo int
  var input int
  var jugando bool
  var numeroJugador int32
  
  //Definicion de la conexión con el servidor
  conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
  if err != nil {
    log.Fatalf("did not connect: %v", err)
  }
  defer conn.Close()
  cliente := pb.NewCommClient(conn)
  ctx, cancel := context.WithTimeout(context.Background(), 600 * time.Second)
  defer cancel()

  rand.Seed(time.Now().UnixNano())

  fmt.Println("[*] Bienvenido a SquidGame.\n[*] ¿Deseas unirte?.\n[*] (1) Si.\t(2) No.")
  fmt.Scan(&input)

  //Enviar petición al Lider si se desea jugar. En otro caso, cerrar programa.
  if input == 1 {
    //Enviar petición
    fmt.Println("[*] Iniciando contacto con el Líder.")
    respuesta, err := cliente.UnirseJuegoCalamar(ctx, &pb.RequestUnirse{Body: 1})
    if err != nil {
      log.Fatalf("Error when calling SayHello: %v", err)
    }
    log.Printf("Número de jugar: %v", respuesta.NumeroJugador)

    numeroJugador = respuesta.GetNumeroJugador()

    if numeroJugador == 0{
      fmt.Println("[*] Finalizando programa de SquidGame.")
      return
    }

    //Esperar a la respuesta del servidor para iniciar la etapa
    fmt.Println("[*] Esperando a los demás jugadores...")
    _, err = cliente.InicioEtapa(ctx, &pb.RequestEtapa{Body: 1})
    if err != nil {
      log.Fatalf("Error en la conexión con el servidor: %v", err)
    }

    //Comenzar primera etapa
    jugando = juegoEtapa1(cliente, ctx, numeroJugador)

    if !jugando{
      fmt.Println("[*] Finalizando programa de SquidGame.")
      return
    }

    intermedio("primera")

    jugando = juego_etapa_2()

    if !jugando{
      fmt.Println("[*] Finalizando programa de SquidGame.")
      return
    }

    intermedio("segunda")

    jugando = juego_etapa_3()

    if !jugando{
      fmt.Println("[*] Finalizando programa de SquidGame.")
      return
    }

    pozo = 100
    fmt.Println("[*] Feliciticaciones, has sido uno de los ganadores.\n[*] Has ganado", pozo, "KRW")

  }

  //Cerrar programa
  fmt.Println("[*] Finalizando programa de SquidGame.")
  return
}