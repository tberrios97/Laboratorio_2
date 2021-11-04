package main

import (
	"log"
	"os"
	"strconv"
	"encoding/json"
	amqp "github.com/streadway/amqp"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
func archivoExiste(ruta string) bool {
	if _, err := os.Stat(ruta); os.IsNotExist(err) {
		return false
	}
	return true
}
func agregar_eliminado(nombre_archivo string, jugador int, ronda int,monto_acumulado int){
	if (archivoExiste(nombre_archivo)){
		archivo, err := os.OpenFile(nombre_archivo, os.O_WRONLY|os.O_APPEND, 0644)
		check(err)
		defer archivo.Close()
		linea:="Jugador_"+strconv.Itoa(jugador)+" Ronda_"+strconv.Itoa(ronda)+" "+strconv.Itoa(monto_acumulado)+"\n"
		_,err = archivo.WriteString(linea)
		check(err)
	} else {
		archivo, err := os.Create(nombre_archivo)
		check(err)
		defer archivo.Close()
		defer archivo.Close()
		linea:="Jugador_"+strconv.Itoa(jugador)+" Ronda_"+strconv.Itoa(ronda)+" "+strconv.Itoa(monto_acumulado)+"\n"
		_,err = archivo.WriteString(linea)
		check(err)
	}
	return
}
func main() {
	monto_acumulado := 0
	nombre_archivo:="registro_eliminados.txt"
	if (archivoExiste(nombre_archivo)){
		err := os.Remove(nombre_archivo)
		check(err)
	}

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

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	var dat map[string]int

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			if err := json.Unmarshal(d.Body, &dat); err != nil {
		        panic(err)
		    }
		    monto_acumulado+=100000000
    		agregar_eliminado(nombre_archivo,dat["Jugador"],dat["Ronda"],monto_acumulado)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}