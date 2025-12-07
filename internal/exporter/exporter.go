package exporter

import (
	"context"
	"fmt"

	"github.com/mikemrm/masscan-exporter/internal/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
)

var (
	descCollectors = prometheus.NewDesc("masscan_collectors_total", "Reports the number of configured collectors.", nil, nil)
)

type exporter struct {
	logger     *zerolog.Logger
	collectors []*collector.Collector
}

func (e *exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- descCollectors

	collector.Describe(ch)
}

func (e *exporter) Collect(ch chan<- prometheus.Metric) {
	totalCollectors, err := prometheus.NewConstMetric(descCollectors, prometheus.GaugeValue, float64(len(e.collectors)))
	if err != nil {
		e.logger.Err(err).Msg("failed to create total collectors metric")
	} else {
		ch <- totalCollectors
	}

	for _, c := range e.collectors {
		c.Collect(ch)
	}
}

func New(ctx context.Context, opts ...Option) (prometheus.Collector, error) {
	cfg := newConfig(opts...)

	exporter := &exporter{
		logger:     zerolog.Ctx(ctx),
		collectors: cfg.Collectors,
	}

	if err := cfg.Registerer.Register(exporter); err != nil {
		return nil, fmt.Errorf("cannot register the exporter: %w", err)
	}

	return exporter, nil
}
