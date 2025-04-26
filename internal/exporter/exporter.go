package exporter

import (
	"context"
	"fmt"

	"github.com/mikemrm/masscan-exporter/internal/collector"
	"github.com/prometheus/client_golang/prometheus"
)

type exporter struct {
	collectors []*collector.Collector
}

func (e *exporter) Describe(ch chan<- *prometheus.Desc) {
	collector.Describe(ch)
}

func (e *exporter) Collect(ch chan<- prometheus.Metric) {
	for _, c := range e.collectors {
		c.Collect(ch)
	}
}

func New(ctx context.Context, opts ...Option) (prometheus.Collector, error) {
	cfg := newConfig(opts...)

	exporter := &exporter{
		collectors: cfg.Collectors,
	}

	if err := cfg.Registerer.Register(exporter); err != nil {
		return nil, fmt.Errorf("cannot register the exporter: %w", err)
	}

	return exporter, nil
}
