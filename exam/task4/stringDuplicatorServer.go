package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/streadway/amqp"
)

func main() {
	conn1, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Error to connect to RabbitMQ")
	defer conn1.Close()

	conn2, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Error to connect to RabbitMQ")
	defer conn2.Close()

	cho, err := conn1.Channel()
	failOnError(err, "Error to open a channel")
	defer cho.Close()
	chi, err := conn2.Channel()
	failOnError(err, "Error to open a channel")
	defer chi.Close()

	err = cho.ExchangeDeclare("p4Exchange", "direct", false, true, false, false, nil)
	failOnError(err, "Error to declare an exchange")

	q, err := chi.QueueDeclare("stringQueue", false, true, false, false, nil)
	failOnError(err, "Error to declare a queue")

	err = chi.QueueBind(q.Name, "string", "p4Exchange", false, nil)
	failOnError(err, "Error to bind a queue")

	msgs, err := chi.Consume(q.Name, "", false, false, false, false, nil)
	failOnError(err, "Error to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			_, err := strconv.Atoi(string(d.Body))
			if err == nil {
				continue
			}
			result := string(d.Body) + string(d.Body)
			err = cho.Publish(
				"p4Exchange", d.ReplyTo, false, false,
				amqp.Publishing{
					ContentType:   "text/plain",
					CorrelationId: d.CorrelationId,
					Body:          []byte(result),
				})
			failOnError(err, "Error to publish a message")
			fmt.Println("Received job:", string(d.Body), "Published response:", result)
			d.Ack(false)
		}
	}()
	fmt.Println("Waiting for jobs")
	<-forever
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
