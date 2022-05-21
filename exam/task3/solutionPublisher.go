package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	printErrorAndExit(err, "Connect to RabbitMQ is Failed!")
	defer conn.Close()

	ch, err := conn.Channel()
	printErrorAndExit(err, "Opening of channel is Failed")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"p3fanoutExchange",
		"fanout",
		false,
		true,
		false,
		false,
		nil,
	)
	printErrorAndExit(err, "Invalid operation to declare an exchange!")
	for i := 0; i <= 10; i++ {
		publishMsg(ch, "p3fanoutExchange", "", i)
	}
	publishMsg(ch, "p3fanoutExchange", "", "END")
}

func printErrorAndExit(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, ":", err)
	}
}

func publishMsg(c amqp.Channel, ex string, key string, msg string) {
	body := msg
	err := (c).Publish(
		ex,
		key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	printErrorAndExit(err, "Error to publish a message")
	fmt.Println("Sent: ", body)
}
