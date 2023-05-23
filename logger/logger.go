package logger

import (
	"log"
	"os"
)

var (
	Logger *log.Logger
)

func InitLogger() {
	file, err := os.OpenFile("logging.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}

	Logger = log.New(file, "", log.LstdFlags)
}
