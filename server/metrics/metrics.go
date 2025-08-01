package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/tun-io/tun-io/server/ws"
	"time"
)

var (
	currentConnectionsGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "tunio_current_connections",
		Help: "Current number of active WebSocket connections",
	})

	requestsPerSubdomain = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tunio_requests_per_subdomain",
			Help: "Number of requests per subdomain",
		},
		[]string{"subdomain"},
	)
)

func NewSubdomainRequest(subdomain string) {
	requestsPerSubdomain.WithLabelValues(subdomain).Inc()
}

func updateCurrentConnections() {
	for {
		currentConnections := len(ws.Connections) // Assuming Connections is a map of active connections
		currentConnectionsGauge.Set(float64(currentConnections))

		time.Sleep(10 * time.Second) // Update every 10 seconds
	}
}

func init() {
	go updateCurrentConnections()

}
