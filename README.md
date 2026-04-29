# Distributed Web Crawler

A scalable distributed web crawler built in Go using Redis, PostgreSQL, and Docker Compose.
The system uses multiple worker nodes that coordinate through a shared priority queue, respect crawl politeness rules, and store structured crawl data.

---

## 🚀 Features

* Distributed multi-worker crawling
* Redis-backed priority queue
* URL normalization for relative links
* Duplicate URL detection
* Per-domain rate limiting
* robots.txt compliance
* PostgreSQL storage for pages and discovered links
* Docker Compose horizontal scaling
* CI pipeline for formatting and build checks

---

## 🧱 Architecture

```text
                +------------------+
                |    Scheduler     |
                | seed URLs loader |
                +--------+---------+
                         |
                         v
                +------------------+
                |   Redis Queue    |
                | priority URLs    |
                +--------+---------+
                         |
          +--------------+--------------+
          |              |              |
          v              v              v
   +-------------+ +-------------+ +-------------+
   |  Worker 1   | |  Worker 2   | |  Worker N   |
   +------+------+ +------+------+ +------+------+
          |               |               |
          +-------+-------+-------+-------+
                  |               |
                  v               v
         +----------------+   +------------------+
         |     Redis      |   |   PostgreSQL     |
         | visited / rate |   | pages / links    |
         | robots cache   |   | crawl metadata   |
         +----------------+   +------------------+
```

---

## ⚙️ Tech Stack

* **Language:** Go
* **Queue / Coordination:** Redis
* **Database:** PostgreSQL
* **Containerization:** Docker + Docker Compose
* **CI:** GitHub Actions

---

## 📂 Project Structure

```text
cmd/
  scheduler/   -> pushes seed URLs into queue
  worker/      -> distributed crawler workers

internal/
  config/      -> env config loader
  fetcher/     -> HTTP fetching
  parser/      -> title + link extraction
  queue/       -> Redis priority queue
  limiter/     -> per-domain rate limiting
  robots/      -> robots.txt compliance
  storage/     -> PostgreSQL persistence
```

---

## ▶️ How To Run

### 1. Start Services + Workers

```bash
docker compose up --build --scale worker=3 -d
```

### 2. Load Seed URLs

```bash
go run cmd/scheduler/main.go
```

### 3. Watch Logs

```bash
docker compose logs -f worker
```

---

## 📈 Example Output

```text
[worker-1] Processing: https://golang.org
[worker-2] Processing: https://news.ycombinator.com
[worker-3] Rate limited for domain: news.ycombinator.com

[worker-1] Title: The Go Programming Language
[worker-2] Found Links: 226
[worker-3] Queue empty, waiting...
```

---

## 🗄️ Database Schema

### pages

Stores crawled page metadata:

* url
* domain
* title
* status_code
* crawled_at

### links

Stores discovered relationships:

* source_url
* target_url
* created_at

---

## 🧠 Key Design Decisions

### Shared Redis Queue

Allows multiple workers to consume tasks concurrently.

### Duplicate Detection

Prevents repeated crawling of already visited URLs.

### Distributed Rate Limiting

Ensures multiple workers do not overload the same domain.

### robots.txt Support

Respects common web crawling rules.

### Horizontal Scaling

Increase workers instantly:

```bash
docker compose up --scale worker=10 -d
```

---

## 🧪 CI Checks

GitHub Actions pipeline validates:

* `gofmt`
* `go test ./...`
* successful build

---

## 🔮 Future Plans

### Intelligent Crawling

* Better URL scoring and crawl prioritization
* Domain-specific crawl strategies
* Freshness-based re-crawling schedules

### Reliability Improvements

* Retry failed requests with exponential backoff
* Worker health checks and heartbeat monitoring
* Graceful shutdown with in-progress job recovery

### Discovery Enhancements

* Sitemap.xml parsing
* RSS/Atom feed discovery
* Canonical URL handling

### Data & Search

* Content hashing for duplicate page detection
* Full-text search indexing
* Export crawled data to JSON / CSV

### Observability

* Metrics dashboard
* Crawl throughput monitoring
* Error rate and latency tracking

### Scaling

* Kubernetes deployment support
* Auto-scaling workers based on queue size
* Multi-region crawling architecture

---

## 📌 Status

Production-style MVP complete.
Designed to be extended into a larger-scale crawler platform.

---

## 👨‍💻 Author

Built as a distributed systems / backend engineering project in Go.
