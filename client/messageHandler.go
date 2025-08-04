package client

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/tun-io/tun-io/internal/http/Helpers"
	"github.com/tun-io/tun-io/internal/http/headers"
	"github.com/tun-io/tun-io/pkg"
	"io"
	"net/http"
	"strings"
)

var runningRequests = make(map[int64]bool)

func replaceBodyDomain(body []byte, resp *http.Response) []byte {
	if len(body) <= 0 {
		return body
	}

	contentEncoding := resp.Header.Get("Content-Encoding")
	var isGzip = false
	if contentEncoding == "gzip" {
		isGzip = true
		gzipReader, err := gzip.NewReader(bytes.NewReader(body))
		if err != nil {
			println("Error creating gzip reader:", err.Error())
			return body
		}
		defer gzipReader.Close()

		body, err = io.ReadAll(gzipReader)
		if err != nil {
			println("Error reading gzip response body:", err.Error())
			return body
		}
	}

	body = bytes.ReplaceAll(body, []byte(localDomain), []byte(remoteDomain))

	if !isGzip {
		return body
	}

	var buffer bytes.Buffer
	gzipWriter := gzip.NewWriter(&buffer)
	_, err := gzipWriter.Write(body)
	if err != nil {
		println("Error writing gzip response body:", err.Error())
		return body
	}

	err = gzipWriter.Close()
	if err != nil {
		println("Error closing gzip writer:", err.Error())
		return body
	}

	return buffer.Bytes()
}

func httpRequest(c *websocket.Conn, command pkg.Command) {

	if _, exists := runningRequests[command.EventId]; exists {
		println("[httpRequest] Command already running:", command.EventId)
		return
	}

	runningRequests[command.EventId] = true
	defer func() {
		delete(runningRequests, command.EventId)
		println("Removed event:", command.EventId)
	}()

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

	if payload.Url == "" {
		println("Received command with empty URL")
		return
	}

	payload.Url = strings.ReplaceAll(payload.Url, remoteDomain, localDomain)

	println("Received HTTP request command:", payload.Url, payload.Method)
	// Create a new HTTP request
	req, err := http.NewRequest(payload.Method, payload.Url, bytes.NewBuffer([]byte(payload.Body)))
	if err != nil {
		println("Error creating HTTP request:", err.Error())
		return
	}

	// Set headers if provided
	for key, value := range payload.Headers {
		if headers.IsDisallowedHeader(key) {
			continue
		}
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

	body = replaceBodyDomain(body, resp)

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

	delete(runningRequests, command.EventId)
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
