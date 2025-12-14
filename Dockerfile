FROM debian:13-slim AS build

ARG MASSCAN_REPO=https://github.com/robertdavidgraham/masscan
# renovate: datasource=github-releases depName=masscan packageName=robertdavidgraham/masscan
ARG MASSCAN_VERSION=1.3.2

ARG MASSCAN_BUILD_ATTEMPTS=5

RUN apt update \
    && DEBIAN_FRONTEND="noninteractive" apt install -y git gcc make libpcap-dev \
    && git clone "${MASSCAN_REPO}" \
    && cd masscan \
    && git checkout "${MASSCAN_VERSION}" \
    && for attempt in $(seq 0 $((MASSCAN_BUILD_ATTEMPTS-1))); do \
        echo "Build attempt $((attempt+1))/$MASSCAN_BUILD_ATTEMPTS"; \
        make && break; \
    done && \
    make test

FROM debian:13-slim

ARG TARGETPLATFORM

RUN apt update \
    && DEBIAN_FRONTEND="noninteractive" apt install -y libpcap0.8 \
    && rm -rf /var/lib/apt/lists/*

COPY --from=build /masscan/bin/masscan /usr/bin/

COPY $TARGETPLATFORM/masscan-exporter /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/masscan-exporter"]
