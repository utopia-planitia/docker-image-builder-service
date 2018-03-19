package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	var err error

	addr := flag.String("address", ":2375", "default server address, ':2375'")
	target := flag.String("docker", "http://docker_1:2375", "docker url, 'http://docker_1:2375'")

	flag.Parse()

	log.Printf("server will run on: %s\n", *addr)
	log.Printf("use %v as docker\n", *target)

	endpoint, err := url.Parse(*target)
	if err != nil {
		log.Fatalf("failed to parse target: %s\n", err)
	}

	server := newBuilder(endpoint, addr)

	// wait for an exit signal
	stop := make(chan os.Signal, 2)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-stop
		err = server.Shutdown(context.Background())
		if err != nil {
			log.Fatalf("server shutdown failed: %s\n", err)
		}
	}()

	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
