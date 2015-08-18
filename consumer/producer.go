package producer

import (
	"fmt"
	"time"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		// log.Fatalf("%s: %s", msg, err)
		// panic(fmt.Sprintf("%s: %s", msg, err))
		fmt.Println("Error:", err)
	}
}

func Initialize() {
	// conn, err := amqp.Dial("amqp://test:letmein@10.5.212.156:5672/")
	conn, err := amqp.Dial("amqp://fahoobagats:Fytg*pv2c,iCUT@rabbitmq-statsnode-3a5n:5672/")
	//failOnError(err, "Failed to connect to RabbitMQ")

	if err != nil {
		//log.Fatalf("%s: %s", msg, err)
		//panic(fmt.Sprintf("%s: %s", msg, err))
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Rabbit MQ Producer Connected!!")
		defer conn.Close()

		ch, err := conn.Channel()
		failOnError(err, "Failed to open a channel")
		defer ch.Close()

		q, err := ch.QueueDeclare(
			"hello", // name
			false,   // durable
			false,   // delete when usused
			false,   // exclusive
			false,   // no-wait
			nil,     // arguments
		)
		failOnError(err, "Failed to declare a queue")

		fmt.Sprintf("\n", time.Now(), ": starting to send\n")

		i := 0
		for i = 0; i < 100000; i++ {
			body := fmt.Sprintf("hello %d", i)
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
		}

		fmt.Sprintf("\n", time.Now(), ": done sending\n")

		body := fmt.Sprintf("done")
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
	}
}
