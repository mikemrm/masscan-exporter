package exporter

import (
	"fmt"
	"time"

	"github.com/mikemrm/masscan-exporter/internal/masscan"
	"github.com/prometheus/client_golang/prometheus"
)

type exporter struct {
	timeout time.Duration
	masscan *masscan.Masscan

	descScrapeSuccess *prometheus.Desc
	descPortsOpen     *prometheus.Desc
}

func (c *exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.descScrapeSuccess
	ch <- c.descPortsOpen
}

func (c *exporter) Collect(ch chan<- prometheus.Metric) {
	addMetric(ch, c.descScrapeSuccess, prometheus.GaugeValue, c.collect(ch))
}

func New(masscan *masscan.Masscan, opts ...Option) (prometheus.Collector, error) {
	cfg := newConfig(opts...)

	exporter := &exporter{
		timeout: cfg.Timeout,
		masscan: masscan,

		descScrapeSuccess: prometheus.NewDesc("masscan_scrape_collector_success", "Reports if the scrape was successful.", nil, cfg.Labels),
		descPortsOpen:     prometheus.NewDesc("masscan_ports_open", "Masscan port status report", []string{"ip", "port", "proto", "reason"}, cfg.Labels),
	}

	if err := cfg.Registerer.Register(exporter); err != nil {
		return nil, fmt.Errorf("cannot register the exporter: %w", err)
	}

	return exporter, nil
}
