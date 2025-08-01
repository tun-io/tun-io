package client

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/tun-io/tun-io/internal/http/Helpers"
	"github.com/tun-io/tun-io/pkg"
	"io"
	"net/http"
	"strings"
)

func httpRequest(c *websocket.Conn, command pkg.Command) {
	var payload, err = command.GetHttpRequestPayload()
	if err != nil {
		println("Error getting HTTP request payload:", err.Error())
		return
	}
	// turn the body into an io.Reader
	if payload == nil {
		println("Received command with no payload")
		return
	}

	// replace payload.Url the original host to localhost:8000
	if payload.Url == "" {
		println("Received command with empty URL")
		return
	}

	payload.Url = strings.ReplaceAll(payload.Url, "a.tunio.test", "localhost:8000")

	println("Received HTTP request command:", payload.Url, payload.Method)
	// Create a new HTTP request
	req, err := http.NewRequest(payload.Method, payload.Url, bytes.NewBuffer([]byte(payload.Body)))
	if err != nil {
		println("Error creating HTTP request:", err.Error())
		return
	}

	// Set headers if provided
	for key, value := range payload.Headers {
		req.Header.Set(key, value)
	}

	// Send the HTTP request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		println("Error sending HTTP request:", err.Error())
		return
	}

	var body []byte
	if resp.Body != nil {
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			println("Error reading response body:", err.Error())
			return
		}
	}

	var encodedBody = base64.StdEncoding.EncodeToString(body)

	// Prepare the response payload
	responsePayload := pkg.HttpResponsePayload{
		StatusCode: resp.StatusCode,
		Headers:    Helpers.HeadersToMap(resp.Header),
		Body:       encodedBody,
	}

	var responseCommand = pkg.Command{
		Version: "1.0",
		Event:   "http_response",
		EventId: command.EventId,
		Payload: responsePayload,
		SendAt:  command.SendAt,
	}

	c.WriteJSON(responseCommand)
}

func messageHandler(c *websocket.Conn) {
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			println("Error reading message:", err.Error())
			break
		}

		var command pkg.Command
		err = json.Unmarshal(message, &command)

		if err != nil {
			println("Error unmarshalling message:", err.Error())
			continue
		}

		switch command.Event {
		case "http_request":
			httpRequest(c, command)
		default:
			println("Unknown command received:", command.Event)
		}
	}
}
