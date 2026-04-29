package robot

import (
	"bufio"
	"context"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type Checker struct {
	client *redis.Client
	ctx    context.Context
	http   *http.Client
}

func New(client *redis.Client) *Checker {
	return &Checker{
		client: client,
		ctx:    context.Background(),
		http: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *Checker) Allowed(rawURL string) (bool, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return false, err
	}

	domain := u.Scheme + "://" + u.Host
	cacheKey := "robots:" + u.Host

	// Check cached robots rules
	rules, err := c.client.Get(c.ctx, cacheKey).Result()
	if err == redis.Nil {
		rules, err = c.fetchRobots(domain)
		if err != nil {
			return true, nil // fail open for MVP
		}

		c.client.Set(c.ctx, cacheKey, rules, 1*time.Hour)
	} else if err != nil {
		return true, nil
	}

	path := u.Path
	return allowedByRules(path, rules), nil
}

func (c *Checker) fetchRobots(domain string) (string, error) {
	resp, err := c.http.Get(domain + "/robots.txt")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", nil
	}

	var lines []string
	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return strings.Join(lines, "\n"), nil
}

func allowedByRules(path, content string) bool {
	lines := strings.Split(content, "\n")

	inStarBlock := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(strings.ToLower(line), "user-agent:") {
			ua := strings.TrimSpace(strings.SplitN(line, ":", 2)[1])
			inStarBlock = ua == "*"
			continue
		}

		if inStarBlock && strings.HasPrefix(strings.ToLower(line), "disallow:") {
			rule := strings.TrimSpace(strings.SplitN(line, ":", 2)[1])

			if rule == "" {
				continue
			}

			if strings.HasPrefix(path, rule) {
				return false
			}
		}
	}

	return true
}
