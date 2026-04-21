run

```bash
    docker compose up -d
```

then for db

```bash
  migrate -path migrations -database "postgres://postgres:password123@localhost:5432/db?sslmode=disable" up
```

to verify

```bash
  docker exec -it postgres_crawler psql -U postgres -d db
```

then

```bash
  db=# \dt
```
