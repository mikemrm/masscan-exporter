package exporter

import (
	"context"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func (e *exporter) collectMetrics(ch chan<- prometheus.Metric) float64 {
	e.logger.Info().Msg("collection started")

	start := time.Now()

	defer func() {
		e.logger.Info().Msgf("finished collecting in %s", time.Since(start))
	}()

	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)

	defer cancel()

	report, err := e.masscan.Run(ctx)
	if err != nil {
		e.logger.Err(err).Msg("failed to execute masscan")

		return 0
	}

	for ip, results := range report.Results {
		for _, port := range results.Ports {
			var value float64 = 0

			if port.Status == "open" {
				value = 1
			}

			e.addMetric(ch, e.descPortsOpen, prometheus.GaugeValue, value,
				ip, strconv.Itoa(port.Port), port.Proto, port.Reason,
			)
		}
	}

	return 1
}

func (e *exporter) addMetric(ch chan<- prometheus.Metric, desc *prometheus.Desc, valueType prometheus.ValueType, value float64, labelValues ...string) {
	metric, err := prometheus.NewConstMetric(desc, valueType, value, labelValues...)
	if err != nil {
		e.logger.Err(err).Msg("failed creating prometheus metric")

		return
	}

	ch <- metric
}
