package server

import (
	"github.com/charmbracelet/log"
	_ "github.com/gorilla/websocket"
	"github.com/tun-io/tun-io/internal/http/Helpers"
	"github.com/tun-io/tun-io/server/ws"
	"net/http"
	"strconv"
)

func StartServer() {
	var port = 8080

	mux := http.NewServeMux()
	mux.HandleFunc("/_tunio/_internal/api/client/ws", ws.WsUpgradeRoute)

	// create a subdomain handler
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		log.Infof("New HTTP request received: %s %s (requester: %s)", r.Method, r.URL.Path, r.RemoteAddr)

		subdomain := Helpers.GetSubdomain(r)
		if subdomain == "" {
			http.Error(w, "Subdomain not found", http.StatusBadRequest)
			log.Warnf("Subdomain not found in request: %s (requester: %s)", r.Host, r.RemoteAddr)
			return
		}

		_, subdomainExists := ws.Connections[subdomain]
		if !subdomainExists {
			http.Error(w, "Target not found", http.StatusNotFound)
			log.Warnf("Subdomain %s not found in connections", subdomain)
			return
		}

		ws.SendTunnelRequest(subdomain, r, w)
	})

	http.ListenAndServe(":"+strconv.Itoa(port), mux)
}
