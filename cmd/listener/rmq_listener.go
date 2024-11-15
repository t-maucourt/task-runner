package listener

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"task-runner/cmd/tasks"
	"task-runner/cmd/transport"
	"task-runner/cmd/utils"

	amqp "github.com/rabbitmq/amqp091-go"
)

var rqmURI string
var queueName string

type ContextTask string

func Listen(tr *tasks.TaskRepository) {
	loadConfiguration()

	conn, err := amqp.Dial(rqmURI)
	utils.PanicOnError(err, "error while dialing rabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	utils.PanicOnError(err, "error while opening a channel")
	defer ch.Close()

	queue, err := ch.QueueDeclare(queueName, false, false, false, false, nil)
	utils.PanicOnError(err, "error while declaring a queue")

	rmqMsgs, err := ch.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	utils.PanicOnError(err, "error while starting consuming the channel")

	for rmqMsg := range rmqMsgs {
		handleMessage(rmqMsg.Body, tr)
	}
}

func loadConfiguration() {
	rqmURI = os.Getenv("RMQ_URI")
	if rqmURI == "" {
		panic("missing RMQ_URI env variable")
	}

	queueName = os.Getenv("RMQ_QUEUE_NAME")
	if queueName == "" {
		panic("missing RMQ_QUEUE_NAME env variable")
	}
}

func handleMessage(rmqMsgBody []byte, tr *tasks.TaskRepository) {
	var msg transport.Message
	json.Unmarshal(rmqMsgBody, &msg)

	t, err := tr.GetTaskFromName(msg.TaskName)
	if err != nil {
		log.Println("Can't find task ", msg.TaskName)
		return
	}

	ctx := context.WithValue(context.TODO(), ContextTask("task"), t)

	t.Fn(ctx, msg)
}
