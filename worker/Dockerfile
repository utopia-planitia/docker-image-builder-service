FROM golang:1.14.4-alpine3.12 AS compile

RUN apk --no-cache add git
RUN mkdir /worker
WORKDIR /worker

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . /worker
RUN CGO_ENABLED=0 GOOS=linux go install -a -installsuffix cgo .


FROM scratch
COPY --from=compile /go/bin/worker /worker
ENTRYPOINT [ "/worker"]
