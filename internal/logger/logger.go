package logger

import (
	"io"
	"log"
	"os"
)

func Init(path string) (*os.File, error) {

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(io.MultiWriter(f, os.Stdout))

	return f, nil
}
