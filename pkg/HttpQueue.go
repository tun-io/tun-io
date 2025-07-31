package pkg

import "net/http"

type PendingHttpRequest struct {
	EventId  int64
	Request  *http.Request
	Response http.ResponseWriter
}
