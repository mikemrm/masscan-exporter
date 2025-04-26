package collector

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/adhocore/gronx"
	"github.com/mikemrm/masscan-exporter/internal/masscan"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
)

type Collector struct {
	logger zerolog.Logger

	name        string
	schedule    string
	scanOnStart bool
	startDelay  time.Duration
	masscan     *masscan.Masscan
	timeout     time.Duration

	mu sync.RWMutex

	collecting bool
	start      time.Time
	cache      []prometheus.Metric
	nextCache  []prometheus.Metric

	doneCh chan struct{}
}

func (c *Collector) run() error {
	logger := c.logger.With().Str("schedule", c.schedule).Logger()

	nextTick, err := gronx.NextTick(c.schedule, false)
	if err != nil {
		logger.Err(err).Msg("Error calculating next tick")

		c.Stop()

		return err
	}

	go func() {
		if c.scanOnStart {
			nextTick = time.Now()

			if c.startDelay > 0 {
				nextTick = time.Now().Add(c.startDelay)
			}
		}

		logger.Info().
			Msgf("First scan at %s (%s)", nextTick.Format(time.RFC3339), time.Until(nextTick))

		for {
			select {
			case <-time.After(time.Until(nextTick)):
			case <-c.doneCh:
				return
			}

			c.refresh()

			for {
				nextTick, err = gronx.NextTick(c.schedule, false)
				if err == nil {
					break
				}

				logger.Err(err).Msg("Error calculating next tick")

				select {
				case <-time.After(time.Minute):
				case <-c.doneCh:
					return
				}
			}

			logger.Debug().Msgf("Next scan scheduled for %s (%s)", nextTick.Format(time.RFC3339), time.Until(nextTick))
		}
	}()

	return nil
}

func (c *Collector) getNextTick() (time.Time, error) {
	for {
		nextTick, err := gronx.NextTick(c.schedule, false)
		if err == nil {
			return nextTick, nil
		}

		c.logger.Err(err).Msg("error calculating next tick")

		select {
		case <-time.After(time.Minute):
		case <-c.doneCh:
			return time.Time{}, err
		}
	}
}

func (c *Collector) Stop() {
	close(c.doneCh)
}

func (c *Collector) refresh() {
	c.mu.Lock()

	c.collecting = true

	c.mu.Unlock()

	start := time.Now()

	defer func() {
		c.mu.Lock()
		defer c.mu.Unlock()

		c.collecting = false
		c.start = start
		c.cache = c.nextCache
		c.nextCache = nil
	}()

	result := c.doCollection()
	duration := time.Since(start)

	c.addMetric(descScrapeSuccess, prometheus.GaugeValue, result, c.name)
	c.addMetric(descScrapeStart, prometheus.CounterValue, float64(start.UnixNano())/float64(time.Second), c.name)
	c.addMetric(descScrapeSeconds, prometheus.GaugeValue, float64(duration)/float64(time.Second), c.name)
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, metric := range c.cache {
		ch <- metric
	}

	var inProgress float64

	if c.collecting {
		inProgress = 1
	}

	if metric := c.buildMetric(descScrapeInProgress, prometheus.GaugeValue, inProgress, c.name); metric != nil {
		ch <- metric
	}
}

func (c *Collector) doCollection() float64 {
	c.logger.Info().Msg("collection started")

	start := time.Now()

	defer func() {
		c.logger.Info().Msgf("finished collecting in %s", time.Since(start))
	}()

	ctx := context.Background()

	if c.timeout > 0 {
		var cancel func()

		ctx, cancel = context.WithTimeout(context.Background(), c.timeout)

		defer cancel()
	}

	report, err := c.masscan.Run(ctx)
	if err != nil {
		c.logger.Err(err).Msg("failed to execute masscan")

		return 0
	}

	for ip, results := range report.Results {
		for _, port := range results.Ports {
			var value float64 = 0

			if port.Status == "open" {
				value = 1
			}

			c.addMetric(descPortsOpen, prometheus.GaugeValue, value,
				c.name, ip, strconv.Itoa(port.Port), port.Proto, port.Reason,
			)
		}
	}

	return 1
}

func (c *Collector) buildMetric(desc *prometheus.Desc, valueType prometheus.ValueType, value float64, labelValues ...string) prometheus.Metric {
	metric, err := prometheus.NewConstMetric(desc, valueType, value, labelValues...)
	if err != nil {
		c.logger.Err(err).Msg("failed creating prometheus metric")

		return nil
	}

	return metric
}

func (c *Collector) addMetric(desc *prometheus.Desc, valueType prometheus.ValueType, value float64, labelValues ...string) {
	metric := c.buildMetric(desc, valueType, value, labelValues...)

	if metric == nil {
		return
	}

	c.nextCache = append(c.nextCache, metric)
}

func NewCollector(ctx context.Context, opts ...Option) (*Collector, error) {
	cfg := newConfig(opts...)

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	logger := zerolog.Ctx(ctx).With().
		Str("component", "collector").
		Str("collector", cfg.Name).
		Logger()

	ctx = logger.WithContext(ctx)

	masscan, err := masscan.New(ctx, masscan.WithConfig(cfg.Masscan))
	if err != nil {
		return nil, err
	}

	collector := &Collector{
		logger: logger,

		name:        cfg.Name,
		schedule:    cfg.Schedule,
		scanOnStart: cfg.ScanOnStart,
		startDelay:  cfg.StartDelay,
		masscan:     masscan,
		timeout:     cfg.Timeout,

		doneCh: make(chan struct{}),
	}

	if err := collector.run(); err != nil {
		return nil, err
	}

	return collector, nil
}
