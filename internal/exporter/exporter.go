package exporter

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/mikemrm/masscan-exporter/internal/masscan"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
)

type exporter struct {
	logger zerolog.Logger

	timeout time.Duration
	masscan *masscan.Masscan

	descScrapeSuccess    *prometheus.Desc
	descScrapeStart      *prometheus.Desc
	descScrapeSeconds    *prometheus.Desc
	descScrapeInProgress *prometheus.Desc
	descPortsOpen        *prometheus.Desc

	collecting      bool
	collectingStart time.Time
	collectingMu    sync.RWMutex
	cachedMetrics   []prometheus.Metric

	cacheTTL time.Duration
}

func (e *exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.descScrapeSuccess
	ch <- e.descScrapeStart
	ch <- e.descScrapeSeconds
	ch <- e.descScrapeInProgress
	ch <- e.descPortsOpen
}

func (e *exporter) Collect(ch chan<- prometheus.Metric) {
	if !e.useCache() {
		e.collect()
	}

	e.collectCache(ch)
}

func (e *exporter) collect() {
	e.collectingMu.Lock()

	e.collecting = true
	start := time.Now()

	e.collectingMu.Unlock()

	go func() {
		cacheCh := make(chan prometheus.Metric)
		doneCh := make(chan struct{})

		var cachedMetrics []prometheus.Metric

		defer func() {
			e.collectingMu.Lock()
			defer e.collectingMu.Unlock()

			<-doneCh

			e.collecting = false
			e.collectingStart = start
			e.cachedMetrics = cachedMetrics
		}()

		go func() {
			defer close(doneCh)

			for metric := range cacheCh {
				cachedMetrics = append(cachedMetrics, metric)
			}
		}()

		defer close(cacheCh)

		result := e.collectMetrics(cacheCh)
		duration := time.Since(start)

		e.addMetric(cacheCh, e.descScrapeSuccess, prometheus.GaugeValue, result)
		e.addMetric(cacheCh, e.descScrapeStart, prometheus.CounterValue, float64(start.Unix()))
		e.addMetric(cacheCh, e.descScrapeSeconds, prometheus.GaugeValue, float64(duration)/float64(time.Second))
	}()
}

func (e *exporter) useCache() bool {
	e.collectingMu.RLock()
	defer e.collectingMu.RUnlock()

	return e.collecting || !e.collectingStart.IsZero() && e.cacheTTL > 0 && !e.collectingStart.Add(e.cacheTTL).Before(time.Now())
}

func (e *exporter) collectCache(ch chan<- prometheus.Metric) {
	e.collectingMu.RLock()
	defer e.collectingMu.RUnlock()

	for _, metric := range e.cachedMetrics {
		ch <- metric
	}

	var inProgress float64

	if e.collecting {
		inProgress = 1
	}

	e.addMetric(ch, e.descScrapeInProgress, prometheus.GaugeValue, inProgress)
}

func New(ctx context.Context, masscan *masscan.Masscan, opts ...Option) (prometheus.Collector, error) {
	cfg := newConfig(opts...)

	exporter := &exporter{
		logger:  *zerolog.Ctx(ctx),
		timeout: cfg.Timeout,
		masscan: masscan,

		descScrapeSuccess:    prometheus.NewDesc("masscan_scrape_collector_success", "Reports if the scrape was successful.", nil, cfg.Labels),
		descScrapeStart:      prometheus.NewDesc("masscan_scrape_start_time", "Reports the start time of the scrape.", nil, cfg.Labels),
		descScrapeSeconds:    prometheus.NewDesc("masscan_scrape_seconds", "Reports how long a scrape took in seconds.", nil, cfg.Labels),
		descScrapeInProgress: prometheus.NewDesc("masscan_scrape_in_progress", "Reports if a scrape is in progress.", nil, cfg.Labels),
		descPortsOpen:        prometheus.NewDesc("masscan_ports_open", "Masscan port status report", []string{"ip", "port", "proto", "reason"}, cfg.Labels),

		cacheTTL: cfg.CacheTTL,
	}

	if err := cfg.Registerer.Register(exporter); err != nil {
		return nil, fmt.Errorf("cannot register the exporter: %w", err)
	}

	return exporter, nil
}
