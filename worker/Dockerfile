FROM golang:1.10.0-alpine3.7 AS compile
WORKDIR /go/src/github.com/utopia-planitia/docker-image-builder/
RUN apk --no-cache add git
COPY . worker
RUN CGO_ENABLED=0 GOOS=linux go install -a -installsuffix cgo ./worker

FROM alpine:3.7
RUN apk --no-cache add ca-certificates
RUN apk --no-cache add bash
WORKDIR /
COPY --from=compile /go/bin/worker /usr/bin/worker
CMD ["/usr/bin/worker"]
