package main

import (
	"context"
	"fmt"
	"log"
	"test/internal/config"
	"test/internal/fetcher"
	"test/internal/parser"
	"test/internal/queue"
	"test/internal/storage"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.Load()

	fmt.Println("Loaded config")

	//db init

	db := storage.NewPostgresDB(cfg.PostgresURL)

	//connect to redis

	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})

	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		log.Fatal("Redis connection failed: ", err)
	}

	fmt.Println("Connected to Redis")

	q := queue.NewQueue(rdb)
	for {
		url, err := q.PopHighest()

		if err != nil {
			fmt.Println("Queue empty, waiting ....")
			time.Sleep(3 * time.Second)
		}

		fmt.Println("Processing: ", url)

		status, body, err := fetcher.Fetch(url)
		if err != nil {
			fmt.Println("Fetch failed: ", err)
			continue
		}
		fmt.Println("Status: ", status)
		fmt.Println("Body bytes: ", len(body))

		title, links, err := parser.Parse(body)
		if err != nil {
			fmt.Println("parse failed: ", err)
			continue
		}

		fmt.Println("Title: ", title)
		fmt.Println("Found Links: ", len(links))

		err = storage.SavePage(db, url, title, status)
		if err != nil {
			fmt.Println("Failed to save the page: ", err)
		}
		err = storage.SaveLinks(db, url, links)
		if err != nil {
			fmt.Println("Failed to save links:", err)
		}

		// parse html
		// save data
		//enqueue discovered links

		time.Sleep(1 * time.Second)

	}

}
