package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"test/internal/config"
	"test/internal/fetcher"
	"test/internal/parser"
	"test/internal/queue"
	"test/internal/storage"

	"github.com/redis/go-redis/v9"
)

func getWorkerID() string {
	host, err := os.Hostname()
	if err != nil || host == "" {
		return "worker-unknown"
	}
	return host
}

func logWorker(workerID string, format string, args ...any) {
	prefix := "[" + workerID + "] "
	fmt.Printf(prefix+format+"\n", args...)
}

func main() {
	workerID := getWorkerID()

	cfg := config.Load()

	logWorker(workerID, "Loaded config")

	// PostgreSQL init
	db := storage.NewPostgresDB(cfg.PostgresURL)
	logWorker(workerID, "Connected to PostgreSQL")

	// Redis init
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})

	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		log.Fatal("Redis connection failed:", err)
	}

	logWorker(workerID, "Connected to Redis")

	q := queue.NewQueue(rdb)

	for {
		url, err := q.PopHighest()
		if err != nil {
			logWorker(workerID, "Queue empty, waiting...")
			time.Sleep(3 * time.Second)
			continue
		}

		logWorker(workerID, "Processing: %s", url)

		status, body, err := fetcher.Fetch(url)
		if err != nil {
			logWorker(workerID, "Fetch failed: %v", err)
			continue
		}

		logWorker(workerID, "Status: %d", status)
		logWorker(workerID, "Body bytes: %d", len(body))

		title, links, err := parser.Parse(url, body)
		if err != nil {
			logWorker(workerID, "Parse failed: %v", err)
			continue
		}

		logWorker(workerID, "Title: %s", title)
		logWorker(workerID, "Found Links: %d", len(links))

		err = storage.SavePage(db, url, title, status)
		if err != nil {
			logWorker(workerID, "Failed to save page: %v", err)
		}

		err = storage.SaveLinks(db, url, links)
		if err != nil {
			logWorker(workerID, "Failed to save links: %v", err)
		}

		time.Sleep(1 * time.Second)
	}
}
