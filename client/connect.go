package client

import (
	"github.com/gorilla/websocket"
	"net/url"
	"strconv"
)

func Connect() {
	var domain = "a.tunio.test"
	var port = 8080
	var target = domain + ":" + strconv.Itoa(port)

	println("Connecting to target:", target)

	u := url.URL{Scheme: "ws", Host: target, Path: "/_tunio/_internal/api/client/ws"}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		println("Failed to connect to server:", err.Error())
		return
	}

	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			println("Failed to close connection:", err.Error())
		} else {
			println("Connection closed successfully")
		}
	}(conn)

	println("Connected to server:", u.String())
	messageHandler(conn)
	println("Connection closed, exiting client.")
}
