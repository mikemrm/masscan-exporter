# masscan-exporter

Provides a simple prometheus exporter for masscan.

```
$ curl localhost:9187/metrics
# HELP masscan_ports_open Masscan port status report
# TYPE masscan_ports_open gauge
masscan_ports_open{ip="10.0.0.1",port="179",proto="tcp",reason="syn-ack",scan_job="some-job"} 1
masscan_ports_open{ip="10.0.0.1",port="443",proto="tcp",reason="syn-ack",scan_job="some-job"} 1
masscan_ports_open{ip="10.0.0.1",port="80",proto="tcp",reason="syn-ack",scan_job="some-job"} 1
masscan_ports_open{ip="10.0.0.123",port="80",proto="tcp",reason="syn-ack",scan_job="some-job"} 1
masscan_ports_open{ip="10.0.0.219",port="443",proto="tcp",reason="syn-ack",scan_job="some-job"} 1
masscan_ports_open{ip="10.0.0.219",port="80",proto="tcp",reason="syn-ack",scan_job="some-job"} 1
masscan_ports_open{ip="10.0.0.28",port="161",proto="tcp",reason="syn-ack",scan_job="some-job"} 1
masscan_ports_open{ip="10.0.0.5",port="443",proto="tcp",reason="syn-ack",scan_job="some-job"} 1
masscan_ports_open{ip="10.0.0.5",port="80",proto="tcp",reason="syn-ack",scan_job="some-job"} 1
masscan_ports_open{ip="10.0.0.6",port="161",proto="tcp",reason="syn-ack",scan_job="some-job"} 1
# HELP masscan_scrape_collector_success Reports if the scrape was successful.
# TYPE masscan_scrape_collector_success gauge
masscan_scrape_collector_success{scan_job="some-job"} 1
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
  labels:
    scan_job: some-job
```
