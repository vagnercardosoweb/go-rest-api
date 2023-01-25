package errors

import (
	"log"
)

func CheckPanicError(message string, err error) {
	if err != nil {
		log.Println(message)
		panic(err)
	}
}
