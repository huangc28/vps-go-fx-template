# VPS Go (FX + Chi) Boilerplate

A VPS-friendly Go service skeleton using:
- Uber FX for dependency injection
- chi for HTTP routing
- zap for logging

## Quick start (run locally)

```bash
cp .env.example .env
make start
curl http://localhost:${APP_PORT:-8080}/health
```

## Start a brand new project

1) Copy this repo as your new repo (or `git clone` it).
2) Update the module path:

```bash
go mod edit -module github.com/<you>/<your-service>
go mod tidy
```

3) Run it:

```bash
make start
```

## Adopt into an existing project (non-destructive)

This writes guidance files into your existing repo without changing runtime code.

From the root of the repo you want to adopt into:

```bash
go run /path/to/vps-go-fx-template/cmd/adopt --dir .
```

Outputs:
- `AGENTS.md`
- `architecture/vps-go-fx-template.md`

Use `--force` to overwrite if those files already exist.

## Optional integrations

- Postgres: enabled when `DB_HOST` + `DB_NAME` are set (`DB_USER`/`DB_PASSWORD` as needed).
- Redis: enabled when `REDIS_HOST` is set (`REDIS_SCHEME` can be `redis` or `rediss`).

## sqlc (optional)

This repo includes a ready-to-run `sqlc.yaml` that reads:
- schema: `supabase/schema.sql`
- queries: `db/query/*.sql`
- output: `db/sqlc` (package `dbsqlc`)

Install `sqlc` and run:

```bash
make sqlc
```
