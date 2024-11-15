package utils

import "fmt"

func PanicOnError(err error, msg string) {
	if err == nil {
		return
	}

	var errorMsg string
	if msg == "" {
		errorMsg = err.Error()
	} else {
		errorMsg = fmt.Sprintf("%s: %s", msg, err.Error())
	}

	panic(errorMsg)
}
