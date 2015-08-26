package producer

import (
	"fmt"

	"github.com/streadway/amqp"
)

var (
	rabbitmquri = ""
	queuename   = ""
	conn        *amqp.Connection
	channel     *amqp.Channel
	queue       amqp.Queue

	connectionChannel = make(chan bool)

	//queuename   string
)

//Exposed for calling go code by capitalizing the first character.
var (
	Connected = false
)

func failOnError(err error, msg string) {
	if err != nil {
		// log.Fatalf("%s: %s", msg, err)
		// panic(fmt.Sprintf("%s: %s", msg, err))
		fmt.Printf("\nRabbit MQ Producer Error: %s %s\n", err, msg)
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
		fmt.Println("Rabbit MQ Producer Connected!!")
		// defer conn.Close()

		channel, err = conn.Channel()
		failOnError(err, "Failed to open a channel")
		// defer channel.Close()

		queue, err = channel.QueueDeclare(
			queuename, // name
			false,     // durable
			false,     // delete when usused
			false,     // exclusive
			false,     // no-wait
			nil,       // arguments
		)
		failOnError(err, "Failed to declare a queue")

		fmt.Printf("\nQueue created: %s\n", queue.Name)

		i := 0
		for i = 0; i < 10; i++ {
			Produce(fmt.Sprintf("hello %d", i))
		}
	}
}

func Produce(msg string) {
	if Connected == true {
		body := msg
		err := channel.Publish(
			"",         // exchange
			queue.Name, // routing key
			false,      // mandatory
			false,      // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		failOnError(err, "Failed to publish a message")
	}
}
