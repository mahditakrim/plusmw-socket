package shutdown

import (
	"os"
	"os/signal"
	"syscall"
)

func Graceful(f func(), sigChan chan os.Signal) {

	signal.Notify(sigChan, os.Interrupt,
		syscall.SIGHUP, syscall.SIGINT,
		syscall.SIGTERM, syscall.SIGQUIT)
	<-sigChan

	if f != nil {
		f()
	}
}
