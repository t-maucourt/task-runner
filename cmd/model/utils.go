package model

import (
	"os"
)

var taskDB string

func GetTaskDBName() string {
	if taskDB != "" {
		return taskDB
	}

	taskDB = os.Getenv("DB_FILENAME")
	if taskDB == "" {
		panic("missing DB_FILENAME env variable")
	}

	return taskDB
}
