package server

import (
	"log/slog"
	"net"

	"github.com/calebportodev/golang/hare-st/src/internal/controller"
)

func Start(addr string) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		slog.Error("Error starting server", "error", err.Error())
		return
	}
	defer listener.Close()

	slog.Info("Server listening on", "address", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			slog.Error("Error accepting connection", "error", err.Error())
			continue
		}
		controller := &controller.Controller{
			Conn: conn,
		}
		go controller.HandleConnection()
	}
}
