package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Config struct {
	RMQ *Pool
}

func main() {
	log.Println("Starting data-processor")
	ctx := context.Background()
	if err := run(ctx, os.Stdout, os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, w io.Writer, args []string) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	var app Config

	pool, err := NewPool(5, 5)
	if err != nil {
		return err
	}
	defer pool.Close()

	app.RMQ = pool

	var wg sync.WaitGroup
	stopChan := make(chan struct{})
	wg.Add(1)
	go app.worker(&wg, stopChan)

	select {}
}

func (app *Config) worker(wg *sync.WaitGroup, stopChan <-chan struct{}) {
	defer wg.Done()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:

			log.Println("checking messages")
			app.RMQ.Consume(func(url string) {
				log.Printf("Received message: %s", url)
			})

		case <-stopChan:
			log.Println("Stopping worker")
			return
		}
	}
}
