package luncher

import (
	"log"
	"os"

	"github.com/mahditakrim/plusmw-socket/internal/shutdown"
)

func Start(items ...Runnable) {

	for _, item := range items {
		go func(i Runnable) {
			if err := i.Run(); err != nil {
				log.Panic(err)
			}
		}(item)

		defer func(i Runnable) {
			if err := i.Shutdown(); err != nil {
				log.Println(err)
			}
		}(item)
	}

	shutdown.Graceful(func() { log.Println("shutting down") }, make(chan os.Signal))
}
