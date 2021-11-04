package main

import (
	"log"
	"fmt"
	"os"
	"encoding/json"
	"io/ioutil"
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

func agregar_eliminado(archivo string,info string){
	b := []byte(info+"\n")
    err := ioutil.WriteFile(archivo, b, 0644)
    if err != nil {
        log.Fatal(err)
    }
}
func main() {
	//var monto_acumulado int
	//monto := 0
	f, err := os.Create("./registro.txt")
    check(err)
    defer f.Close()
	n3, err := f.WriteString("Hola mundo\n")
    check(err)
    f.Sync()
    n3, err = f.WriteString("Hace sueño\n")
    check(err)
	_ = n3
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
		    fmt.Println(dat)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}