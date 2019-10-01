package dropv6

import (
	"context"

	"github.com/coredns/coredns/plugin"

	"github.com/miekg/dns"
)

type DropV6 struct {
	Next plugin.Handler
}

func (p DropV6) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	// TODO check for matching prefix

	dw := &Dropper{w, ""}
	return plugin.NextOrFailure(p.Name(), p.Next, ctx, dw, r)
}

func (p DropV6) Name() string {
	return "dropv6"
}

type Dropper struct {
	dns.ResponseWriter

	Prefix string
}

func (w *Dropper) WriteMsg(r *dns.Msg) error {
	answers := make([]dns.RR, 0)
	for _, a := range r.Answer {
		if _, ipv6 := a.(*dns.AAAA); !ipv6 {
			answers = append(answers, a)
		}
	}

	if len(answers) == len(r.Answer) {
		return w.ResponseWriter.WriteMsg(r)
	}

	nr := r.Copy()
	nr.Answer = answers

	n := len(r.Answer) - len(answers)
	droppedQueriesCount.WithLabelValues(w.Prefix).Inc()
	droppedAnswersCount.WithLabelValues(w.Prefix).Add(float64(n))

	return w.ResponseWriter.WriteMsg(nr)
}

func (w *Dropper) Write(buf []byte) (int, error) {
	return w.ResponseWriter.Write(buf)
}
