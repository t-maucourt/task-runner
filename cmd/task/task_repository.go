package task

import (
	"context"
	"fmt"
)

type TaskRepository struct {
	tasks []*Task
}

func (tr *TaskRepository) RegisterTask(name string, fn func(context.Context, Message)) {
	tr.tasks = append(tr.tasks, newTask(name, fn))
}

func (tr *TaskRepository) ListTasks() {
	for _, task := range tr.tasks {
		fmt.Println(task)
	}
}

func (tr *TaskRepository) GetTaskFromName(name string) (*Task, error) {
	for _, task := range tr.tasks {
		if task.Name == name {
			return task, nil
		}
	}

	return &Task{}, fmt.Errorf("can't find task with name %s", name)
}

func NewTaskRepository() *TaskRepository {
	return &TaskRepository{}
}
