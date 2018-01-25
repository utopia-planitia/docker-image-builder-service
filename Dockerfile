FROM golang:1.9.3-alpine3.7 AS compile
WORKDIR /go/src/github.com/damoon/docker-image-builder-service/
RUN apk --update add git
RUN go get -d -v golang.org/x/sync/semaphore
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o proxy .

FROM alpine:3.7
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=compile /go/src/github.com/damoon/docker-image-builder-service/proxy .
ENTRYPOINT ["./proxy"]  
