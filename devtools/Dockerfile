FROM docker:17.12.1-ce-dind as docker-client

FROM alpine:3.7
COPY --from=docker-client /usr/local/bin/docker /usr/bin/docker
RUN apk --no-cache add ca-certificates bats curl
WORKDIR /
