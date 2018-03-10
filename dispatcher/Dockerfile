FROM golang:1.10.0-alpine3.7 AS compile
WORKDIR /go/src/github.com/utopia-planitia/docker-image-builder/
COPY . dispatcher
RUN CGO_ENABLED=0 GOOS=linux go install -a -installsuffix cgo ./dispatcher

FROM alpine:3.7
RUN apk --no-cache add ca-certificates
WORKDIR /
COPY --from=compile /go/bin/dispatcher /usr/bin/dispatcher
ENTRYPOINT ["/usr/bin/dispatcher"]
