package main

import (
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/streadway/amqp"
)

func main() {
	var wg sync.WaitGroup

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	printErrorAndExit(err, "Error to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	printErrorAndExit(err, "Error to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare("p3fanoutExchange", "fanout", false, true, false, false, nil)
	printErrorAndExit(err, "Error to declare an exchange")

	q, err := ch.QueueDeclare("", false, true, true, false, nil)
	printErrorAndExit(err, "Error to declare a queue")
	err = ch.QueueBind(q.Name, "", "p3fanoutExchange", false, nil)
	printErrorAndExit(err, "Error to bind a queue")

	msgs, err := ch.Consume(q.Name, "Inc", false, false, false, false, nil)
	printErrorAndExit(err, "Error to register a consumer")

	go func() {
		for i := 0; i < 3; i++ {
			wg.Add(1)
			for d := range msgs {
				bodyString := string(d.Body)
				fmt.Println(bodyString)
				if bodyString == "END" {
					ch.Close()
					return
				}
				res, err := strconv.Atoi(string(bodyString))
				if err != nil {
					printErrorAndExit(err, "Error to convert")
				}
				fmt.Println("Received:", bodyString, "Result:", res+1)
				d.Ack(false)
			}
			defer wg.Done()
		}
	}()

	wg.Wait()

	fmt.Println("Waiting for messages")
	forever := make(chan bool)
	<-forever
	fmt.Println("Finish")
}

func printErrorAndExit(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, ":", err)
	}
}
