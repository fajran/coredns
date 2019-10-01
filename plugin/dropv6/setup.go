package dropv6

import (
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/metrics"
	clog "github.com/coredns/coredns/plugin/pkg/log"

	"github.com/caddyserver/caddy"
	"github.com/prometheus/client_golang/prometheus"
)

var log = clog.NewWithPlugin("dropv6")

var (
	droppedAnswersCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: plugin.Namespace,
		Subsystem: "dropv6",
		Name:      "dropped_answers_total",
		Help:      "Counter of dropped answers",
	}, []string{"prefix"})

	droppedQueriesCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: plugin.Namespace,
		Subsystem: "dropv6",
		Name:      "dropped_queries_total",
		Help:      "Counter of dropped queries",
	}, []string{"prefix"})
)

func init() {
	plugin.Register("dropv6", setup)
}

func setup(c *caddy.Controller) error {
	// TODO load prefix from config

	c.OnStartup(func() error {
		metrics.MustRegister(c, droppedAnswersCount, droppedQueriesCount)
		return nil
	})

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return DropV6{Next: next}
	})

	return nil
}
