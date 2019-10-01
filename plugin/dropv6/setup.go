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
	}, []string{"suffix"})

	droppedQueriesCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: plugin.Namespace,
		Subsystem: "dropv6",
		Name:      "dropped_queries_total",
		Help:      "Counter of dropped queries",
	}, []string{"suffix"})
)

func init() {
	plugin.Register("dropv6", setup)
}

func setup(c *caddy.Controller) error {
	// TODO load prefix from config

	suffixes := readSuffixes(c)
	if len(suffixes) > 0 {
		log.Info("Dropping IPv6 results from", suffixes)
	}

	c.OnStartup(func() error {
		metrics.MustRegister(c, droppedAnswersCount, droppedQueriesCount)
		return nil
	})

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return DropV6{
			Next:     next,
			Suffixes: suffixes,
		}
	})

	return nil
}

func readSuffixes(c *caddy.Controller) []string {
	suffixes := make(map[string]struct{})
	for c.Next() {
		for _, arg := range c.RemainingArgs() {
			suffixes[arg] = struct{}{}
		}
	}

	result := make([]string, 0)
	for suffix := range suffixes {
		result = append(result, suffix)
	}

	return result
}
