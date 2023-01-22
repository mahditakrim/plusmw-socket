package websocket

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gobwas/ws"
	"github.com/mahditakrim/plusmw-socket/luncher"
	"github.com/mahditakrim/plusmw-socket/service"
)

type websocket struct {
	server http.Server
	addr   string
}

func NewWebsocket(socket service.Service, addr string) (luncher.Runnable, error) {

	if socket == nil {
		return nil, errors.New("nil socket reference")
	}

	server := http.Server{
		Addr: addr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			conn, _, _, err := ws.UpgradeHTTP(r, w)
			if err != nil {
				log.Println(err)
				return
			}

			userID, err := strconv.Atoi(r.Header.Get("user_id"))
			if err != nil {
				conn.Close()
				log.Println(err)
				return
			}

			err = socket.StoreConn(r.Context(), int64(userID), conn)
			if err != nil {
				conn.Close()
				log.Println(err)
			}
		}),
	}

	server.RegisterOnShutdown(socket.TerminateClients)
	return &websocket{server, addr}, nil
}

func (w *websocket) Run() error {

	log.Println("websocket server started at", w.addr)
	if err := w.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (w *websocket) Shutdown() error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	return w.server.Shutdown(ctx)
}
