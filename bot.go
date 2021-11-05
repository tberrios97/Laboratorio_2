package main

import (
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

func juegoEtapa1(cliente pb.CommClient, ctx context.Context, numeroJugador int32) (bool, bool, int32, int32){
  var jugada int32
  var estado bool
  var ganador bool
  var terminoJuego bool
  var ronda int32 = 1

  //Esperar a la respuesta del servidor para iniciar la etapa
  _, err := cliente.InicioEtapa(ctx, &pb.RequestEtapa{Etapa: 1, NumeroJugador: numeroJugador})
  if err != nil {
    log.Fatalf("Error en la conexión con el servidor: %v", err)
  }
  for ronda = 1; ronda <= 4; ronda ++{

    //Comprobar que el jugador no haya ganada ya la etapa
    if !ganador {
      //Lectura de la jugada en cada ronda, hasta un máximo de 4
      jugada = int32(random(5, 6))

      //Enviar al Líder la jugada y recibir respuesta
      respuesta, err := cliente.JugadaPrimeraEtapa(ctx, &pb.RequestPrimeraEtapa{Jugada: jugada, Ronda:ronda, Jugador:numeroJugador})
      if err != nil {
        log.Fatalf("Error: %v", err)
      }
      estado = respuesta.GetEstado()
      ganador = respuesta.GetGanador()

      //Comprobar si sigue ha sido eliminado el jugador
      if(!estado){
        return false, false, 0, 0
      }
    }
    
    //Espera del inicio de la siguiente ronda
    if ronda == 4 {
      respuestaRonda, error := cliente.TerminoRonda(ctx, &pb.RequestRonda{Etapa: 1, Ronda: ronda, RondaFinal: true, TerminoJuego: false})
      if error != nil {
        log.Fatalf("Error en la conexión con el servidor: %v", error)
      }
      terminoJuego = respuestaRonda.GetTerminoJuego()

      //Comprobar si el juego ya se ha acabado por solo quedar un único jugador
      if terminoJuego {
        ganador = true
        return ganador, terminoJuego, respuestaRonda.GetMontoAcumulado(), respuestaRonda.GetJugadores()
      }

    } else {
      respuestaRonda, error := cliente.TerminoRonda(ctx, &pb.RequestRonda{Etapa: 1, Ronda: ronda, RondaFinal: false, TerminoJuego: false})
      if error != nil {
        log.Fatalf("Error en la conexión con el servidor: %v", error)
      }
      terminoJuego = respuestaRonda.GetTerminoJuego()

      //Comprobar si el juego ya se ha acabado por solo quedar un único jugador
      if terminoJuego {
        ganador = true
        return ganador, terminoJuego, respuestaRonda.GetMontoAcumulado(), respuestaRonda.GetJugadores()
      }
    }
    
  }
  
  return ganador, terminoJuego, 0, 0
}

func juegoEtapa2(cliente pb.CommClient, ctx context.Context, numeroJugador int32) (bool, bool, int32, int32){
  var equipo int32
  var jugada int32
  var estado bool
  var terminoJuego bool
  jugando := true

  //Realizar petición para entrar en la etapa
  //Mostrar en que equipo se encuentra o si fue eliminado al azar

  //Esperando respuesta del servidor para iniciar la etapa
  respuesta, err := cliente.InicioEtapa(ctx, &pb.RequestEtapa{Etapa: 2, NumeroJugador: numeroJugador})
  if err != nil {
    log.Fatalf("Error en la conexión con el servidor: %v", err)
  }
  equipo = respuesta.GetBody()
  if equipo == -1{
    jugando = false
    return jugando, false, 0, 0
  }

  jugada = int32(random(1, 4))
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
    return jugando, false, 0, 0
  }

  //Aviso de termino de ronda y etapa
  respuestaRonda, err := cliente.TerminoRonda(ctx, &pb.RequestRonda{Etapa: 2, Ronda: 1, RondaFinal: true, TerminoJuego: false})
  if err != nil {
    log.Fatalf("Error en la conexión con el servidor: %v", err)
  }
  terminoJuego = respuestaRonda.GetTerminoJuego()

  return jugando, terminoJuego, respuestaRonda.GetMontoAcumulado(), respuestaRonda.GetJugadores()
}

func juegoEtapa3(cliente pb.CommClient, ctx context.Context, numeroJugador int32) (bool, int32, int32){
  var jugada int32
  var contrincante int32
  var estado bool
  jugando := true

  //Realizar petición para entrar en la etapa.
  //Mostrar contrincante a vencer o si fue eliminado al azar.

  //Esperando respuesta del servidor para iniciar la etapa
  respuesta, err := cliente.InicioEtapa(ctx, &pb.RequestEtapa{Etapa: 3, NumeroJugador: numeroJugador})
  if err != nil {
    log.Fatalf("Error en la conexión con el servidor: %v", err)
  }
  contrincante = respuesta.GetBody()
  if contrincante == -1 {
    jugando = false
    return jugando, 0, 0
  }

  jugada = int32(random(1, 10))
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
    return jugando, 0, 0
  }

  //Aviso de termino de ronda, etapa y juego
  respuestaRonda, err := cliente.TerminoRonda(ctx, &pb.RequestRonda{Etapa: 3, Ronda: 1, RondaFinal: true, TerminoJuego: true})
  if err != nil {
    log.Fatalf("Error en la conexión con el servidor: %v", err)
  }

  return jugando, respuestaRonda.GetMontoAcumulado(), respuestaRonda.GetJugadores()
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

  input = 1

  //Enviar petición al Lider si se desea jugar. En otro caso, cerrar programa.
  if input == 1 {

    //Enviar petición para unirse al juego si hay cupos disponibles
    respuesta, err := cliente.UnirseJuegoCalamar(ctx, &pb.RequestUnirse{Body: 1})
    if err != nil {
      log.Fatalf("Error when calling SayHello: %v", err)
    }
    //Registrar número del jugador que se ha elegido
    numeroJugador = respuesta.GetNumeroJugador()

    //Si el número de jugador es 0 significa que el juego ya esta en transcurso
    if numeroJugador == 0{
      return
    }
    
    /*
    *
    * Sección Etapa 1
    *
    */
    
    //Comienzo de la primera etapa
    jugando, terminoJuego, _, _ = juegoEtapa1(cliente, ctx, numeroJugador)

    //Verificar si el jugador sigue jugando
    if !jugando{
      return
    }

    //Verificar si ya solo queda un único jugador
    if terminoJuego {
      return
    }
    
    /*
    *
    * Sección Etapa 2
    *
    */

    //Comienzo de la segunda etapa
    jugando, terminoJuego, _, _ = juegoEtapa2(cliente, ctx, numeroJugador)

    //Verificar si el jugador sigue jugando
    if !jugando{
      return
    }

    //Verificar si ya solo queda un único jugador
    if terminoJuego {
      return
    }
    
    /*
    *
    * Sección Etapa 3
    *
    */

    //Comienzo de la tercera etapa
    jugando, _, _ = juegoEtapa3(cliente, ctx, numeroJugador)

    if !jugando{
      return
    }

  }
  //Cerrar programa
  return
}