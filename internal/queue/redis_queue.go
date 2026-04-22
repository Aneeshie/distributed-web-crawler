package queue

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
)

const QUEUE_KEY = "queue_crawler"

type Queue struct {
	client *redis.Client
	ctx    context.Context
}

func NewQueue(client *redis.Client) *Queue {
	return &Queue{
		client: client,
		ctx:    context.Background(),
	}
}

// push url with score

func (q *Queue) Enqueue(url string, score float64) error {
	return q.client.ZAdd(q.ctx, QUEUE_KEY, redis.Z{
		Score:  score,
		Member: url,
	}).Err()
}

// pop the highest and return it

func (q *Queue) PopHighest() (string, error) {
	items, err := q.client.ZPopMax(q.ctx, QUEUE_KEY, 1).Result()
	if err != nil {
		return "", err
	}

	if len(items) == 0 {
		return "", errors.New("Queue is empty")
	}

	url, ok := items[0].Member.(string)
	if !ok {
		return "", errors.New("invalid queue member type")
	}

	return url, nil
}

// returns the number of queued urls

func (q *Queue) QueueSize() (int64, error) {
	return q.client.ZCard(q.ctx, QUEUE_KEY).Result()
}
