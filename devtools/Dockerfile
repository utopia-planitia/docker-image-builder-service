FROM docker:19.03.11-dind as docker-client

FROM ubuntu:20.04
COPY --from=docker-client /usr/local/bin/docker /usr/bin/docker
RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates bats curl jq dnsutils \
    && rm -rf /var/lib/apt/lists/*

COPY devtools/untilfail /usr/local/bin/untilfail
RUN chmod +x /usr/local/bin/untilfail

COPY tests /tests
