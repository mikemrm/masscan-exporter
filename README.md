# masscan-exporter

Provides a simple prometheus exporter for masscan.

This processes the scans asynchronously (not when the /metrics endpoint is requested).
This is due to the time it can take for scans to complete.

Scan times are configured with a cron style expression supporting 5, 6 and 7 segment formats.
See [here](https://github.com/adhocore/gronx/blob/main/README.md#cron-expression) for more details.

Import the grafana dashboard with id `23344`.

```
$ curl localhost:9187/metrics
# HELP masscan_collectors_total Reports the number of configured collectors.
# TYPE masscan_collectors_total gauge
masscan_collectors_total 2
# HELP masscan_ports_open Masscan port status report
# TYPE masscan_ports_open gauge
masscan_ports_open{collector="network0",ip="10.0.0.1",port="179",proto="tcp",reason="syn-ack"} 1
masscan_ports_open{collector="network0",ip="10.0.0.1",port="443",proto="tcp",reason="syn-ack"} 1
masscan_ports_open{collector="network0",ip="10.0.0.1",port="80",proto="tcp",reason="syn-ack"} 1
masscan_ports_open{collector="network0",ip="10.0.0.123",port="80",proto="tcp",reason="syn-ack"} 1
masscan_ports_open{collector="network0",ip="10.0.0.219",port="443",proto="tcp",reason="syn-ack"} 1
masscan_ports_open{collector="network0",ip="10.0.0.219",port="80",proto="tcp",reason="syn-ack"} 1
masscan_ports_open{collector="network0",ip="10.0.0.28",port="161",proto="tcp",reason="syn-ack"} 1
masscan_ports_open{collector="network0",ip="10.0.0.5",port="443",proto="tcp",reason="syn-ack"} 1
masscan_ports_open{collector="network0",ip="10.0.0.5",port="80",proto="tcp",reason="syn-ack"} 1
masscan_ports_open{collector="network0",ip="10.0.0.6",port="161",proto="tcp",reason="syn-ack"} 1
masscan_ports_open{collector="network1",ip="10.1.0.1",port="179",proto="tcp",reason="syn-ack"} 1
masscan_ports_open{collector="network1",ip="10.1.0.1",port="443",proto="tcp",reason="syn-ack"} 1
masscan_ports_open{collector="network1",ip="10.1.0.1",port="80",proto="tcp",reason="syn-ack"} 1
masscan_ports_open{collector="network1",ip="10.1.0.28",port="80",proto="tcp",reason="syn-ack"} 1
# HELP masscan_scrape_collector_success Reports if the scrape was successful.
# TYPE masscan_scrape_collector_success gauge
masscan_scrape_collector_success{collector="network0"} 1
masscan_scrape_collector_success{collector="network1"} 1
# HELP masscan_scrape_in_progress Reports if a scrape is in progress.
# TYPE masscan_scrape_in_progress gauge
masscan_scrape_in_progress{collector="network0"} 0
masscan_scrape_in_progress{collector="network1"} 0
# HELP masscan_scrape_next_start_time Reports the start time for the next scrape.
# TYPE masscan_scrape_next_start_time gauge
masscan_scrape_next_start_time{collector="network0"} 1.7456961000000699e+09
masscan_scrape_next_start_time{collector="network1"} 1.7456961300006979e+09
# HELP masscan_scrape_seconds Reports how long a scrape took in seconds.
# TYPE masscan_scrape_seconds gauge
masscan_scrape_seconds{collector="network0"} 67.674926113
masscan_scrape_seconds{collector="network1"} 66.65523849
# HELP masscan_scrape_start_time Reports the start time of the scrape.
# TYPE masscan_scrape_start_time counter
masscan_scrape_start_time{collector="network0"} 1.7456958000000699e+09
masscan_scrape_start_time{collector="network1"} 1.7456958300006979e+09
# HELP masscan_scrapes_failed_current The number of consecutive scrapes which have failed.
# TYPE masscan_scrapes_failed_current gauge
masscan_scrapes_failed_current{collector="network0"} 0
masscan_scrapes_failed_current{collector="network1"} 0
# HELP masscan_scrapes_total Total number of scrapes executed for the collector.
# TYPE masscan_scrapes_total counter
masscan_scrapes_total{collector="network0",result="failed"} 0
masscan_scrapes_total{collector="network0",result="success"} 3
masscan_scrapes_total{collector="network1",result="failed"} 0
masscan_scrapes_total{collector="network1",result="success"} 3
```

### Example Config:

```yaml
loglevel: info # default: info
collectors:
  - name: network0
    schedule: '*/5 * * * *'
    masscan:
      ranges:
        - 10.0.0.0/24
      ports:
        - 80
        - 443
        - 100-200
  - name: network1
    schedule: '30 */5 * * * * *'
    scan_on_start: true
    start_delay: 10s
    timeout: 10m
    masscan:
      max_rate: 500
      ranges:
        - 10.1.0.0/24
      config: |
        ports = 80,443,100-200
  - name: network3
    schedule: '@hourly'
    scan_on_start: true
    start_delay: 10s
    timeout: 10m
    masscan:
      ranges:
        url: https://net.example.com/ranges?zone=external
        url_config:
          auth:
            bearer: file:///run/secrets/net-token
      ports: https://net.example.com/ports
# - name: collector-name          # required
#   schedule: '30 */5 * * * * *'  # required
#   scan_on_start: false          # scans on start
#   start_delay: 0s               # delays scan on start
#   timeout: 0s                   # sets a timeout for a scan (default: disabled)
#   masscan:                      # masscan config
#     temp_dir: /tmp              # temp directory for masscan runs
#     bin_path: /usr/bin/masscan  # path to masscan
#     wait_delay: 20s             # delay to wait for command to exit when cancelled
#     max_rate: 100               # masscan scan rate
#     ranges: []                  # ip ranges (overrides config ranges) (dynamic value, see below)
#     ports: []                   # port ranges (overrides config ports) (dynamic value, see below)
#     config_path: ""             # path to an existing masscan config (overrides config option)
#     config: ""                  # provide a masscan config as a string (overrides config_source) (dynamic value, see below)
server:
  listen: :9187 # default: :9187
  # The number of times a collector can fail before /readyz will report unhealthy.
  # default is 5, set to 0 to disable.
  unhealthy_failed_scrapes: 5
```

### Dynamic Value Configuration

Dynamic fields support a number of configuration methods.

1. You may statically set the raw value, for example:
   ```yaml
   ranges: [10.1.0.0/24, 10.2.0.0/24]
   config: |
     ports = 80,443
   ```
2. Dynamically using the env/file/http(s) prefixes, for example:
   ```yaml
   ranges: https://net.example.com/ranges
   ports: file:///data/ports
   config: env://MASSCAN_CONFIG
   ```
3. Or define the whole dynamic config directly.
   ```yaml
   ranges:
     url: https://net.example.com
     url_config:
       auth:
         bearer: file:///path/to/token
   ports:
     file: /data/ports
   config:
     env: MASSCAN_CONFIG
   ```

For fields which are a list of values, values may be in the form of JSON, comma separated or newline separate responses.
The following examples all provide the same results.

- `["10.1.0.0/24", "10.2.0.0/24"]`
- `10.1.0.0/24,10.2.0.0/24`
- ```
  10.1.0.0/24
  10.2.0.0/24
  ```

See below for all available dynamic configuration options:

```yaml
field_name:
  value: ...       # static value
  env: ""          # environment variable to source value from
  file: ""         # direct file path to source
  url: ""          # remote file path (http/https). If no scheme is provided, https is assumed.
  url_config:
    method: ""     # http method (default: GET)
    auth:
      username: "" # provide basic auth username
      password: "" # provide basic auth password
      bearer: ""   # provide a bearer token (overrides basic auth if not empty)
    headers: {}    # string key/value headers
    body: ""       # body contents to send with non GET requests.
```

All string values may be prefixed with either `env://` or `file://` to source the value from environment variables or a file.
These values are loaded on each masscan run, so a value may be changed on the fly.

## Development

In addition to [`go`], some `make` commands use [`docker`] and [`jq`].

```shell
$ make help
help    Show this help.
all     Tests and builds the binary and container images.
test    Runs go tests
build   Builds a binary for the current os/arch.
image   Builds all container images.
```

***note**: `make image` builds containers for all platforms.
Ensure your buildx environment is configured to support amd64 and arm64 platforms.*

[`go`]: https://go.dev
[`docker`]: https://docker.com
[`jq`]: https://jqlang.org
