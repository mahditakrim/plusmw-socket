package gc

import (
	"errors"
	"log"
	"net"
	"time"

	"github.com/gobwas/ws"
	"github.com/mahditakrim/plusmw-socket/entity"
	"github.com/mahditakrim/plusmw-socket/luncher"
	"github.com/mahditakrim/plusmw-socket/service"
)

type gc struct {
	socket               service.Service
	stopChan, stopedChan chan struct{}
	interval             int
	cleaner              func(int64, entity.SocketValue)
}

func NewGC(socket service.Service, interval int) (luncher.Runnable, error) {

	if socket == nil {
		return nil, errors.New("nil socket reference")
	}
	if interval == 0 {
		return nil, errors.New("interval should be greater than 0")
	}

	return &gc{
		socket,
		make(chan struct{}),
		make(chan struct{}),
		interval,
		func(key int64, value entity.SocketValue) {
			if time.Now().After(value.LastDelivery.Add(time.Hour)) || !isAlive(value.Conn) {
				socket.DeleteConn(key)
			}
		},
	}, nil
}

func isAlive(conn net.Conn) bool {

	for range "01" {
		if ws.WriteFrame(conn, ws.NewPingFrame(nil)) != nil {
			return false
		}
		time.Sleep(time.Millisecond)
	}

	return true
}

func (gc *gc) Run() error {

	log.Printf("gc started with interval %ds\n", gc.interval)
	defer close(gc.stopedChan)

	for {
		for i := 0; i < gc.interval; i++ {
			select {
			case <-gc.stopChan:
				return nil
			default:
				time.Sleep(time.Second)
			}
		}

		log.Println("gc cleaner started")
		gc.socket.IterateOnConns(gc.cleaner)
		log.Println("gc cleaner finished")
	}
}

func (gc *gc) Shutdown() error {

	close(gc.stopChan)
	<-gc.stopedChan
	return nil
}
