package consumers

import (
	"context"
	"encoding/json"
	"log"
	"task-runner/cmd/tasks"
	"task-runner/cmd/transport"
	"task-runner/cmd/utils"

	amqp "github.com/rabbitmq/amqp091-go"
)

const RabbitMQ = "RabbitMQ"

type rabbitMQConfig struct {
	URI       string `json:"uri"`
	QueueName string `json:"queue_name"`
}

func consumeRabbitMQ(tr *tasks.TaskRepository) {
	config := readConfig(RabbitMQ).(rabbitMQConfig)

	conn, err := amqp.Dial(config.URI)
	utils.PanicOnError(err, "error while dialing rabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	utils.PanicOnError(err, "error while opening a channel")
	defer ch.Close()

	queue, err := ch.QueueDeclare(config.QueueName, false, false, false, false, nil)
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
		go handleMessage(rmqMsg.Body, tr)
	}
}

func handleMessage(messageBody []byte, tr *tasks.TaskRepository) {
	var msg transport.Message
	json.Unmarshal(messageBody, &msg)

	t, err := tr.GetTaskFromName(msg.TaskName)
	if err != nil {
		log.Println("Can't find task ", msg.TaskName)
		return
	}

	ctx := context.WithValue(context.TODO(), ContextTask("task"), t)

	t.Fn(ctx, msg)
}
