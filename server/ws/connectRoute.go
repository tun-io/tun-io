package ws

import (
	"github.com/charmbracelet/log"
	"github.com/gorilla/websocket"
	"github.com/tun-io/tun-io/internal/http/Helpers"
	"github.com/tun-io/tun-io/pkg"
	"net/http"
)

var upgrader = websocket.Upgrader{} // use default options

// Connections key: string subdomain, value: *websocket.Conn
var Connections = map[string]*pkg.SyncSafeSocket{}

func WsUpgradeRoute(w http.ResponseWriter, r *http.Request) {
	subdomain := Helpers.GetSubdomain(r)
	if subdomain == "" {
		http.Error(w, "Subdomain not found", http.StatusBadRequest)
		log.Warnf("Subdomain not found in request: %s (requester: %s)", r.Host, r.RemoteAddr)
		return
	}

	var _, subdomainInUse = Connections[subdomain]

	if subdomainInUse {
		http.Error(w, "Subdomain already in use", http.StatusConflict)
		log.Warnf("Subdomain %s already in use", subdomain)
		return
	}

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("Failed to upgrade connection for subdomain %s: %v (requester: %s)", subdomain, err, r.RemoteAddr)
		return
	}

	// Clean up the connection when done
	defer func(c *websocket.Conn, subdomain string) {
		delete(Connections, subdomain)
		c.Close()
		log.Infof("Connection closed for subdomain: %s (requester: %s)", subdomain, c.RemoteAddr().String())
	}(c, subdomain)

	log.Infof("New connection established for subdomain: %s (requester: %s)", subdomain, r.RemoteAddr)
	Connections[subdomain] = pkg.NewSyncSafeSocket(c)

	messageHandler(c)
}
