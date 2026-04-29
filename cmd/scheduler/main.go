package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"test/internal/config"
	"test/internal/queue"

	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.Load()

	fmt.Println("Loaded config")

	ctx := context.Background()

	// Redis connection
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})

	err := rdb.Ping(ctx).Err()
	if err != nil {
		log.Fatal("Redis connection failed:", err)
	}

	fmt.Println("Connected to Redis")

	// Queue instance
	q := queue.NewQueue(rdb)

	// Open seeds file
	file, err := os.Open("seeds.txt")
	if err != nil {
		log.Fatal("failed to open seeds.txt:", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	loaded := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines
		if line == "" {
			continue
		}

		// Validate URL
		if !isValidURL(line) {
			fmt.Println("Skipping invalid URL:", line)
			continue
		}

		// Enqueue with default priority
		err := q.Enqueue(line, 100)
		if err != nil {
			fmt.Println("Failed to enqueue:", line)
			continue
		}

		fmt.Println("Queued:", line)
		loaded++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("error reading seeds.txt:", err)
	}

	size, err := q.QueueSize()
	if err != nil {
		log.Fatal("failed to get queue size:", err)
	}

	fmt.Printf("\nLoaded %d seed URLs\n", loaded)
	fmt.Printf("Queue size: %d\n", size)
}

func isValidURL(raw string) bool {
	u, err := url.ParseRequestURI(raw)
	if err != nil {
		return false
	}

	return u.Scheme != "" && u.Host != ""
}
