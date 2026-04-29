run

```bash
    docker compose up --build --scale worker=3 -d
```


For checking worker logs
```bash
     docker compose logs -f worker
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
