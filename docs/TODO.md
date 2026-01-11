# VPS Go (FX + Chi) Template — TODO

Source of truth: `architecture/vps-go-fx-template.md`.

## Status Legend
- [ ] Not started
- [~] In progress
- [x] Done

## Milestone 1: Bootable long-lived server

- [ ] Add `cmd/server/main.go` entrypoint (single FX app; listens until shutdown).
- [ ] Add core infra modules (config + logger; defaults; zero required env vars).
- [ ] Add `internal/router` module (Chi mux + middleware + FX grouped route registration).
- [ ] Add `internal/server` module (`http.Server` with `fx.Lifecycle` start/stop + timeouts).
- [ ] Add `health` domain module (`GET /health` -> `{ "ok": true }`).
- [ ] Validate locally:
  - [ ] `go test ./...`
  - [ ] `go build ./...`
  - [ ] `go run ./cmd/server` then `curl http://localhost:${APP_PORT:-8080}/health`

## Milestone 2: Production workflow artifacts

- [ ] Add `Makefile` targets:
  - [ ] `make start` (`go run ./cmd/server`)
  - [ ] `make build/prod` (compose build)
  - [ ] `make start/prod` (compose up)
  - [ ] `make push/prod` (compose push)
- [ ] Add `Dockerfile` (multi-stage build; single server binary runtime image).
- [ ] Add `docker-compose.prod.yaml` (image name/tag via `IMAGE`/`TAG`; env passthrough).
- [ ] Add `.env.example` (document optional env vars).
- [ ] Update `README.md` (exact run locally + prod commands; optional env vars).

## Milestone 3: Optional integrations

- [ ] Add optional Postgres wiring (disabled unless `DB_HOST` + `DB_NAME` set; must not block startup).
- [ ] Add optional Redis wiring (disabled when `REDIS_URL` empty; must not block startup).
- [ ] Add optional Inngest endpoint module:
  - [ ] Route: `POST /api/inngest` (or `/inngest`)
  - [ ] Missing keys should return `501/503` (no startup failure)
  - [ ] Minimal example function registration (so endpoint demonstrates something)

## Milestone 4: Template instantiation/adoption (optional)

- [ ] Support `gonew` flow (module path + tag strategy documented).
- [ ] Add non-destructive `cmd/adopt` CLI to write guidance files only (no runtime code changes).
