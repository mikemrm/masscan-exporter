package exporter

import (
	"context"
	"fmt"
	"time"

	"github.com/mikemrm/masscan-exporter/internal/masscan"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
)

type exporter struct {
	logger zerolog.Logger

	timeout time.Duration
	masscan *masscan.Masscan

	descScrapeSuccess *prometheus.Desc
	descScrapeSeconds *prometheus.Desc
	descPortsOpen     *prometheus.Desc
}

func (e *exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.descScrapeSuccess
	ch <- e.descScrapeSeconds
	ch <- e.descPortsOpen
}

func (e *exporter) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()
	result := e.collect(ch)
	duration := time.Since(start)

	e.addMetric(ch, e.descScrapeSuccess, prometheus.GaugeValue, result)
	e.addMetric(ch, e.descScrapeSeconds, prometheus.GaugeValue, float64(duration)/float64(time.Second))
}

func New(ctx context.Context, masscan *masscan.Masscan, opts ...Option) (prometheus.Collector, error) {
	cfg := newConfig(opts...)

	exporter := &exporter{
		logger:  *zerolog.Ctx(ctx),
		timeout: cfg.Timeout,
		masscan: masscan,

		descScrapeSuccess: prometheus.NewDesc("masscan_scrape_collector_success", "Reports if the scrape was successful.", nil, cfg.Labels),
		descScrapeSeconds: prometheus.NewDesc("masscan_scrape_seconds", "Reports how long a scrape took in seconds.", nil, cfg.Labels),
		descPortsOpen:     prometheus.NewDesc("masscan_ports_open", "Masscan port status report", []string{"ip", "port", "proto", "reason"}, cfg.Labels),
	}

	if err := cfg.Registerer.Register(exporter); err != nil {
		return nil, fmt.Errorf("cannot register the exporter: %w", err)
	}

	return exporter, nil
}
