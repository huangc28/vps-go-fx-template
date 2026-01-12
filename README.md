# VPS Go (FX + Chi) Boilerplate

A VPS-friendly Go service skeleton using:
- Uber FX for dependency injection
- chi for HTTP routing
- zap for logging (use `*zap.SugaredLogger` in app code)

## Quick start (this repo)

```bash
cp .env.example .env
make start
curl http://localhost:${APP_PORT:-8080}/health
```

## 1) Brand new project

### Option A (recommended): generate into an empty folder

```bash
mkdir my-service && cd my-service
go mod init github.com/<you>/my-service

GOPROXY=direct go run github.com/huangc28/vps-go-fx-template/cmd/adopt@latest --dir . --scaffold
go mod tidy

make start
curl http://localhost:8080/health
```

### Option B: clone this repo and edit the module path

```bash
git clone https://github.com/huangc28/vps-go-fx-template.git my-service
cd my-service
rm -rf .git

go mod edit -module github.com/<you>/my-service
go mod tidy

make start
```

## 2) Existing project (adopt into it)

From the root of the repo you want to adopt into:

### A) Write guidance only (no runtime code changes)

```bash
GOPROXY=direct go run github.com/huangc28/vps-go-fx-template/cmd/adopt@latest --dir .
```

Outputs:
- `AGENTS.md`
- `architecture/vps-go-fx-template.md`
- `codex/skills/adopt/SKILL.md`

### B) Write the actual scaffold (deterministic db/cache/etc)

This creates the canonical folders (`db/`, `cache/`, `internal/`, `cmd/server`, etc) and refuses to overwrite existing files unless `--force` is set.

```bash
GOPROXY=direct go run github.com/huangc28/vps-go-fx-template/cmd/adopt@latest --dir . --scaffold
go mod tidy
```

Then run:

```bash
make start
curl http://localhost:8080/health
```

Use `--force` to overwrite if the target files already exist.

For local development in this repo:

```bash
cd cmd/adopt
go run . --dir ../../some/other/repo --scaffold
```

To enable `/adopt` in Codex, copy the skill into your Codex home (commonly `~/.codex/skills/adopt`).

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
