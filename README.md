# masscan-exporter

Provides a simple prometheus exporter for masscan.

```
$ curl localhost:9187/metrics
# HELP masscan_ports_open Masscan port status report
# TYPE masscan_ports_open gauge
masscan_ports_open{ip="10.0.0.1",port="179",proto="tcp",reason="syn-ack"} 1
masscan_ports_open{ip="10.0.0.1",port="443",proto="tcp",reason="syn-ack"} 1
masscan_ports_open{ip="10.0.0.1",port="80",proto="tcp",reason="syn-ack"} 1
masscan_ports_open{ip="10.0.0.123",port="80",proto="tcp",reason="syn-ack"} 1
masscan_ports_open{ip="10.0.0.219",port="443",proto="tcp",reason="syn-ack"} 1
masscan_ports_open{ip="10.0.0.219",port="80",proto="tcp",reason="syn-ack"} 1
masscan_ports_open{ip="10.0.0.28",port="161",proto="tcp",reason="syn-ack"} 1
masscan_ports_open{ip="10.0.0.5",port="443",proto="tcp",reason="syn-ack"} 1
masscan_ports_open{ip="10.0.0.5",port="80",proto="tcp",reason="syn-ack"} 1
masscan_ports_open{ip="10.0.0.6",port="161",proto="tcp",reason="syn-ack"} 1
# HELP masscan_scrape_collector_success Reports if the scrape was successful.
# TYPE masscan_scrape_collector_success gauge
masscan_scrape_collector_success 1
# HELP masscan_scrape_in_progress Reports if a scrape is in progress.
# TYPE masscan_scrape_in_progress gauge
masscan_scrape_in_progress 1
# HELP masscan_scrape_seconds Reports how long a scrape took in seconds.
# TYPE masscan_scrape_seconds gauge
masscan_scrape_seconds 296.448288709
# HELP masscan_scrape_start_time Reports the start time of the scrape.
# TYPE masscan_scrape_start_time counter
masscan_scrape_start_time 1.745622339e+09
```

### Example Config:

```yaml
masscan:
  ranges:
    - 10.0.0.0/24
  ports:
    - 80
    - 443
    - 100-200

exporter:
  cache_ttl: 10m
```

You can also provide a masscan config.

```yaml
masscan:
  ranges:
    - 10.0.0.0/24
  config: |
    ports = 80,443,100-200

exporter:
  cache_ttl: 10m
```
