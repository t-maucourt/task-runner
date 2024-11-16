package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"task-runner/cmd/consumers"
	"task-runner/cmd/model"
	"task-runner/cmd/tasks"
	"task-runner/cmd/transport"

	"github.com/joho/godotenv"
)

func task1(ctx context.Context, m transport.Message) {
	tsk := ctx.Value(consumers.ContextTask("task")).(*tasks.Task)

	defer model.NewTaskModel().Save(tsk)

	var d map[string]string
	if err := json.Unmarshal(m.Data, &d); err != nil {
		tsk.Status = tasks.STATUS_ERROR
		log.Printf("Couldn't unmarshall data - %s", tsk)
		return
	}

	for k, v := range d {
		fmt.Println(k, v)
	}

	tsk.Status = tasks.STATUS_DONE
	log.Println("Task done", tsk)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("couldn't load .env file: " + err.Error())
	}

	log.Println("task-runner started")

	taskRepository := tasks.NewTaskRepository()
	taskRepository.RegisterTask("task-1", task1)

	consumers.Run(taskRepository)
}
