FROM golang:1.9.3-alpine3.7 AS compile
WORKDIR /go/src/github.com/damoon/docker-image-builder-service/
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o scheduler ./cmd/scheduler

FROM alpine:3.7
RUN apk --no-cache add ca-certificates
WORKDIR /
COPY --from=compile /go/src/github.com/damoon/docker-image-builder-service/scheduler .
ENTRYPOINT ["./scheduler"]
