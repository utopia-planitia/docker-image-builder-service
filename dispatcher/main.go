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
)

func main() {

	var err error

	addr := flag.String("address", ":2375", "default server address, ':2375'")
	cpu := flag.Int64("cpu", 0, "cpu microseconds used per second per build")
	memory := flag.Int64("memory", 0, "memory used per build in Byte")
	network := flag.String("network", "", "network to use for build")
	target := flag.String("workers", "", "worker urls, 'http://worker_1:2375,...'")

	flag.Parse()

	targets := strings.Split(*target, ",")
	endpoints := make([]*url.URL, len(targets))
	for i, e := range targets {
		endpoints[i], err = url.Parse(e)
		if err != nil {
			log.Fatalf("failed to parse target: %s\n", err)
		}
	}
	if len(targets) == 0 {
		log.Fatalln("no endpoints provided")
	}

	log.Printf("server will run on: %s\n", *addr)
	log.Printf("use %v microseconds of cpu per build\n", *cpu)
	log.Printf("use %v bytes memory per build\n", *memory)
	log.Printf("use %v as workers\n", *target)

	server := newDispatcher(endpoints, cpu, memory, network, addr)

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
