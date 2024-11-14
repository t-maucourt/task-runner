package model

import (
	"log"
	"task-runner/cmd/task"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const (
	TASKS_DB = "tasks.db"

	INITIALIZE_TABLE_STMT = `CREATE TABLE IF NOT EXISTS tasks (
	"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	"task_name" TEXT,
	"status" TEXT,
	"task_id" TEXT,
	"result" TEXT
	);`
)

func InitDB() *sqlx.DB {
	db, err := sqlx.Open("sqlite3", TASKS_DB)
	if err != nil {
		log.Panicln(err)
	}

	_, err = db.Exec(INITIALIZE_TABLE_STMT)
	if err != nil {
		log.Panicln(err)
	}

	return db
}

type TaskModel struct {
	db *sqlx.DB
}

func NewTaskModel(db *sqlx.DB) *TaskModel {
	return &TaskModel{db}
}

func (tm *TaskModel) Save(t *task.Task) error {
	stmt := `INSERT INTO tasks (task_name, status, task_id, result) VALUES (?, ?, ?, ?)`
	_, err := tm.db.Exec(stmt, t.Name, t.Status, t.TaskId, t.Result)

	if err != nil {
		log.Println("Can't insert task in DB : ", err)
	}

	return err
}

func (tm *TaskModel) UpdateStatus(taskId string, status string) error {
	stmt := `UPDATE tasks SET status = ? WHERE task_id = ?`
	_, err := tm.db.Exec(stmt, status, taskId)

	if err != nil {
		log.Println("Can't update task status : ", err)
	}

	return err
}

func (tm *TaskModel) GetByTaskId(taskId string) (task.Task, error) {
	var tsk task.Task
	if err := tm.db.Get(&tsk, `SELECT * FROM tasks WHERE task_id = ?`, taskId); err != nil {
		return task.Task{}, err
	}

	return tsk, nil
}
