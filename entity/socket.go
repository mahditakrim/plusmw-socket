package entity

import (
	"net"
	"time"
)

type SocketValue struct {
	LastDelivery time.Time
	Conn         net.Conn
}
