#  Task runner

Program listening to RabbitMQ messages and executing corresponding functions based on the message content.

## Usage

```go
// main.go

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"task-runner/cmd/listener"
	"task-runner/cmd/model"
	"task-runner/cmd/tasks"
	"task-runner/cmd/transport"

	"github.com/joho/godotenv"
)

func task1(ctx context.Context, m transport.Message) {
	tsk := ctx.Value(listener.ContextTask("task")).(*tasks.Task)

	defer model.NewTaskModel().Save(tsk)

	var d map[string]string
	if err := json.Unmarshal(m.Data, &d); err != nil {
		tsk.Status = tasks.STATUS_ERROR
		log.Printf("Couldn't unmarshall data - %s", tsk)
		return
	}

	for k, v := range d {
		fmt.Println(k, v)
	}

	tsk.Status = tasks.STATUS_DONE
	log.Println("Task done", tsk)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("couldn't load .env file: " + err.Error())
	}

	log.Println("task-runner started")

    // Creating a new task repository that will hold all the tasks that can be executed in response to a received message from RabbitMQ
	taskRepository := tasks.NewTaskRepository()
	
    // Registering a task with the `task-1` name and targeting the `task1()` function
    taskRepository.RegisterTask("task-1", task1)

    // Starting to listen to RabbitMQ messages and reacting to them
	listener.Listen(taskRepository)
}
```

This example creates a new task named `task-1` that will run the function `task1(...)`.

To trigger this task we can simply insert a payload in RabbitMQ (using the UI or another script) with the following payload:

```json
{
    "client": "thomas",
    "task_name": "task-1",
    "data": {
        "hello": "world"
    }
}
```

When received the message will be sent to the `task1(...)` function and print out the keys & values received in the `data` field. When the function is done executing, the `Task` object will be stored in a SQLite database with the status of the task set to `DONE`

If the unmarshalling of the data fails, the task's status is updated to `ERROR` and the function stops executing 

In order to run the program, you need a JSON configuration file, and the filename needs to be added to your env variables under the name `CONFIG_FILE`

# TODO

* Unit tests
* Allow the addition of custom consumers
* ...
