package main

import (
  "fmt"
  "log"
  "time"
  "context"
  "strconv"
  //"strings"
  "math/rand"
  "google.golang.org/grpc"
  pb "example.com/go-comm-grpc/comm"
)

const (
  address = "localhost:9000"
)

func printSeparador(){
  fmt.Println("[*][*][*][*][*][*][*][*][*][*][*][*][*][*][*][*][*][*][*][*][*][*][*][*][*][*]")
}

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

func menu_prints(esEtapa bool, etapa int, cliente pb.CommClient, ctx context.Context){
  var opcion int
  if esEtapa {
    for {
      fmt.Println("[*] ¿Qué acción desea realizar?")
      fmt.Println("[*] (1) Pedir monto acumulado.\t (2) Pedir acceso a la Etapa", etapa)
      fmt.Print("[*] Respuesta: ")
      fmt.Scan(&opcion)
      if opcion == 1 {
        fmt.Println("[*] Solicitando el monto total acumulado...")
        //Solicitar monto acumulado

        respuesta, err := cliente.PedirMonto(ctx, &pb.RequestPedirMonto{Body: 1})
        if err != nil {
          log.Fatalf("Error en la conexión con el servidor: %v", err)
        }

        var monto int32 = respuesta.GetMonto()
        fmt.Println("[*] Monto acumulado", monto, "KRW")
      } else {
        break
      }
    }
  } else {
    for {
      fmt.Println("[*] ¿Qué acción desea realizar?")
      fmt.Println("[*] (1) Pedir monto acumulado.\t (2) Continuar con la siguiente ronda.")
      fmt.Print("[*] Respuesta: ")
      fmt.Scan(&opcion)
      if opcion == 1 {
        fmt.Printf("[*] Solicitando el monto total acumulado...")
        //Solicitar monto acumulado

        respuesta, err := cliente.PedirMonto(ctx, &pb.RequestPedirMonto{Body: 1})
        if err != nil {
          log.Fatalf("Error en la conexión con el servidor: %v", err)
        }

        var monto int32 = respuesta.GetMonto()
        fmt.Printf("[*] Monto acumulado", monto, "KRW")
      } else {
        break
      }
    }
  }

  fmt.Println("\n")
  printSeparador()
  printSeparador()
  printSeparador()
  fmt.Println("\n")

  return
}

func juegoEtapa1(cliente pb.CommClient, ctx context.Context, numeroJugador int32) (bool, bool){
  var jugada int32
  var estado bool
  var ganador bool
  var terminoJuego bool
  var ronda int32 = 1

  menu_prints(true, 1, cliente, ctx)

  fmt.Println("[*] Comenzando primera etapa. Próximo juego:")
  fmt.Println("[*] || Juego Luz Roja, Luz Verde ||")

  //Esperar a la respuesta del servidor para iniciar la etapa
  fmt.Println("[*] Esperando inicio de la etapa por parte del servidor...")
  _, err := cliente.InicioEtapa(ctx, &pb.RequestEtapa{Etapa: 1, NumeroJugador: numeroJugador})
  if err != nil {
    log.Fatalf("Error en la conexión con el servidor: %v", err)
  }

  fmt.Println("[*] Para jugar debes elegir un número entre 1 y 10.")
  for ronda = 1; ronda <= 4; ronda ++{

    //Comprobar que el jugador no haya ganada ya la etapa
    fmt.Println("")
    printSeparador()
    fmt.Println("")
    fmt.Print("[*] Ronda ", ronda, ".\n[*] Realice su jugada: ")
    if !ganador {
      //Lectura de la jugada en cada ronda, hasta un máximo de 4
      fmt.Scan(&jugada)
      fmt.Println("[*] Esperando jugadas de los demás jugadores...")

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
        return false, false
      }
    } else {
      fmt.Println("[*] Esperando jugadas de los demás jugadores...")
    }
    
    //Espera del inicio de la siguiente ronda
    if ronda == 4 {
      respuestaRonda, error := cliente.TerminoRonda(ctx, &pb.RequestRonda{Etapa: 1, Ronda: ronda, RondaFinal: true, TerminoJuego: false})
      if error != nil {
        log.Fatalf("Error en la conexión con el servidor: %v", error)
      }
      terminoJuego = respuestaRonda.GetTerminoJuego()
    } else {
      respuestaRonda, error := cliente.TerminoRonda(ctx, &pb.RequestRonda{Etapa: 1, Ronda: ronda, RondaFinal: false, TerminoJuego: false})
      if error != nil {
        log.Fatalf("Error en la conexión con el servidor: %v", error)
      }
      terminoJuego = respuestaRonda.GetTerminoJuego()
    }
    
    //Comprobar si el juego ya se ha acabado por solo quedar un único jugador
    if terminoJuego {
      ganador = true
      //fmt.Println("[*] Juego Terminado, has ganado")
      return ganador, terminoJuego
    }
  }
  
  return ganador, terminoJuego
}

