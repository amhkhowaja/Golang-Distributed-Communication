package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/rs/xid"
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

	var jobCorr = make(map[string]string)
	var mu sync.Mutex

	msgs, err := chi.Consume(q.Name, "", false, false, false, false, nil)
	failOnError(err, "Error to register a consumer")

	go func() {
		for d := range msgs {
			mu.Lock()
			v, ok := jobCorr[d.CorrelationId]
			if ok {
				fmt.Println("Job:", v, "Got response:"+string(d.Body))
				delete(jobCorr, d.CorrelationId)
			} else {
				fmt.Println("Got a not related msg")
			}
			mu.Unlock()
			d.Ack(false)
		}
	}()

	jobs := []string{"aaa", "bbb", "777", "ccc", "999"}

	for _, j := range jobs {
		var corrId = randomString()
		err := cho.Publish("p4Exchange", "string", false, false,
			amqp.Publishing{
				ContentType:   "text/plain",
				CorrelationId: corrId,
				ReplyTo:       q.Name,
				Body:          []byte(j),
			})

		failOnError(err, "Error to publish")
		fmt.Println("Published job:" + j)
		mu.Lock()
		jobCorr[corrId] = j
		mu.Unlock()
	}

	forever := make(chan bool)
	<-forever
}

func randomString() string {
	guid := xid.New()
	return guid.String()
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
