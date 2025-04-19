FROM debian:12-slim AS build

ARG MASSCAN_REPO=https://github.com/robertdavidgraham/masscan
ARG MASSCAN_VERSION=1.3.2

RUN apt update \
    && apt install -y git gcc make libpcap-dev \
    && git clone "${MASSCAN_REPO}" \
    && cd masscan \
    && git checkout "${MASSCAN_VERSION}" \
    && make

FROM debian:12-slim

RUN apt update && apt install -y libpcap0.8 && rm -rf /var/lib/apt/lists/*

COPY --from=build /masscan/bin/masscan /usr/bin/

COPY masscan-exporter /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/masscan-exporter"]
