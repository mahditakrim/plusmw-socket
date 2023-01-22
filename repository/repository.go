package repository

import (
	"context"
	"net"

	"github.com/mahditakrim/plusmw-socket/entity"
)

type Repository interface {
	Insert(ctx context.Context, key int64, value entity.SocketValue) error
	Update(key int64, value entity.SocketValue) error
	GetConn(key int64) net.Conn
	LoadEach(handler func(key int64, value entity.SocketValue))
	Delete(key int64) *entity.SocketValue
}
