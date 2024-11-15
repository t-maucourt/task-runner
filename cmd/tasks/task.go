package tasks

import (
	"context"
	"fmt"
	"task-runner/cmd/transport"

	"github.com/google/uuid"
)

const (
	STATUS_CREATED = "CREATED"
	STATUS_PENDING = "PENDING"
	STATUS_RUNNING = "RUNNING"
	STATUS_ERROR   = "ERROR"
	STATUS_DONE    = "DONE"
)

type TaskFn func(context.Context, transport.Message)

type Task struct {
	Name   string    `db:"task_name"`
	Status string    `db:"status"`
	TaskId uuid.UUID `db:"task_id"`
	Result string    `db:"result"`
	Fn     TaskFn
}

func newTask(name string, fn TaskFn) *Task {
	return &Task{Name: name, Status: STATUS_CREATED, TaskId: uuid.New(), Fn: fn}
}

func (t *Task) String() string {
	return fmt.Sprintf("Name : %s - ID : %s - Status : %s", t.Name, t.TaskId, t.Status)
}
