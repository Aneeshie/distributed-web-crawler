package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"test/internal/config"
	"test/internal/fetcher"
	"test/internal/limiter"
	"test/internal/parser"
	"test/internal/queue"
	"test/internal/robot"
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

const defaultEnqueueScore = 100.0

// crawlableURL reports whether raw is suitable to place on the crawl queue (http/https with host).
func crawlableURL(raw string) bool {
	u, err := url.Parse(raw)
	if err != nil || u.Host == "" {
		return false
	}
	switch strings.ToLower(u.Scheme) {
	case "http", "https":
		return true
	default:
		return false
	}
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
	rl := limiter.New(rdb)
	robotsChecker := robot.New(rdb)

	for {
		urlStr, err := q.PopHighest()
		if err != nil {
			logWorker(workerID, "Queue empty, waiting...")
			time.Sleep(3 * time.Second)
			continue
		}

		parsedURL, err := url.Parse(urlStr)
		if err != nil || parsedURL.Host == "" {
			logWorker(workerID, "Invalid URL: %s", urlStr)
			continue
		}

		domain := parsedURL.Host

		allowedByRobots, err := robotsChecker.Allowed(urlStr)
		if err != nil {
			logWorker(workerID, "robots.txt check failed: %v", err)
		}

		if !allowedByRobots {
			logWorker(workerID, "Blocked by robots.txt: %s", urlStr)
			continue
		}

		allowed, err := rl.Allow(domain)
		if err != nil {
			logWorker(workerID, "Rate limiter failed: %v", err)
			continue
		}

		if !allowed {
			logWorker(workerID, "Rate limited for domain: %s", domain)
			// Don't drop work just because we're throttled: requeue and retry later.
			if err := q.Enqueue(urlStr, 100); err != nil {
				logWorker(workerID, "Failed to requeue (rate limited): %v", err)
			}
			time.Sleep(1 * time.Second)
			continue
		}

		logWorker(workerID, "Processing: %s", urlStr)

		status, body, err := fetcher.Fetch(urlStr)
		if err != nil {
			logWorker(workerID, "Fetch failed: %v", err)
			continue
		}

		logWorker(workerID, "Status: %d", status)
		logWorker(workerID, "Body bytes: %d", len(body))

		title, links, err := parser.Parse(urlStr, body)
		if err != nil {
			logWorker(workerID, "Parse failed: %v", err)
			continue
		}

		logWorker(workerID, "Title: %s", title)
		logWorker(workerID, "Found Links: %d", len(links))

		err = storage.SavePage(db, urlStr, title, status)
		if err != nil {
			logWorker(workerID, "Failed to save page: %v", err)
		}

		err = storage.SaveLinks(db, urlStr, links)
		if err != nil {
			logWorker(workerID, "Failed to save links: %v", err)
		}

		for _, discovered := range links {
			if discovered == urlStr || !crawlableURL(discovered) {
				continue
			}
			already, err := storage.PageExists(db, discovered)
			if err != nil {
				logWorker(workerID, "PageExists check failed for %s: %v (enqueueing anyway)", discovered, err)
			} else if already {
				continue
			}
			if err := q.Enqueue(discovered, defaultEnqueueScore); err != nil {
				logWorker(workerID, "Failed to enqueue discovered URL %s: %v", discovered, err)
			}
		}

		time.Sleep(1 * time.Second)
	}
}
