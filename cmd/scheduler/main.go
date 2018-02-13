package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/damoon/docker-image-builder-service/dibs"
)

func main() {

	var err error

	addr := flag.String("address", ":2375", "default server address, ':2375'")
	cpu := flag.Int64("cpu", 1000, "cpu milliseconds used per second per build")
	memory := flag.Int64("memory", 2147483648, "memory used per build in Byte")
	target := flag.String("workers", "http://worker_1:2375,http://worker_2:2375", "redirect urls, 'http://worker_1:2375,...'")

	flag.Parse()

	log.Printf("server will run on: %s\n", *addr)
	log.Printf("use %v millisecounds of cpu per build\n", *cpu)
	log.Printf("use %v bytes memory per build\n", *memory)

	targets := strings.Split(*target, ",")
	endpoints := make([]*url.URL, len(targets))
	for i, e := range targets {
		endpoints[i], err = url.Parse(e)
		if err != nil {
			log.Fatalf("failed to parse target: %s\n", err)
		}
	}

	server := dibs.NewScheduler(endpoints, cpu, memory, addr)

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

	// serve requests
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

}
