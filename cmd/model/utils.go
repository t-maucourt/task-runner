package model

import (
	"task-runner/cmd/configuration"
)

var taskDB string

func GetTaskDBName() string {
	if taskDB != "" {
		return taskDB
	}

	taskDB = configuration.GetConfig().DB.Filename
	if taskDB == "" {
		panic("missing DB_FILENAME env variable")
	}

	return taskDB
}
