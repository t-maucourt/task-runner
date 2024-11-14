package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"task-runner/cmd/listener"
	"task-runner/cmd/model"
	"task-runner/cmd/task"

	"github.com/jmoiron/sqlx"
)

type TaskFn func(context.Context, task.Message)

func task1(ctx context.Context, m task.Message) {
	tsk := ctx.Value(listener.ContextTask("task")).(*task.Task)
	dbo := ctx.Value(listener.ContextDB("db")).(*sqlx.DB)

	defer model.NewTaskModel(dbo).Save(tsk)

	var d map[string]string
	if err := json.Unmarshal(m.Data, &d); err != nil {
		tsk.Status = task.STATUS_ERROR
		log.Printf("Couldn't unmarshall data - %s", tsk)
		return
	}

	for k, v := range d {
		fmt.Println(k, v)
	}

	tsk.Status = task.STATUS_DONE
	log.Println("Task done", tsk)
}

func task2(ctx context.Context, m task.Message) {
	var data int
	json.Unmarshal(m.Data, &data)
	fmt.Println(data)
}

func setupTasks() *task.TaskRepository {

	taskRepository := task.NewTaskRepository()
	taskRepository.RegisterTask("task-1", task1)
	taskRepository.RegisterTask("task-2", task2)

	return taskRepository
}

func main() {
	fmt.Println("Task runner")

	db := model.InitDB()

	taskRepository := setupTasks()

	listener.Listen(taskRepository, db)
}
