package collector

import "github.com/prometheus/client_golang/prometheus"

var (
	descScrapeSuccess    = prometheus.NewDesc("masscan_scrape_collector_success", "Reports if the scrape was successful.", []string{"collector"}, nil)
	descScrapeStart      = prometheus.NewDesc("masscan_scrape_start_time", "Reports the start time of the scrape.", []string{"collector"}, nil)
	descScrapeSeconds    = prometheus.NewDesc("masscan_scrape_seconds", "Reports how long a scrape took in seconds.", []string{"collector"}, nil)
	descScrapeInProgress = prometheus.NewDesc("masscan_scrape_in_progress", "Reports if a scrape is in progress.", []string{"collector"}, nil)
	descPortsOpen        = prometheus.NewDesc("masscan_ports_open", "Masscan port status report", []string{"collector", "ip", "port", "proto", "reason"}, nil)
)

func Describe(ch chan<- *prometheus.Desc) {
	ch <- descScrapeSuccess
	ch <- descScrapeStart
	ch <- descScrapeSeconds
	ch <- descScrapeInProgress
	ch <- descPortsOpen
}