func juegoEtapa2(cliente pb.CommClient, ctx context.Context, numeroJugador int32) (bool, bool){
  var equipo int32
  var jugada int32
  var estado bool
  var terminoJuego bool
  jugando := true

  menu_prints(true, 2, cliente, ctx)

  fmt.Println("[*] Comenzado segunda etapa. Próximo juego")
  fmt.Println("[*] || Juego Tirar la Cuerda ||")

  //Realizar petición para entrar en la etapa
  //Mostrar en que equipo se encuentra o si fue eliminado al azar

  //Esperando respuesta del servidor para iniciar la etapa
  fmt.Println("[*] Esperando inicio de la etapa por parte del servidor...")
  respuesta, err := cliente.InicioEtapa(ctx, &pb.RequestEtapa{Etapa: 2, NumeroJugador: numeroJugador})
  if err != nil {
    log.Fatalf("Error en la conexión con el servidor: %v", err)
  }
  equipo = respuesta.GetBody()
  if equipo == -1{
    jugando = false
    fmt.Println("[*] Lo sentimos, has sido eliminado al azar.\n[*] Gracias por jugar.") 
    return jugando, false
  }

  fmt.Println("[*] Perteneces al equipo:", equipo)
  fmt.Print("[*] Elija un número entre 1 y 4.\n[*] Realice su jugada: ")
  fmt.Scan(&jugada)
  fmt.Println("[*] Esperando jugadas de los demás jugadores...")
  //Enviar al Líder la jugada y recibir respuesta
  respuestaEtapa, err := cliente.JugadaSegundaEtapa(ctx, &pb.RequestSegundaEtapa{Jugada: jugada, Jugador: numeroJugador})
  if err != nil {
    log.Fatalf("Error en la conexión con el servidor: %v", err)
  }

  //Obteniendo el estado del jugador, si ha sido eliminado o no
  estado = respuestaEtapa.GetEstado()
  if estado {
    jugando = true
  }else {
    jugando = false
    return jugando, false
  }

  //Aviso de termino de ronda y etapa
  respuestaRonda, err := cliente.TerminoRonda(ctx, &pb.RequestRonda{Etapa: 2, Ronda: 1, RondaFinal: true, TerminoJuego: false})
  if err != nil {
    log.Fatalf("Error en la conexión con el servidor: %v", err)
  }
  terminoJuego = respuestaRonda.GetTerminoJuego()

  return jugando, terminoJuego
}

func juegoEtapa3(cliente pb.CommClient, ctx context.Context, numeroJugador int32) bool{
  var jugada int32
  var contrincante int32
  var estado bool
  jugando := true

  menu_prints(true, 3, cliente, ctx)

  fmt.Println("[*] Comenzado última etapa.")
  fmt.Println("[*] || Juego Todo o Nada ||")
  //Realizar petición para entrar en la etapa.
  //Mostrar contrincante a vencer o si fue eliminado al azar.

  //Esperando respuesta del servidor para iniciar la etapa
  fmt.Println("[*] Esperando inicio de la etapa por parte del servidor...")
  respuesta, err := cliente.InicioEtapa(ctx, &pb.RequestEtapa{Etapa: 3, NumeroJugador: numeroJugador})
  if err != nil {
    log.Fatalf("Error en la conexión con el servidor: %v", err)
  }
  contrincante = respuesta.GetBody()
  if contrincante == -1 {
    jugando = false
    fmt.Println("[*] Lo sentimos, has sido eliminado al azar.\n[*] Gracias por jugar.") 
    return jugando
  }

  fmt.Println("[*] Tu contrincante es el Jugador", contrincante)
  fmt.Print("[*] Elija un número entre el 1 y 10.\n[*] Realice su jugada: ")
  fmt.Scan(&jugada)
  fmt.Println("[*] Esperando jugadas de los demás jugadores...")
  //Enviar al Líder la jugada y recibir respuesta
  respuestaEtapa, err := cliente.JugadaTerceraEtapa(ctx, &pb.RequestTerceraEtapa{Jugada: jugada, Jugador: numeroJugador})
  if err != nil {
    log.Fatalf("Error en la conexión con el servidor: %v", err)
  }

  //Obteniendo el estado del jugador, si ha sido eliminado o no
  estado = respuestaEtapa.GetEstado()
  if estado {
    jugando = true
  }else {
    jugando = false
    return jugando
  }

  //Aviso de termino de ronda, etapa y juego
  _, err = cliente.TerminoRonda(ctx, &pb.RequestRonda{Etapa: 3, Ronda: 1, RondaFinal: true, TerminoJuego: true})
  if err != nil {
    log.Fatalf("Error en la conexión con el servidor: %v", err)
  }

  return jugando
}

