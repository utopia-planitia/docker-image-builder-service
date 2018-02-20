FROM minio/mc:RELEASE.2018-02-09T23-07-36Z as minio-client

FROM docker:17.11.0-ce-dind as docker-client

FROM golang:1.9.3-alpine3.7 AS compile
WORKDIR /go/src/github.com/damoon/docker-image-builder-service/
COPY dibs dibs
COPY cmd/builder cmd/builder
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o builder ./cmd/builder

FROM alpine:3.7
RUN apk --no-cache add ca-certificates
WORKDIR /
COPY --from=minio-client /usr/bin/mc /usr/bin/mc
COPY --from=docker-client /usr/local/bin/docker /usr/bin/docker
COPY --from=compile /go/src/github.com/damoon/docker-image-builder-service/builder /usr/bin/builder
COPY caching-scripts/save.sh /usr/bin/save
COPY caching-scripts/load.sh /usr/bin/load
ENV CACHE_MAX_AGE 864000
COPY caching-scripts/deleteOld.sh /usr/bin/deleteOld
COPY caching-scripts/entrypoint.sh /entrypoint.sh
ENTRYPOINT [ "/entrypoint.sh" ]
CMD ["/usr/bin/builder"]
