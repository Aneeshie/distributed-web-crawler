package queue

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
)

func setupTestQueue(t *testing.T) *Queue {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	err := client.Ping(context.Background()).Err()
	if err != nil {
		t.Fatalf("failed to connect to redis: %v", err)
	}

	client.Del(context.Background(), QUEUE_KEY)

	return NewQueue(client)
}

func TestEnqueueAndSize(t *testing.T) {
	q := setupTestQueue(t)

	err := q.Enqueue("https://ada.com", 100)
	if err != nil {
		t.Fatalf("enqueue failed: %v", err)
	}

	err = q.Enqueue("https://babbage.com", 90)
	if err != nil {
		t.Fatalf("enqueue failed: %v", err)
	}

	size, err := q.QueueSize()
	if err != nil {
		t.Fatalf("size failed: %v", err)
	}

	if size != 2 {
		t.Fatalf("expected size 2, got %d", size)
	}
}

func TestPopHighest(t *testing.T) {
	q := setupTestQueue(t)

	q.Enqueue("https://ada.com", 10)
	q.Enqueue("https://babbage.com", 90)
	q.Enqueue("https://curie.com", 100)

	url, err := q.PopHighest()
	if err != nil {
		t.Fatalf("pop failed: %v", err)
	}

	if url != "https://curie.com" {
		t.Fatalf("expected curie.com first, got %s", url)
	}
}

func TestEmptyQueue(t *testing.T) {
	q := setupTestQueue(t)

	_, err := q.PopHighest()
	if err == nil {
		t.Fatal("expected error on empty queue")
	}
}
