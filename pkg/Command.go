package pkg

import (
	"encoding/json"
	"errors"
)

type HttpRequestPayload struct {
	Url     string            `json:"url"`               // The URL to which the HTTP request is made
	Method  string            `json:"method"`            // The HTTP method (GET, POST, etc.)
	Headers map[string]string `json:"headers,omitempty"` // Optional headers for the HTTP request
	Body    string            `json:"body,omitempty"`    // Optional body for the HTTP request
}

type HttpResponsePayload struct {
	StatusCode int               `json:"status_code"`       // HTTP status code of the response
	Headers    map[string]string `json:"headers,omitempty"` // Headers returned in the response
	Body       string            `json:"body,omitempty"`    // Body of the response
}

type Command struct {
	Version string      `json:"version"`           // Version of the command, used for compatibility checks
	Event   string      `json:"event"`             // Event name of the event, used to identify the type of command (e.g. "http_request", "http_response")
	EventId int64       `json:"event_id"`          // EventId Unique identifier for the event, used to match requests and responses (e.g. HTTP request + response)
	Payload interface{} `json:"payload,omitempty"` // Payload of the command, can be of any type depending on the event (e.g. HttpRequestPayload, HttpResponsePayload, etc.)
	SendAt  int64       `json:"send_at,omitempty"` // SendAt is used for checking latency between server and client
}

// GetHttpResponsePayload returns the HttpResponsePayload if the event is "http_response", otherwise returns nil.
func (c *Command) GetHttpResponsePayload() (*HttpResponsePayload, error) {
	if c.Event != "http_response" {
		return nil, errors.New("incorrect event type") // Return nil if the event is not "http_response"
	}
	var payloadJson, err = json.Marshal(c.Payload)

	if err != nil {
		return nil, err
	}

	var payload = &HttpResponsePayload{}
	err = json.Unmarshal(payloadJson, payload)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

// GetHttpRequestPayload returns the HttpRequestPayload if the event is "http_request", otherwise returns nil.
func (c *Command) GetHttpRequestPayload() (*HttpRequestPayload, error) {
	if c.Event != "http_request" {
		return nil, errors.New("incorrect event type") // Return nil if the event is not "http_request"
	}

	var payloadJson, err = json.Marshal(c.Payload)
	if err != nil {
		return nil, err
	}

	var payload = &HttpRequestPayload{}
	err = json.Unmarshal(payloadJson, payload)
	if err != nil {
		return nil, err
	}

	return payload, nil
}
