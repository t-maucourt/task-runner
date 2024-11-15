package model

import (
	"log"
	"task-runner/cmd/tasks"

	"github.com/jmoiron/sqlx"
)

type TaskModel struct {
	db *sqlx.DB
}

func NewTaskModel() *TaskModel {
	db := getDB()
	return &TaskModel{db}
}

func (tm *TaskModel) Save(t *tasks.Task) error {
	query := `INSERT INTO tasks (task_name, status, task_id, result) VALUES (?, ?, ?, ?)`
	_, err := tm.db.Exec(query, t.Name, t.Status, t.TaskId, t.Result)

	if err != nil {
		log.Println("Can't insert task in DB: ", err)
	}

	return err
}

func (tm *TaskModel) UpdateStatus(taskId string, status string) error {
	query := `UPDATE tasks SET status = ? WHERE task_id = ?`
	_, err := tm.db.Exec(query, status, taskId)

	if err != nil {
		log.Println("Can't update task status: ", err)
	}

	return err
}

func (tm *TaskModel) GetByTaskId(taskId string) (tasks.Task, error) {
	var t tasks.Task
	if err := tm.db.Get(&t, `SELECT * FROM tasks WHERE task_id = ?`, taskId); err != nil {
		return tasks.Task{}, err
	}

	return t, nil
}
