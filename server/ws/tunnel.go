package ws

import (
	"github.com/tun-io/tun-io/internal/http/Helpers"
	"github.com/tun-io/tun-io/pkg"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var PendingRequests map[string]pkg.PendingHttpRequest = make(map[string]pkg.PendingHttpRequest)

func randBetween(min int, max int) int64 {
	return int64(min + rand.Intn(max-min))
}

func SendTunnelRequest(subdomain string, r *http.Request, w http.ResponseWriter) {
	// Send the request to the tunnel server
	conn, exists := Connections[subdomain]
	if !exists {
		http.Error(w, "Tunnel not found", http.StatusNotFound)
		return
	}

	var scheme = "http"
	if r.TLS != nil {
		scheme = "https"
	}

	if r.URL.Scheme != "" {
		scheme = r.URL.Scheme
	}

	var eventId = time.Now().Unix() + int64(len(PendingRequests)) + randBetween(1000, 999999)
	var command = pkg.Command{
		Version: "1.0",
		Event:   "http_request",
		EventId: eventId,
		Payload: pkg.HttpRequestPayload{
			Url:     scheme + "://" + Helpers.GetDomain(r) + r.URL.Path + "?" + r.URL.RawQuery,
			Method:  strings.ToUpper(r.Method),
			Headers: Helpers.HeadersToMap(r.Header),
			Body:    string(Helpers.GetRequestBody(r)),
		},
		SendAt: time.Now().UnixMilli(),
	}

	err := conn.WriteJSON(command)

	if err != nil {
		http.Error(w, "Failed to send request to tunnel", http.StatusInternalServerError)
		return
	}

	PendingRequests[strconv.FormatInt(eventId, 10)] = pkg.PendingHttpRequest{
		EventId:  eventId,
		Request:  r,
		Response: w,
	}

	// keep the connection alive (TODO: implement a proper keep-alive mechanism)
	var i int = 0
	for {
		if i > 1000 {
			break
		}
		i++
		time.Sleep(100 * time.Millisecond)
		if _, exists := PendingRequests[strconv.FormatInt(eventId, 10)]; !exists {
			break
		}

	}
}
