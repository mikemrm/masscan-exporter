package exporter

import (
	"context"
	"log"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

func (e *exporter) collect(ch chan<- prometheus.Metric) float64 {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)

	defer cancel()

	report, err := e.masscan.Run(ctx)
	if err != nil {
		log.Printf("failed to execute masscan: %s", err.Error())

		return 0
	}

	for ip, results := range report.Results {
		for _, port := range results.Ports {
			var value float64 = 0

			if port.Status == "open" {
				value = 1
			}

			addMetric(ch, e.descPortsOpen, prometheus.GaugeValue, value,
				ip, strconv.Itoa(port.Port), port.Proto, port.Reason,
			)
		}
	}

	return 1
}

func addMetric(ch chan<- prometheus.Metric, desc *prometheus.Desc, valueType prometheus.ValueType, value float64, labelValues ...string) {
	metric, err := prometheus.NewConstMetric(desc, valueType, value, labelValues...)
	if err != nil {
		log.Printf("failed creating prometheus metric: %s", err.Error())

		return
	}

	ch <- metric
}
