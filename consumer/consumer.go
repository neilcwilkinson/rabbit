package consumer

import (
	"fmt"

	"github.com/streadway/amqp"
)

var (
	rabbitmquri = ""
	queuename   = ""

	conn           *amqp.Connection
	channel        *amqp.Channel
	queue          amqp.Queue
	messageChannel = make(chan []byte)

	connectionChannel = make(chan bool)
)

//Exposed for calling go code by capitalizing the first character.
var (
	Connected = false
)

func failOnError(err error, msg string) {
	if err != nil {
		//log.Fatalf("%s: %s", msg, err)
		//panic(fmt.Sprintf("%s: %s", msg, err))
		fmt.Printf("\nRabbit MQ Consumer Error: %s %s\n", err, msg)
		connectionChannel <- false
	}
}

func receiveConnectionStatus() {
	go func() {
		for {
			isConnected := <-connectionChannel

			if isConnected != Connected {
				fmt.Printf("State changed: %b\n", isConnected)
			}
			Connected = isConnected

			if Connected == false {
				//this only works in a go routine
				fmt.Println("Attempting Connection\n")
				go Initialize(rabbitmquri, queuename)
			}
		}
	}()
}

func Initialize(uri string, queue_name string) {
	go receiveConnectionStatus()

	rabbitmquri = uri
	queuename = queue_name

	conn, err := amqp.Dial(rabbitmquri)
	// conn, err := amqp.Dial("amqp://fahoobagats:Fytg*pv2c,iCUT@rabbitmq-statsnode-3a5n:5672/")
	//failOnError(err, "Failed to connect to RabbitMQ")

	if err != nil {
		//log.Fatalf("%s: %s", msg, err)
		//panic(fmt.Sprintf("%s: %s", msg, err))
		failOnError(err, "Connection unavailable")
	} else {
		connectionChannel <- true
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

		go receiveMessages()
		Consume()

	}
}

func Consume() {
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

		for d := range msgs {
			messageChannel <- d.Body
		}

	}
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
