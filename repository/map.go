package repository

import (
	"context"
	"errors"
	"net"
	"sync"

	"github.com/mahditakrim/plusmw-socket/entity"
)

type connStore struct {
	connMap sync.Map
}

func NewConnStore() Repository {

	return &connStore{sync.Map{}}
}

func (cs *connStore) Update(key int64, value entity.SocketValue) error {

	if value.Conn == nil {
		return errors.New("nil conn reference")
	}

	cs.connMap.Store(key, value)
	return nil
}

func (cs *connStore) LoadEach(handler func(int64, entity.SocketValue)) {

	cs.connMap.Range(func(k, v interface{}) bool {
		handler(k.(int64), v.(entity.SocketValue))
		return true
	})
}

func (cs *connStore) GetConn(key int64) net.Conn {

	v, _ := cs.connMap.Load(key)
	if v == nil {
		return nil
	}

	return v.(entity.SocketValue).Conn
}

func (cs *connStore) Insert(ctx context.Context, key int64, value entity.SocketValue) error {

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		if value.Conn == nil {
			return errors.New("nil conn reference")
		}
		cs.connMap.Store(key, value)
		return nil
	}
}

func (cs *connStore) Delete(key int64) *entity.SocketValue {

	v, _ := cs.connMap.LoadAndDelete(key)
	if v == nil {
		return nil
	}

	value := v.(entity.SocketValue)
	return &value
}
