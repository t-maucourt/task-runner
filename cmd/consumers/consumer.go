package consumers

import (
	"encoding/json"
	"task-runner/cmd/configuration"
	"task-runner/cmd/tasks"
	"task-runner/cmd/utils"
)

type ContextTask string

var taskRunnerConfig *configuration.TaskRunnerConfiguration

type BrokerConfig interface{}

func init() {
	taskRunnerConfig = configuration.GetConfig()
}

func Run(tr *tasks.TaskRepository) {
	switch taskRunnerConfig.Broker.Type {
	case RabbitMQ:
		consumeRabbitMQ(tr)
	default:
		panic("broker not supported " + taskRunnerConfig.Broker.Type)
	}
}

func readConfig(brokerType string) BrokerConfig {
	switch brokerType {
	case RabbitMQ:
		var config rabbitMQConfig
		err := json.Unmarshal(taskRunnerConfig.Broker.Configuration, &config)
		utils.PanicOnError(err, "wrong broker configuration")
		return config
	default:
		panic("broker configuration not supported")
	}
}
