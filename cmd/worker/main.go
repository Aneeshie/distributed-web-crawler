package main

import (
	"context"
	"fmt"
	"log"
	"test/internal/config"
	"test/internal/fetcher"
	"test/internal/parser"
	"test/internal/queue"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.Load()

	fmt.Println("Loaded config")

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

		// parse html
		// save data
		//enqueue discovered links

		time.Sleep(1 * time.Second)

	}

}
