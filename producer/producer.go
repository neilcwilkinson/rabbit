package producer

import (
	"fmt"

	"github.com/streadway/amqp"
)

var (
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue

	//queuename   string
)

func failOnError(err error, msg string) {
	if err != nil {
		// log.Fatalf("%s: %s", msg, err)
		// panic(fmt.Sprintf("%s: %s", msg, err))
		fmt.Println("Error:", err, msg)
	}
}

func Initialize(uri string, queuename string) {
	conn, err := amqp.Dial(uri)
	// conn, err := amqp.Dial("amqp://fahoobagats:Fytg*pv2c,iCUT@rabbitmq-statsnode-3a5n:5672/")
	//failOnError(err, "Failed to connect to RabbitMQ")

	if err != nil {
		//log.Fatalf("%s: %s", msg, err)
		//panic(fmt.Sprintf("%s: %s", msg, err))
		fmt.Println("Rabbit MQ Producer Error:", err)
	} else {
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
