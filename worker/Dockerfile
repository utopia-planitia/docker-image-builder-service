FROM minio/mc:RELEASE.2018-02-09T23-07-36Z as minio-client

FROM docker:17.11.0-ce-dind as docker-client

FROM golang:1.10.0-alpine3.7 AS compile
WORKDIR /go/src/github.com/utopia-planitia/docker-image-builder/
RUN apk --no-cache add git
RUN go get code.cloudfoundry.org/bytefmt
COPY . worker
RUN CGO_ENABLED=0 GOOS=linux go install -a -installsuffix cgo ./worker

FROM alpine:3.7
RUN apk --no-cache add ca-certificates
RUN apk --no-cache add bash
WORKDIR /
COPY --from=minio-client /usr/bin/mc /usr/bin/mc
COPY --from=docker-client /usr/local/bin/docker /usr/bin/docker
COPY --from=compile /go/bin/worker /usr/bin/worker
COPY caching-scripts/save.sh /usr/bin/save
COPY caching-scripts/load.sh /usr/bin/load
COPY caching-scripts/uncachedBytes.sh /usr/bin/uncachedBytes
ENV CACHE_MAX_AGE 864000
COPY caching-scripts/deleteOld.sh /usr/bin/deleteOld
COPY caching-scripts/entrypoint.sh /entrypoint.sh
ENTRYPOINT [ "/entrypoint.sh" ]
CMD ["/usr/bin/worker"]
