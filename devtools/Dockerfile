FROM docker:17.12.1-ce-dind as docker-client

FROM ubuntu:16.04
COPY --from=docker-client /usr/local/bin/docker /usr/bin/docker
RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates bats curl \
    && rm -rf /var/lib/apt/lists/*
