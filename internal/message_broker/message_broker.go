package messagebroker

import (
	"errors"
	"log"
	"strconv"
	"sync"

	"github.com/mahditakrim/plusmw-socket/luncher"
	"github.com/mahditakrim/plusmw-socket/service"
	"github.com/nats-io/nats.go"
)

type messageBroker struct {
	conn    *nats.Conn
	msgChan chan *nats.Msg
	stoped  chan struct{}
	subject string
	handler func()
}

var conn *nats.Conn

func newConn(url string) (*nats.Conn, error) {

	if url == "" {
		return nil, errors.New("invalid connection address")
	}

	if conn != nil && !conn.IsClosed() {
		return conn, nil
	}

	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	conn = nc

	return conn, nil
}

func NewSubscriber(socket service.Service, addr, user, pass, subject string) (luncher.Runnable, error) {

	if socket == nil {
		return nil, errors.New("nil socket reference")
	}

	nc, err := newConn("nats://" + user + ":" + pass + "@" + addr)
	if err != nil {
		return nil, err
	}

	msgChan := make(chan *nats.Msg, 64)
	_, err = nc.ChanSubscribe(subject, msgChan)
	if err != nil {
		return nil, err
	}

	stopChan := make(chan struct{})
	return &messageBroker{
		nc,
		msgChan,
		stopChan,
		subject,
		func() {
			defer close(stopChan)
			handlersWaiter := sync.WaitGroup{}

			for {
				msg, ok := <-msgChan
				if !ok {
					handlersWaiter.Wait()
					return
				}

				handlersWaiter.Add(1)
				go func(msg *nats.Msg) {
					defer handlersWaiter.Done()

					receiver, err := strconv.Atoi(msg.Header.Get("receiver"))
					if err != nil {
						log.Println(err)
						return
					}

					err = socket.Send(msg.Data, int64(receiver))
					if err != nil {
						log.Println(err)
						socket.DeleteConn(int64(receiver))
					}
				}(msg)
			}
		},
	}, nil
}

func (mb *messageBroker) Run() error {

	log.Println("message-broker started consuming on", mb.subject)
	mb.handler()
	return nil
}

func (mb *messageBroker) Shutdown() error {

	mb.conn.Close()
	close(mb.msgChan)
	<-mb.stoped
	return nil
}
