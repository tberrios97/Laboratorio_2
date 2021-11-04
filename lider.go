package main

import (
  "fmt"
  "time"
  //"strings"
  "math/rand"
  "log"
  "strconv"
  "github.com/streadway/amqp"
)

func random(min, max int) int {
  //49152-65535
  return rand.Intn(max-min) + min
}
/*
func juego_etapa_1() bool{

}

func juego_etapa_2() bool{

}

func juego_etapa_3() bool{

}
*/
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

  response := registrar_jugada_nameNode(6 ,1 ,4, "localhost")
  
  var input int
  //var pozo int
  //var jugando bool

  fmt.Println("[*] Bienvenido a SquidGame.\n[*] ¿Deseas unirte?.\n[*] (1) Si.\t(2) No.")
  fmt.Scan(&input)

  //Enviar petición al Lider si se desea jugar. En otro caso, cerrar programa.
  if input == 1 {
    //Enviar petición
    fmt.Println("[*] Iniciando contacto con el Líder.")
    //Establecer conexión al juego y recibir cuando este listo para jugar

    //...

    //Comenzar primera etapa
    /*
    jugando = juego_etapa_1()

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
    */

  } else{
    informar_jugador_eliminado(input)
  }

  //Cerrar programa
  fmt.Println("[*] Finalizando programa de SquidGame.")
  return
}
