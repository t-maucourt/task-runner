package model

import (
	"task-runner/cmd/utils"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var db *sqlx.DB

var initTasksTableSQL = `CREATE TABLE IF NOT EXISTS tasks (
	"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	"task_name" TEXT,
	"status" TEXT,
	"task_id" TEXT,
	"result" TEXT
	);`

func getDB() *sqlx.DB {
	if db != nil {
		return db
	}

	db, err := sqlx.Open("sqlite3", GetTaskDBName())
	utils.PanicOnError(err, "error while opening database file")

	_, err = db.Exec(initTasksTableSQL)
	utils.PanicOnError(err, "error while executing initializing the tasks table")

	return db
}
