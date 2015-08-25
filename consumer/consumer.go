package consumer

import (
	"fmt"

	"github.com/streadway/amqp"
)

var (
	conn           *amqp.Connection
	channel        *amqp.Channel
	queue          amqp.Queue
	messageChannel = make(chan []byte)
)

func failOnError(err error, msg string) {
	if err != nil {
		//log.Fatalf("%s: %s", msg, err)
		//panic(fmt.Sprintf("%s: %s", msg, err))
		fmt.Println("Error:", err)
	}
}

func Initialize(uri string, queuename string) {
	conn, err := amqp.Dial(uri)
	// conn, err := amqp.Dial("amqp://fahoobagats:Fytg*pv2c,iCUT@rabbitmq-statsnode-3a5n:5672/")
	//failOnError(err, "Failed to connect to RabbitMQ")

	if err != nil {
		//log.Fatalf("%s: %s", msg, err)
		//panic(fmt.Sprintf("%s: %s", msg, err))
		fmt.Println("Rabbit MQ Consumer Error:", err)
	} else {
		fmt.Println("Rabbit MQ Consumer Connected!!")
		defer conn.Close()

		channel, err = conn.Channel()

		failOnError(err, "Failed to open a channel")
		//defer ch.Close()

		queue, err = channel.QueueDeclare(
			queuename, // name
			false,     // durable
			false,     // delete when usused
			false,     // exclusive
			false,     // no-wait
			nil,       // arguments
		)

		failOnError(err, "Failed to declare a queue")

		// msgs, err := channel.Consume(
		// 	queue.Name, // queue
		// 	"",         // consumer
		// 	true,       // auto-ack
		// 	false,      // exclusive
		// 	false,      // no-local
		// 	false,      // no-wait
		// 	nil,        // args
		// )
		// failOnError(err, "Failed to register a consumer")
		go receiveMessages()
		Consume()

		// forever := make(chan bool)

		// go func() {
		// 	for d := range msgs {
		// 		//if string(d.Body) == "done" {
		// 		fmt.Printf("\nReceived all messages: %s\n", d.Body)
		// 		//}
		// 	}
		// }()

		// fmt.Printf(" [*] Waiting for messages. To exit press CTRL+C")
		// <-forever
	}
}

func Consume() {
	// forever := make(chan bool)
	for {
		msgs, err := channel.Consume(
			queue.Name, // queue
			"",         // consumer
			true,       // auto-ack
			false,      // exclusive
			false,      // no-local
			false,      // no-wait
			nil,        // args
		)
		failOnError(err, "Failed to register a consumer")

		go func() {
			//for d := range msgs {
			for d := range msgs {
				messageChannel <- d.Body
				// if string(d.Body) == "done" {
				//fmt.Printf("\nReceived message: %s\n", d.Body)
				// }
			}
		}()
	}
	// <-forever
}

func receiveMessages() {
	go func() {
		for {
			_ = <-messageChannel
			//message := <-messageChannel
			//fmt.Printf("State changed: %s\n", message)
		}
	}()
}
