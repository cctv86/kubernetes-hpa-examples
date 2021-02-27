package prometheus

import (
	"strconv"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	reqsName        = "http_requests_total"
	latencyName     = "http_request_duration_seconds"
	connectionsName = "tcp_connections_total"
)

// Prometheus is a handler that exposes prometheus metrics for the number of requests,
// the latency and the response size, partitioned by status code, method and HTTP path.
//
// Usage: pass its `ServeHTTP` to a route or globally.
type Prometheus struct {
	reqs        *prometheus.CounterVec
	latency     *prometheus.HistogramVec
	connections *prometheus.GaugeVec
}

// New returns a new prometheus middleware.
func New(name string) *Prometheus {
	p := Prometheus{}
	p.reqs = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:        reqsName,
			Help:        "How many HTTP requests processed, partitioned by status code, method and HTTP path.",
			ConstLabels: prometheus.Labels{"test": name},
		},
		[]string{"code", "method", "path"},
	)
	prometheus.MustRegister(p.reqs)

	p.latency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        latencyName,
		Help:        "How long it took to process the request, partitioned by status code, method and HTTP path.",
		ConstLabels: prometheus.Labels{"test": name},
		//Buckets:     buckets,
	},
		[]string{"code", "method", "path"},
	)
	prometheus.MustRegister(p.latency)

	p.connections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        connectionsName,
			Help:        "How many TCP connections, partitioned by  method and HTTP path.",
			ConstLabels: prometheus.Labels{"test": name},
		},
		[]string{},
	)
	prometheus.MustRegister(p.connections)
	return &p
}

func (p *Prometheus) ServeHTTP(ctx iris.Context) {
	start := time.Now()
	ctx.Next()
	r := ctx.Request()
	statusCode := strconv.Itoa(ctx.GetStatusCode())

	p.reqs.WithLabelValues(statusCode, r.Method, r.URL.Path).
		Inc()

	p.latency.WithLabelValues(statusCode, r.Method, r.URL.Path).
		Observe(float64(time.Since(start).Nanoseconds()) / 1000000000)

	p.connections.WithLabelValues().
		Set(getConnections())
}
