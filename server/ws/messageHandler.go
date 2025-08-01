package ws

import (
	"encoding/base64"
	"encoding/json"
	"github.com/charmbracelet/log"
	"github.com/gorilla/websocket"
	"github.com/tun-io/tun-io/pkg"
	"strconv"
)

func httpResponse(c *websocket.Conn, command pkg.Command) {
	var payload, err = command.GetHttpResponsePayload()
	if err != nil {
		log.Warnf("Error getting HTTP response payload: %v", err.Error())
		return
	}

	if payload == nil {
		println("Received command with no payload")
		return
	}

	PendingRequest, exists := PendingRequests[strconv.FormatInt(command.EventId, 10)]
	if !exists {
		log.Warnf("No pending request found for event ID: %v", command.EventId)
		return
	}

	log.Infof("Received HTTP response command for event ID: %v", command.EventId)

	var responseWriter = PendingRequest.Response
	if responseWriter == nil {
		log.Errorf("Response writer is nil for event ID: %v", command.EventId)
		return
	}

	if len(payload.Body) == 0 && payload.StatusCode == 0 {
		log.Warnf("Received empty payload for event ID: %v", command.EventId)
		return
	}

	for key, value := range payload.Headers {
		responseWriter.Header().Set(key, value)
	}

	responseWriter.Header().Set("Content-Length", strconv.Itoa(len(payload.Body)))

	responseWriter.WriteHeader(payload.StatusCode)

	// base64 decode the body if it starts with "base64:"
	decodedBytes, err := base64.StdEncoding.DecodeString(payload.Body)
	if err != nil {
		log.Warnf("Error decoding base64 body for event ID %v: %v", command.EventId, err.Error())
	}

	if payload.Body != "" {
		_, err = responseWriter.Write(decodedBytes)
		if err != nil {
			log.Errorf("Error writing response body for event ID %v: %v", command.EventId, err.Error())
			return
		}
	}

	log.Infof("HTTP response sent for event ID: %v with status code: %d", command.EventId, payload.StatusCode)
	delete(PendingRequests, strconv.FormatInt(command.EventId, 10))
}

func messageHandler(c *websocket.Conn) {
	for {
		// Read messages from the connection
		_, message, err := c.ReadMessage()
		if err != nil {
			c.Close()
			return
		}

		var command pkg.Command
		err = json.Unmarshal(message, &command)

		if err != nil {
			println("Error unmarshalling message:", err.Error())
			continue
		}

		println("Received command:", command.Event, "with ID:", command.EventId)
		switch command.Event {
		case "http_response":
			httpResponse(c, command)
		default:
			log.Warnf("Unknown command received: %s", command.Event)
		}

	}
}
