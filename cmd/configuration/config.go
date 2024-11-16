package configuration

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"task-runner/cmd/utils"

	"github.com/joho/godotenv"
)

type TaskRunnerConfiguration struct {
	Broker broker `json:"broker"`
	DB     db
}

type broker struct {
	Type          string          `json:"type"`
	Configuration json.RawMessage `json:"configuration"`
}

type db struct {
	Filename string `json:"filename"`
}

var Configuration *TaskRunnerConfiguration

func GetConfig() *TaskRunnerConfiguration {
	if Configuration != nil {
		return Configuration
	}

	err := godotenv.Load()
	if err != nil {
		panic("couldn't load .env file: " + err.Error())
	}

	configFileName := os.Getenv("CONFIG_FILE")
	if configFileName == "" {
		panic("missing CONFIG_FILE env variable")
	}

	Configuration = loadConfig(configFileName)
	return Configuration
}

func loadConfig(configFileName string) *TaskRunnerConfiguration {
	file, err := os.Open(configFileName)
	utils.PanicOnError(err, fmt.Sprintf("error while opening config file %s", configFileName))

	v, _ := io.ReadAll(file)

	var config TaskRunnerConfiguration
	json.Unmarshal(v, &config)

	return &config
}
