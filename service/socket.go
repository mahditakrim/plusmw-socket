package service

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/gobwas/ws/wsutil"
	"github.com/mahditakrim/plusmw-socket/entity"
	"github.com/mahditakrim/plusmw-socket/repository"
)

type socket struct {
	repo repository.Repository
}

func NewSocket(repo repository.Repository) (Service, error) {

	if repo == nil {
		return nil, errors.New("nil repository reference")
	}

	return &socket{repo}, nil
}

func (s *socket) IterateOnConns(handler func(int64, entity.SocketValue)) {

	s.repo.LoadEach(handler)
}

func (s *socket) TerminateClients() {

	s.repo.LoadEach(func(_ int64, value entity.SocketValue) {
		value.Conn.Close()
	})
}

func (s *socket) StoreConn(ctx context.Context, userID int64, conn net.Conn) error {

	return s.repo.Insert(ctx, userID, entity.SocketValue{LastDelivery: time.Now(), Conn: conn})
}

func (s *socket) Send(data []byte, receiver int64) error {

	conn := s.repo.GetConn(receiver)
	if conn == nil {
		return nil
	}

	if err := wsutil.WriteServerText(conn, data); err != nil {
		return err
	}

	return s.repo.Update(receiver, entity.SocketValue{LastDelivery: time.Now(), Conn: conn})
}

func (s *socket) DeleteConn(userID int64) {

	value := s.repo.Delete(userID)
	if value != nil {
		value.Conn.Close()
	}
}