func main(){
  var input int
  var jugando bool
  var terminoJuego bool
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

  fmt.Print("[*] Bienvenido a SquidGame.\n[*] ¿Deseas unirte?.\n[*] (1) Si.\t(2) No.\n[*] Respuesta: ")
  fmt.Scan(&input)

  //Enviar petición al Lider si se desea jugar. En otro caso, cerrar programa.
  if input == 1 {

    //Enviar petición para unirse al juego si hay cupos disponibles
    fmt.Println("[*] Iniciando contacto con el Líder...")
    respuesta, err := cliente.UnirseJuegoCalamar(ctx, &pb.RequestUnirse{Body: 1})
    if err != nil {
      log.Fatalf("Error when calling SayHello: %v", err)
    }
    fmt.Println("[*] Ingreso realizado correctamente.")
    fmt.Println("[*] Número de jugador:", respuesta.NumeroJugador)
    //Registrar número del jugador que se ha elegido
    numeroJugador = respuesta.GetNumeroJugador()

    //Si el número de jugador es 0 significa que el juego ya esta en transcurso
    if numeroJugador == 0{
      fmt.Println("[*] No hay cupos disponibles para unirse al juego.\n[*] Finalizando programa de SquidGame.")
      return
    }
    
    /*
    *
    * Sección Etapa 1
    *
    */
    
    fmt.Println("\n")
    printSeparador()
    printSeparador()
    printSeparador()
    fmt.Println("\n")

    //Comienzo de la primera etapa
    jugando, terminoJuego = juegoEtapa1(cliente, ctx, numeroJugador)

    //Verificar si el jugador sigue jugando
    if !jugando{
      fmt.Println("[*] Finalizando programa de SquidGame.")
      return
    }

    //Verificar si ya solo queda un único jugador
    if terminoJuego {
      respuesta, err := cliente.PedirMonto(ctx, &pb.RequestPedirMonto{Body: 1})
      if err != nil {
        log.Fatalf("Error en la conexión con el servidor: %v", err)
      }

      var monto int32 = respuesta.GetMonto()
      fmt.Println("[*] Felicitaciones jugador " + strconv.Itoa(int(numeroJugador)) + ", has gando el Juego del Calamar.\n[*] Has ganado", monto, "KRW")
      fmt.Println("[*] Finalizando programa de SquidGame.")
      return
    }
    
    /*
    *
    * Sección Etapa 2
    *
    */
    
    fmt.Println("\n")
    printSeparador()
    printSeparador()
    printSeparador()
    fmt.Println("\n")

    //Comienzo de la segunda etapa
    jugando, terminoJuego = juegoEtapa2(cliente, ctx, numeroJugador)

    //Verificar si el jugador sigue jugando
    if !jugando{
      fmt.Println("[*] Finalizando programa de SquidGame.")
      return
    }

    //Verificar si ya solo queda un único jugador
    if terminoJuego {
      respuesta, err := cliente.PedirMonto(ctx, &pb.RequestPedirMonto{Body: 1})
      if err != nil {
        log.Fatalf("Error en la conexión con el servidor: %v", err)
      }

      var monto int32 = respuesta.GetMonto()
      fmt.Println("[*] Felicitaciones jugador " + strconv.Itoa(int(numeroJugador)) + ", has gando el Juego del Calamar.\n[*] Has ganado", monto, "KRW")
      fmt.Println("[*] Finalizando programa de SquidGame.")
      return
    }
    
    /*
    *
    * Sección Etapa 3
    *
    */

    fmt.Println("\n")
    printSeparador()
    printSeparador()
    printSeparador()
    fmt.Println("\n")

    //Comienzo de la tercera etapa
    jugando = juegoEtapa3(cliente, ctx, numeroJugador)

    if !jugando{
      fmt.Println("[*] Has sido eliminado\n[*] Finalizando programa de SquidGame.")
      return
    }
    respuestaMonto, err := cliente.PedirMonto(ctx, &pb.RequestPedirMonto{Body: 1})
    if err != nil {
      log.Fatalf("Error en la conexión con el servidor: %v", err)
    }

    var monto int32 = respuestaMonto.GetMonto()
    fmt.Println("[*] Feliciticaciones, has sido uno de los ganadores.\n[*] Has ganado", monto, "KRW")

  }
  //Cerrar programa
  fmt.Println("[*] Finalizando programa de SquidGame.")
  return
}