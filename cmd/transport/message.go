package transport

import "encoding/json"

type Message struct {
	Client   string          `json:"client"`
	TaskName string          `json:"task_name"`
	Data     json.RawMessage `json:"data"`
}
