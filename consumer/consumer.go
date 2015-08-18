package consumer

import (
	"fmt"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		//log.Fatalf("%s: %s", msg, err)
		//panic(fmt.Sprintf("%s: %s", msg, err))
		fmt.Println("Error:", err)
	}
}

func Initialize() {
	// conn, err := amqp.Dial("amqp://test:letmein@10.5.212.156:5672/")
	conn, err := amqp.Dial("amqp://fahoobagats:Fytg*pv2c,iCUT@rabbitmq-queue-nodes-9a13e:5672/")
	//failOnError(err, "Failed to connect to RabbitMQ")

	if err != nil {
		//log.Fatalf("%s: %s", msg, err)
		//panic(fmt.Sprintf("%s: %s", msg, err))
		fmt.Println("Error:", err)
	} else {
		fmt.Println("We are connected!!")
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

		go func() {
			for d := range msgs {
				if string(d.Body) == "done" {
					fmt.Printf("Received a message: %s", d.Body)
				}
			}
		}()

		fmt.Printf(" [*] Waiting for messages. To exit press CTRL+C")
		<-forever
	}
}
