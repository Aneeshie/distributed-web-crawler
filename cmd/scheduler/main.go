package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"test/internal/config"

	_ "github.com/lib/pq"

	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.Load()

	fmt.Println("Loaded config")

	ctx := context.Background()

	//redis

	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})

	err := rdb.Ping(ctx).Err()
	if err != nil {
		log.Fatal("Redis connection failed:", err)
	}
	fmt.Println("Connected to Redis")

	//psql

	//later we'll migrate to pgx ( which is modern and faster than lib/pq)

	db, err := sql.Open("postgres", cfg.PostgresURL)
	if err != nil {
		log.Fatal("Postgres open failed:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Postgres ping failed:", err)
	}

	fmt.Println("Connected to PostgreSQL")
}
