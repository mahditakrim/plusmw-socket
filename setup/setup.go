package setup

import (
	"github.com/mahditakrim/plusmw-socket/config"
	"github.com/mahditakrim/plusmw-socket/internal/gc"
	messagebroker "github.com/mahditakrim/plusmw-socket/internal/message_broker"
	"github.com/mahditakrim/plusmw-socket/internal/websocket"
	"github.com/mahditakrim/plusmw-socket/luncher"
	"github.com/mahditakrim/plusmw-socket/repository"
	"github.com/mahditakrim/plusmw-socket/service"
)

func Init(conf *config.Config) ([]luncher.Runnable, error) {

	service, err := service.NewSocket(repository.NewConnStore())
	if err != nil {
		return nil, err
	}

	ws, err := websocket.NewWebsocket(service, conf.Transport.Websocket.Addr)
	if err != nil {
		return nil, err
	}

	mb, err := messagebroker.NewSubscriber(service,
		conf.Transport.MessageBroker.Nats.Addr,
		conf.Transport.MessageBroker.Nats.Username,
		conf.Transport.MessageBroker.Nats.Password,
		conf.Transport.MessageBroker.Nats.Subjects["messages"])
	if err != nil {
		return nil, err
	}

	gc, err := gc.NewGC(service, 5*60)
	if err != nil {
		return nil, err
	}

	return []luncher.Runnable{ws, mb, gc}, nil
}
