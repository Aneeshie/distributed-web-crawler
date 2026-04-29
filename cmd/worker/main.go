package main

import (
	"context"
	"fmt"
	"log"
	"test/internal/config"
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

		//TODO
		// fetch page
		// parse html
		// save data
		//enqueue discovered links

		time.Sleep(1 * time.Second)

	}

}
