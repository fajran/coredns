package dropv6

import (
	"context"
	"strings"

	"github.com/coredns/coredns/plugin"

	"github.com/miekg/dns"
)

type DropV6 struct {
	Next plugin.Handler

	Suffixes []string
}

func (p DropV6) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	// TODO check for matching suffix
	matchingSuffix := ""
found:
	for _, suffix := range p.Suffixes {
		for _, question := range r.Question {
			if strings.HasSuffix(question.Name, suffix) {
				matchingSuffix = suffix
				break found
			}
		}
	}

	if matchingSuffix == "" {
		return plugin.NextOrFailure(p.Name(), p.Next, ctx, w, r)
	}

	dw := &Dropper{ResponseWriter: w, Suffix: matchingSuffix}
	return plugin.NextOrFailure(p.Name(), p.Next, ctx, dw, r)
}

func (p DropV6) Name() string {
	return "dropv6"
}

type Dropper struct {
	dns.ResponseWriter

	Suffix string
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
	droppedQueriesCount.WithLabelValues(w.Suffix).Inc()
	droppedAnswersCount.WithLabelValues(w.Suffix).Add(float64(n))

	return w.ResponseWriter.WriteMsg(nr)
}

func (w *Dropper) Write(buf []byte) (int, error) {
	return w.ResponseWriter.Write(buf)
}
