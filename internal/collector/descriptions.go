package collector

import "github.com/prometheus/client_golang/prometheus"

var (
	descScrapeSuccess    = prometheus.NewDesc("masscan_scrape_collector_success", "Reports if the scrape was successful.", []string{"collector"}, nil)
	descScrapeStart      = prometheus.NewDesc("masscan_scrape_start_time", "Reports the start time of the scrape.", []string{"collector"}, nil)
	descScrapeNextStart  = prometheus.NewDesc("masscan_scrape_next_start_time", "Reports the start time for the next scrape.", []string{"collector"}, nil)
	descScrapeSeconds    = prometheus.NewDesc("masscan_scrape_seconds", "Reports how long a scrape took in seconds.", []string{"collector"}, nil)
	descScrapeInProgress = prometheus.NewDesc("masscan_scrape_in_progress", "Reports if a scrape is in progress.", []string{"collector"}, nil)
	descScrapesTotal     = prometheus.NewDesc("masscan_scrapes_total", "Total number of scrapes executed for the collector.", []string{"collector", "result"}, nil)
	descScrapesFailed    = prometheus.NewDesc("masscan_scrapes_failed_current", "The number of consecutive scrapes which have failed.", []string{"collector"}, nil)
	descPortsOpen        = prometheus.NewDesc("masscan_ports_open", "Masscan port status report", []string{"collector", "ip", "port", "proto", "reason"}, nil)
)

func Describe(ch chan<- *prometheus.Desc) {
	ch <- descScrapeSuccess
	ch <- descScrapeStart
	ch <- descScrapeNextStart
	ch <- descScrapeSeconds
	ch <- descScrapeInProgress
	ch <- descScrapesTotal
	ch <- descScrapesFailed
	ch <- descPortsOpen
}
