package client

import (
	"github.com/cenkalti/backoff/v4"
	"github.com/gorilla/websocket"
	"github.com/tun-io/tun-io/internal/http/Helpers"
	"net/url"
	"time"
)

func createConnection() *websocket.Conn {
	scheme := "ws"
	if isSecure {
		scheme = "wss"
	}

	u := url.URL{Scheme: scheme, Host: Helpers.NormaliseUrl(remoteDomain), Path: "/_tunio/_internal/api/client/ws"}
	var conn *websocket.Conn
	var err error
	operation := func() error {
		conn, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			println("Failed to connect to server:", err.Error())
			return err
		}
		println("Connected to server:", u.String())
		return nil
	}
	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = 0 // never stop retrying
	backoff.Retry(operation, b)
	return conn
}

func Connect() {
	var connection *websocket.Conn
	for {
		connection = createConnection()
		if connection == nil {
			println("Exiting client due to connection failure.")
			return
		}

		println("Connected to server:", connection.RemoteAddr().String())
		messageHandler(connection)
		println("Connection lost. Attempting to reconnect...")
		// Wait before reconnecting
		time.Sleep(2 * time.Second)
	}
}
