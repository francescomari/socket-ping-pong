package main

import (
	"context"
	"flag"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	var dialer net.Dialer

	client := http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return dialer.DialContext(ctx, "unix", socket)
			},
		},
	}

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	for {
		if err := ping(&client); err != nil {
			log.Printf("ping: %v", err)
		}

		select {
		case <-signalCh:
			return nil
		case <-time.After(1 * time.Second):
			// Send next ping.
		}
	}
}

func ping(client *http.Client) error {
	req, err := http.NewRequest(http.MethodPost, "http://unix/ping", nil)
	if err != nil {
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Printf("close: %v", err)
		}
	}()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	log.Printf("received '%s'", string(data))

	return nil
}
