package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"sync"
	"time"
	"url-frontier/data"
)

type Config struct {
	DB     *sql.DB
	Models data.Models
	RMQ    *Pool
}

func main() {
	log.Println("starting url frontier service")
	var app Config

	// connect postgres
	dbConn := app.connectToDB(5)
	if dbConn == nil {
		log.Panic("could not connect to database")
	}

	// connect rabbitmq
	pool, err := NewPool(5, 5)
	if err != nil {
		log.Fatalf("Failed to create pool: %v", err)
	}
	defer pool.Close()

	// set config
	app.DB = dbConn
	app.Models = data.New(dbConn)
	app.RMQ = pool

	var wg sync.WaitGroup
	stopChan := make(chan struct{})
	wg.Add(1)
	go app.worker(&wg, stopChan)

	// Block the main function forever
	select {}
}

func (app *Config) worker(wg *sync.WaitGroup, stopChan chan struct{}) {
	defer wg.Done()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:

			// fetch configured number of data from db and publish to rmq
			seeds, _ := app.Models.Seed.GetAll(1000)

			// publish seeds
			app.publishSeeds(seeds)

			log.Printf("Total %d Messages published in this session. %s", len(seeds), getEmoji(len(seeds)))

		case <-stopChan:
			fmt.Println("Worker stopped")
			return
		}
	}
}

func (app *Config) publishSeeds(seeds []*data.Seed) {
	for _, seed := range seeds {
		jsonPayload, err := json.Marshal(seed)
		if err != nil {
			return
		}

		// publish to rmq
		app.RMQ.Publish(
			"amqp.direct", // exchange
			"seeds",       // routing key
			false,         // mandatory
			false,         // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        jsonPayload,
			})
	}
}

func getEmoji(count int) string {
	if count > 800 {
		return "\U0001F973"
	} else if count > 600 {
		return "\U0001F60E"
	} else if count > 400 {
		return "\U0001F916"
	} else if count > 200 {
		return "\U0001FAE1"
	} else {
		return "\U0001FAE4"
	}
}
