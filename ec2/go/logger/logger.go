package logger

import (
	"io"
	"log"
	"os"
)

var (
	Info  *log.Logger
	Error *log.Logger
)

func Init(logFilePath string) {
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	multi := io.MultiWriter(os.Stdout, file)

	Info = log.New(multi, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(multi, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
