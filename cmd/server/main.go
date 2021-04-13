package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("error: %v", err)
	}
}

func run() error {
	var socket string

	flag.StringVar(&socket, "socket", "ping.sock", "Ping (receiver) socket")
	flag.Parse()

	listener, err := net.Listen("unix", socket)
	if err != nil {
		return err
	}

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "pong")
	})

	var server http.Server

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalCh

		if err := server.Shutdown(context.Background()); err != nil {
			log.Printf("shutdown: %v", err)
		}
	}()

	if err := server.Serve(listener); err != http.ErrServerClosed {
		return err
	}

	return nil
}
