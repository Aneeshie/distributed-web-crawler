package limiter

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Limiter struct {
	client *redis.Client
	ctx    context.Context
}

func New(client *redis.Client) *Limiter {
	return &Limiter{
		client: client,
		ctx:    context.Background(),
	}
}

func (l *Limiter) Allow(domain string) (bool, error) {
	key := fmt.Sprintf("rate:%s", domain)

	status, err := l.client.SetArgs(
		l.ctx,
		key,
		"1",
		redis.SetArgs{
			Mode: "NX",
			TTL:  time.Second,
		},
	).Result()
	// With NX, Redis returns a nil reply if the key already exists (i.e., throttled).
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return status == "OK", nil
}
