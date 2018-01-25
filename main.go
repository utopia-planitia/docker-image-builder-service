package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
  "net/url"
  "./pkg/proxy"
	"context"
	"os"
  "os/signal"
  "syscall"
)


func main() {

	port := flag.String("port", ":80", "default server port, ':80', ':8080'...")
	target := flag.String("target", "http://127.0.0.1:8080", "default redirect url, 'http://127.0.0.1:8080'")
	parallelism := flag.Int64("parallelism", 4, "default parallelism is '4'")

	flag.Parse()

	log.Printf("server will run on: %s\n", *port)
	log.Printf("redirecting to: %s\n", *target)
	log.Printf("parallelism of: %v\n", *parallelism)

  url, err := url.Parse(*target)
  if err != nil {
    log.Fatalf("failed to parse target: %s\n", err)
  }
  reverseProxy := httputil.NewSingleHostReverseProxy(url)

	queuedProxy := proxy.New(reverseProxy, *parallelism)

	mux := http.NewServeMux()
	mux.HandleFunc("/", queuedProxy.Handle)

  server := http.Server{
    Addr:    *port,
    Handler: mux,
  }

  stop := make(chan os.Signal, 2)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
  go func() {
    <-stop
    err = server.Shutdown(context.Background())
    if err != nil {
      log.Fatalf("server shutdown failed: %s", err)
    }
  }()

  err = server.ListenAndServe();
  if err != nil && err != http.ErrServerClosed {
    log.Fatal(err)
  }

}
