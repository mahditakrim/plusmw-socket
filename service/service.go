package service

import (
	"context"
	"net"

	"github.com/mahditakrim/plusmw-socket/entity"
)

type Service interface {
	StoreConn(ctx context.Context, userID int64, conn net.Conn) error
	Send(data []byte, receiver int64) error
	DeleteConn(userID int64)
	IterateOnConns(handler func(key int64, value entity.SocketValue))
	TerminateClients()
}
