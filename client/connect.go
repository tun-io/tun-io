package client

import (
	"github.com/gorilla/websocket"
	"net/url"
)

func createConnection() *websocket.Conn {

	scheme := "ws"
	if isSecure {
		scheme = "wss"
	}

	u := url.URL{Scheme: scheme, Host: remoteDomain, Path: "/_tunio/_internal/api/client/ws"}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		println("Failed to connect to server:", err.Error())
		return nil
	}

	println("Connected to server:", u.String())
	return conn
}

func Connect() {
	connection := createConnection()
	if connection == nil {
		println("Exiting client due to connection failure.")
		return
	}

	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			println("Failed to close connection:", err.Error())
		} else {
			println("Connection closed successfully")
		}
	}(connection)

	println("Connected to server:", connection.RemoteAddr().String())
	messageHandler(connection)
	println("Connection closed, exiting client.")
}
