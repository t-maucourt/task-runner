package task

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

const (
	STATUS_PENDING = "PENDING"
	STATUS_RUNNING = "RUNNING"
	STATUS_ERROR   = "ERROR"
	STATUS_DONE    = "DONE"
)

type Task struct {
	Name   string    `db:"task_name"`
	Status string    `db:"status"`
	TaskId uuid.UUID `db:"task_id"`
	Result string    `db:"result"`
	Fn     func(context.Context, Message)
}

func newTask(name string, fn func(context.Context, Message)) *Task {
	return &Task{Name: name, Status: STATUS_PENDING, TaskId: uuid.New(), Fn: fn}
}

func (t *Task) String() string {
	return fmt.Sprintf("Name : %s - ID : %s - Status : %s", t.Name, t.TaskId, t.Status)
}
