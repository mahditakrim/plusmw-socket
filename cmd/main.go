package main

import (
	"log"

	"github.com/mahditakrim/plusmw-socket/config"
	"github.com/mahditakrim/plusmw-socket/internal/logger"
	"github.com/mahditakrim/plusmw-socket/luncher"
	"github.com/mahditakrim/plusmw-socket/setup"
)

func main() {

	conf, err := config.Init()
	if err != nil {
		log.Panic(err)
	}

	f, err := logger.Init(conf.LogPath)
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()

	runnables, err := setup.Init(conf)
	if err != nil {
		log.Panic(err)
	}

	luncher.Start(runnables...)
}
