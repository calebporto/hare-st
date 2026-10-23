package controller

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"uuid"

	"github.com/calebportodev/golang/hare-st/src/internal/app"
	"github.com/calebportodev/golang/hare-st/src/internal/protocol"
)

type InvalidRequestResponse struct {
	Error   string `json:"error"`
	Data    string `json:"data"`
	Message string `json:"message"`
}

type Controller struct {
	Conn net.Conn
}

func (c *Controller) HandleConnection() {
	defer c.Conn.Close()
	reader := bufio.NewReader(c.Conn)
	conexao := uuid.NewV4().String()
	for {
		data, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				slog.Info("Connection closed", "connection", conexao)
				return
			}
			slog.Error("Error reading from connection", "error", err.Error())
			c.Response(fmt.Sprintf("Error reading from connection: %s\n", err.Error()))
			continue
		}
		parser := protocol.Parser{}

		cleanData, err := parser.CleanRawRequest(data)
		if err != nil {
			slog.Error("Error cleaning raw request", "error", err.Error())
			c.Response(fmt.Sprintf("Invalid request: %s\n", err.Error()))
			continue
		}

		parsedRequest, parseRequestErr := parser.ParseRawRequest(cleanData)
		if parseRequestErr != nil {
			response := InvalidRequestResponse{
				Error:   "invalid_request",
				Data:    cleanData,
				Message: parseRequestErr.Error(),
			}
			responseJsonBytes, responseJsonBytesErr := json.Marshal(response)
			if responseJsonBytesErr != nil {
				slog.Error("Error marshalling JSON", "error", responseJsonBytesErr.Error())
				c.Response(fmt.Sprintf("Error: %s; Data: %s; Message: %s\n", response.Error, response.Data, response.Message))
				continue
			}
			slog.Error(response.Error, "data", response.Data, "message", response.Message)
			c.Response(string(responseJsonBytes) + "\n")
			continue
		}
		// Handle valid request
		var storageMap = app.Storage
		var app = app.RequestHandler{
			StorageMap: storageMap,
		}
		response, handleErr := app.HandleRequest(parsedRequest)
		if handleErr != nil {
			slog.Error("Error handling request", "error", handleErr.Error())
			c.Response(fmt.Sprintf("Error: %s\n", handleErr.Error()))
			continue
		}

		slog.Info("Request handled successfully", "request", parsedRequest, "response", response)
		c.Response(response + "\n")
	}
}

func (c *Controller) Response(r string) {
	c.Conn.Write([]byte(r))
}
