package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"regexp"
	"syscall"

	"github.com/damoon/docker-image-builder-service/pkg/proxy"
)

func main() {

	port := flag.String("port", ":80", "default server port, ':80', ':8080'...")
	target := flag.String("target", "http://127.0.0.1:8080", "default redirect url, 'http://127.0.0.1:8080'")
	maxCPU := flag.Int64("maxCPU", 4000, "maximum cpu milliseconds used per second")
	maxMemory := flag.Int64("maxMemory", 17179869184, "maximum memory used in byte")
	cpu := flag.Int64("cpu", 1000, "cpu milliseconds used per second per build")
	memory := flag.Int64("memory", 4294967296, "memory used per build")

	flag.Parse()

	parallelism := calculateParallelism(*maxCPU, *maxMemory, *cpu, *memory)

	log.Printf("server will run on: %s\n", *port)
	log.Printf("redirecting to: %s\n", *target)
	log.Printf("allow a maximum of %v cpu\n", *maxCPU)
	log.Printf("allow a maximum of %v bytes memory\n", *maxMemory)
	log.Printf("use %v millisecounds of cpu per build\n", *cpu)
	log.Printf("use %v bytes memory per build\n", *memory)
	log.Printf("parallelism of: %v\n", *parallelism)

	url, err := url.Parse(*target)
	if err != nil {
		log.Fatalf("failed to parse target: %s\n", err)
	}

	reg, err := regexp.Compile("^/[^/]*/build")
	if err != nil {
		log.Fatalf("failed to prepare pattern matching: %s\n", err)
	}

	server := buildServer(url, parallelism, cpu, memory, port, reg)

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

func calculateParallelism(maxCPU, maxMemory, cpu, memory int64) *int64 {
	parallelismByCPU := maxCPU / cpu
	parallelismByMemory := maxMemory / memory
	if parallelismByCPU > parallelismByMemory {
		return &parallelismByCPU
	}
	return &parallelismByMemory
}

func buildServer(url *url.URL, parallelism, cpu, memory *int64, addr *string, reg *regexp.Regexp) *http.Server {

	reverseProxy := httputil.NewSingleHostReverseProxy(url)

	queuedProxy := proxy.New(reverseProxy, *parallelism, *cpu, *memory, reg)

	mux := http.NewServeMux()
	mux.HandleFunc("/", queuedProxy.Handle)

	return &http.Server{
		Addr:    *addr,
		Handler: mux,
	}
}
