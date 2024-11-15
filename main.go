package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"task-runner/cmd/listener"
	"task-runner/cmd/model"
	"task-runner/cmd/tasks"
	"task-runner/cmd/transport"

	"github.com/joho/godotenv"
)

func task1(ctx context.Context, m transport.Message) {
	tsk := ctx.Value(listener.ContextTask("task")).(*tasks.Task)

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

func task2(ctx context.Context, m transport.Message) {
	var data int
	json.Unmarshal(m.Data, &data)
	fmt.Println(data)
}

func task3(ctx context.Context, m transport.Message) {
	type D struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	var d D
	err := json.Unmarshal(m.Data, &d)
	fmt.Println(err)
	fmt.Println(d)
}

func setupTasks() *tasks.TaskRepository {
	taskRepository := tasks.NewTaskRepository()
	taskRepository.RegisterTask("task-1", task1)
	taskRepository.RegisterTask("task-2", task2)
	taskRepository.RegisterTask("task-3", task3)

	return taskRepository
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("couldn't load .env file: " + err.Error())
	}

	log.Println("task-runner started")

	taskRepository := setupTasks()

	listener.Listen(taskRepository)
}
