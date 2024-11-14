package listener

import (
	"context"
	"encoding/json"
	"log"
	"task-runner/cmd/task"

	"github.com/jmoiron/sqlx"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ContextTask string
type ContextDB string

func Listen(tr *task.TaskRepository, db *sqlx.DB) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	queue, err := ch.QueueDeclare("task-runner", false, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	msgs, err := ch.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		panic(err)
	}

	for msg := range msgs {
		handleMessage(msg, tr, db)
	}
}

func handleMessage(msg amqp.Delivery, tr *task.TaskRepository, db *sqlx.DB) {
	var message task.Message
	json.Unmarshal(msg.Body, &message)

	t, err := tr.GetTaskFromName(message.TaskName)
	if err != nil {
		log.Fatalln("can't find task ", message.TaskName)
	}

	ctx := context.WithValue(context.TODO(), ContextTask("task"), t)
	ctx = context.WithValue(ctx, ContextDB("db"), db)

	t.Fn(ctx, message)
}
