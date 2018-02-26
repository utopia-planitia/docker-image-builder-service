
builder: $(shell find ./cmd/builder -type f)
#	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o builder ./cmd/builder
	go get -t ./cmd/builder/...
	go build -o builder ./cmd/builder

dispatcher: $(shell find ./cmd/dispatcher -type f)
#	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o dispatcher ./cmd/dispatcher
	go get -t ./cmd/builder/...
	go build -o dispatcher ./cmd/dispatcher

clean:
	unlink builder
	unlink dispatcher
