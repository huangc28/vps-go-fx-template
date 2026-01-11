# Go VPS Structure Template Plan (FX + Chi)

This document defines a conservative, reusable boilerplate/template plan for deploying Go services to a VPS (or any long-lived environment) while keeping the same core patterns from this repo: **Uber FX** for dependency injection, **Chi** for HTTP routing, and per-domain handler organization.

The key change vs Vercel: the process is **long-lived**, so we **bootstrap FX once** at startup and run an `http.Server` until shutdown.

## 0) Inngest (prebuilt, optional)

The VPS template should include an optional Inngest HTTP endpoint out of the box (same intent as the Vercel variant), implemented as a normal route on the long-lived server:

- Provide a neutral wrapper package: `internal/pkg/inngestclient`.
- Provide an HTTP handler that mounts the Inngest endpoint, e.g. `POST /api/inngest` (or `/inngest`).
- Wire it as a domain-like module: `internal/app/inngest/fx.Module` that contributes its handler via `router.AsRoute(...)`.

The template baseline should still run without Inngest configured; the Inngest handler can respond with a clear `501`/`503` until required env vars are set.

## 1) What Changes vs the Vercel Template

### 1.1 Entry points and lifecycle

Vercel template:
- Entrypoints live at `api/<domain>/core.go`.
- Each request bootstraps an FX app (per-request DI) and executes a handler.
- Routing is mediated by `vercel.json` rewrites.

VPS template:
- Entrypoints live under `cmd/` (for example `cmd/server/main.go`).
- The FX app is created **once** per process and held for the lifetime of the service.
- The HTTP server listens on a port and shuts down gracefully on `SIGINT`/`SIGTERM`.

### 1.2 Folder layout (`lib` -> `internal`)

For VPS services, it’s common to use Go’s `internal/` tree to prevent external imports and keep boundaries clear.

Note: Go convention is `internal/` (singular), not `internals/`. If you strongly prefer `internals/`, it will work, but it won’t get the same tooling/convention benefits.

### 1.3 Routing and “domains”

You can preserve the same domain pattern:
- Keep a router package with a handler interface and grouped registration (via FX).
- Keep handlers under `internal/app/<domain>/...`.
- Register handlers into the router module and mount them under the appropriate path prefix.

### 1.4 Deployment assets

Vercel uses platform-managed build/deploy.
VPS templates typically include one of:
- Dockerfile + `docker-compose.prod.yaml` (recommended for repeatable deploys)

## 2) Template Goal (What we want to extract)

Create a reusable VPS template that mirrors the architecture and ergonomics of this repo:
- FX modules for config/logger/DB/Redis.
- Chi router with handler grouping.
- Minimal example routes (`/health`).
- Optional prebuilt Inngest endpoint.
- Docker-based production build/run workflow.

Non-goals:
- No platform-specific Vercel wiring (`api/<domain>/core.go`, `vercel.json`, `vercel dev`).
- No business-domain code.

## 3) Proposed Template Structure (Conservative)

```
.
├── cmd/
│   ├── adopt/                    # optional: add Codex context to an existing repo
│   │   ├── main.go               # writes AGENTS + architecture plan + Codex skill
│   │   └── assets/               # embedded files for non-destructive adoption
│   └── server/
│       └── main.go                 # process entrypoint (long-lived)
├── config/
│   └── config.go                   # viper + typed config + defaults
├── cache/
│   └── redis.go                    # redis init (optional)
├── db/
│   ├── db.go                       # sqlx postgres init + fx lifecycle close hooks
│   └── tx.go                       # sqlx transaction helper (Tx wrapper)
├── internal/
│   ├── app/
│   │   ├── fx/
│   │   │   └── core.go             # CoreAppOptions (config/logger/db/redis)
│   │   ├── health/
│   │   │   └── handler.go          # example handler implementing router.Handler
│   │   ├── inngest/                # optional (prebuilt)
│   │   │   ├── handler.go          # mounts /api/inngest (or /inngest)
│   │   │   └── fx/
│   │   │       └── module.go       # exports fx.Module
│   │   └── <domain>/
│   │       └── ...                 # domain-specific handlers/services
│   ├── logs/
│   │   └── logs.go                 # logger constructor
│   ├── pkg/
│   │   └── inngestclient/          # optional, generic Inngest wrapper
│   │   └── render/                 # ChiJSON / ChiErr response helpers
│   ├── interfaces/
│   │   └── README.md               # cross-domain interfaces (avoid cycles)
│   ├── router/
│   │   ├── core.go                 # router.Handler + registration loop
│   │   └── fx/
│   │       └── options.go          # CoreRouterOptions (chi mux + handler group)
│   └── server/
│       ├── http.go                 # http.Server construction + lifecycle wiring
│       └── fx/
│           └── options.go          # ServerOptions (start/stop hooks)
├── docker-compose.prod.yaml        # production build/run/push
├── Dockerfile
├── Makefile
└── README.md
```

## 4) Dependency Injection + Router Pattern (VPS)

Keep the same conceptual split:

### 4.0 Domain-owned modules (recommended)

Each domain owns its own FX wiring and dependencies, and exposes a single module that `cmd/server/main.go` aggregates. This keeps the entrypoint minimal while preserving clear ownership boundaries.

- Each domain exports `internal/app/<domain>/fx.Module` as an `fx.Option` (typically `fx.Options(...)`).
- The domain module:
  - Provides domain services/repos/controllers via `fx.Provide(...)`.
  - Registers HTTP handlers via `router.AsRoute(New<Handler>)` so the generic router module can collect them (FX group injection).
- Cross-domain dependencies must be expressed via interfaces in `internal/interfaces/<area>` packages to avoid cyclic imports.

### 4.1 Core app module (infra)
- `config.NewViper` and `config.NewConfig`
- `internal/logs.NewLogger`
- `db.NewSQLXPostgresDB` (optional, enabled when `DB_HOST` + `DB_NAME` are present)
- `cache.NewRedis` (optional, enabled when `REDIS_HOST` is present)

### 4.2 Router module
- `internal/router.Handler` interface:
  - `RegisterRoute(r *chi.Mux)`
  - `Handle(w http.ResponseWriter, r *http.Request)`
- `internal/router.AsRoute(...)` to add handlers to an FX group (e.g. `group:"handlers"`).
- `internal/router/fx.CoreRouterOptions` builds a `*chi.Mux`, attaches middleware, and calls `RegisterRoute` on all grouped handlers.

### 4.3 HTTP server module
VPS adds an explicit server module that:
- Builds an `http.Server` with sensible timeouts.
- Starts listening on `APP_PORT` in an `fx.Lifecycle` `OnStart` hook.
- Shuts down gracefully in `OnStop`.

The server module depends on the constructed `*chi.Mux` and config/logger.

## 5) Makefile + Production Workflow

Proposed Makefile targets (as requested):

- `make start`
  - Runs the server locally (example: `go run ./cmd/server`).

- `make build/prod`
  - Builds the production app image via `docker-compose.prod.yaml`.

- `make start/prod`
  - Builds and starts the production image via `docker-compose.prod.yaml`.

- `make push/prod`
  - Pushes the production image to a remote image registry via `docker-compose.prod.yaml`.

Notes:
- `docker-compose.prod.yaml` should parameterize the image name/tag via env vars (e.g. `IMAGE`, `TAG`) so CI/CD and local flows match.
- If you want “build once, deploy many”, also consider emitting the image digest and pinning it on the VPS.

## 6) Minimal Example Domain (`health`)

Include a minimal health handler to validate wiring:
- Route: `GET /health`
- Response: `200` with `{ "ok": true }` (using `internal/pkg/render.ChiJSON`)

This handler is registered into the router via `router.AsRoute(NewHealthHandler)`.

## 6.1 Optional Inngest endpoint

Include an optional Inngest handler module to provide a “ready-to-wire” endpoint:
- Route: `POST /api/inngest` (or `/inngest`)
- Response helpers: `internal/pkg/render.ChiJSON` / `internal/pkg/render.ChiErr`
- Gating: if required env vars are not set, return a clear `501`/`503` instead of failing app startup

## 7) How to Add a New Domain

1) Create handler(s) under `internal/app/<domain>/...` implementing `internal/router.Handler`.
2) Create `internal/app/<domain>/fx/module.go` exporting `Module fx.Option` and register constructors there via `router.AsRoute(...)`.
3) Add `internal/app/<domain>/fx.Module` to the aggregate app in `cmd/server/main.go`.
4) Implement `RegisterRoute` to mount routes (e.g. `r.Route("/v2/<domain>", ...)`).

No `vercel.json` rewrites are needed; routing is entirely inside the running server.

## 8) Env vars (example baseline)

Keep the same config approach (Viper + defaults). Typical vars:
- `APP_NAME` (default: `go-vps-service`)
- `APP_ENV` (default: `development`)
- `APP_PORT` (default: `8080`)
- `LOG_LEVEL` (default: `info`)
- Postgres (optional; disabled unless `DB_HOST` + `DB_NAME` are set):
  - `DB_USER`
  - `DB_PASSWORD`
  - `DB_HOST`
  - `DB_PORT` (default: `5432`)
  - `DB_NAME`
- Redis (optional; disabled unless `REDIS_HOST` is set):
  - `REDIS_USER`
  - `REDIS_PASSWORD`
  - `REDIS_HOST`
  - `REDIS_PORT` (default: `6379`)
  - `REDIS_SCHEME` (default: `redis`, use `rediss` for TLS)

Minimal to spin up locally: none (use defaults), optionally set `APP_PORT`.

If you enable Inngest, required env vars depend on the chosen Inngest Go SDK wiring. Commonly:
- `INNGEST_EVENT_KEY`
- `INNGEST_SIGNING_KEY`
And optionally an app identifier such as `INNGEST_APP_ID` (depending on how you register/label functions).

## 9) sqlc (optional, manual)

Keep sqlc conventions identical if desired:
- Schema source at `supabase/schema.sql`
- Queries under `db/query/*.sql`
- `sqlc.yaml` defines the generated package location

The template should compile even if generated code is not present (don’t hard-require sqlc output at baseline).

## 10) Validation Checklist

- `go test ./...`
- `go build ./...`
- Local run smoke test:
  - `make start`
  - `curl http://localhost:$APP_PORT/health`
- Production image:
  - `make build/prod`
  - `make start/prod`

## 11) “Low Friction” Dev Experience (Day-1)

To make it easy for a developer (or an AI agent) to implement and use the boilerplate with minimal setup:

- The project must start with **zero required env vars** (sensible defaults).
- Postgres/Redis must be **optional** and disabled unless `DB_HOST` + `DB_NAME` are set / `REDIS_HOST` is set.
- Provide an `.env.example` and document the exact minimal commands to run locally.
- Include a small, non-destructive “adopter” CLI (`cmd/adopt`) so an existing repo can add Codex guidance without reorganizing code.

Recommended “first run” workflow:

1) `make start`
2) `curl http://localhost:${APP_PORT:-8080}/health`

If you include Inngest:

3) Start Inngest dev server pointing at your endpoint (example):
   - `npx inngest-cli@latest dev -u http://localhost:${APP_PORT:-8080}/api/inngest --no-discovery`

## 12) Implementation Checklist (for an AI agent)

This section is intentionally explicit so an agent can scaffold the repo without guessing.

### 12.1 Minimal code artifacts

- `cmd/server/main.go`
  - Builds `fx.New(...)` with: core app options, router options, server options, and domain modules (eg `health`, optional `inngest`).
  - Starts the FX app and blocks on signals (or uses `fx.App.Run()`).

- `internal/server/http.go` + `internal/server/fx/options.go`
  - Provide `*http.Server` and register `fx.Lifecycle` hooks to listen/shutdown.
  - Depend on `*chi.Mux` from the router module.

- `db/tx.go`
  - Provide a small `sqlx` transaction helper (e.g. `db.Tx(db, func(tx *sqlx.Tx) {...})`) to standardize commit/rollback behavior.

- `internal/router/core.go` + `internal/router/fx/options.go`
  - Keep `Handler` interface + FX group pattern (mirrors this repo).
  - Build `*chi.Mux`, attach baseline middleware, then call `RegisterRoute` for each grouped handler.

- `internal/app/health/handler.go` (+ optional `internal/app/health/fx/module.go`)
  - `GET /health` returning `{ "ok": true }` via `internal/pkg/render`.

- `internal/app/inngest/handler.go` + `internal/app/inngest/fx/module.go` (optional but “prebuilt”)
  - Mount `POST /api/inngest` (or `/inngest`) and serve via `inngestgo.Client.Serve().ServeHTTP`.
  - Register at least one example function (cron) so the endpoint demonstrates registration.
  - Gate behavior when keys are missing (respond `501`/`503`), but do not prevent the app from starting.

### 12.2 Minimal non-code artifacts

- `Makefile`
  - `start`: `go run ./cmd/server`
  - `build/prod`, `start/prod`, `push/prod` using `docker-compose.prod.yaml`

- `Dockerfile`
  - Multi-stage build (builder -> runtime) producing a single server binary.

- `docker-compose.prod.yaml`
  - A service for the app image with env var passthrough.
  - Image name/tag driven by env vars (eg `IMAGE`, `TAG`) to reduce CI friction.

- `README.md`
  - Exact “run locally” steps and required env vars (none required; list optional ones).

### 12.3 Minimal env vars to “spin it up”

None required (use defaults). The only commonly-set variable for local dev should be:
- `APP_PORT` (optional)

Optional integrations:
- Postgres: `DB_HOST`, `DB_NAME` (plus `DB_USER`/`DB_PASSWORD` as needed)
- Redis: `REDIS_HOST` (plus `REDIS_USER`/`REDIS_PASSWORD` as needed)
- Inngest (for real use): `INNGEST_EVENT_KEY`, `INNGEST_SIGNING_KEY` (and optional `INNGEST_SIGNING_KEY_FALLBACK`), plus `INNGEST_APP_ID` (optional/labeling)

## 13) Implement Using `gonew` (Adopt / Instantiate)

To match the Vercel template’s onboarding ergonomics, the VPS template should also support a “one command” instantiation flow using `gonew`.

### 13.1 Create and publish the template as a Go module

1) Create a new repo for the VPS template (example name: `go-vps-fx-template`).
2) Ensure `go.mod` module path matches the repo (example: `module github.com/<org>/go-vps-fx-template`).
3) Push a tag (`v0.x.y`) so `gonew` can fetch a stable version.

### 13.2 Install `gonew`

```bash
go install golang.org/x/tools/cmd/gonew@latest
```

### 13.3 Instantiate a new service

```bash
gonew github.com/<org>/go-vps-fx-template github.com/<org>/<service>
```

This produces a new folder (named after the destination module by default) with:
- `cmd/server/main.go` as the long-lived entrypoint
- `internal/` packages wired via FX + Chi
- `Makefile`, `Dockerfile`, and `docker-compose.prod.yaml`

### 13.4 Adopt into an existing repo (non-destructive)

If you only want to add Codex/agent guidance files (no code changes), include the same “adopter” CLI pattern as the Vercel template and run it from inside the target repo:

```bash
go run github.com/<org>/go-vps-fx-template/cmd/adopt@latest --dir .
```

This should write (without touching runtime code):
- `AGENTS.md`
- `architecture/go-vps-reusable-template-plan.md`
- `codex/skills/adopt/SKILL.md`

To enable `/adopt` in Codex, install the skill to your Codex home (commonly `~/.codex/skills/adopt`).

Recommended workflow (keeps history and avoids risky in-place rewrites):

1) Instantiate into a sibling directory via `gonew` (as above).
2) Copy your existing domain code into `internal/app/<domain>/...` (or wrap it there).
3) Expose your domain(s) via `internal/app/<domain>/fx.Module` and register HTTP handlers via `router.AsRoute(...)`.
4) Add those domain modules to `cmd/server/main.go`.

This mirrors the Vercel template’s domain/module pattern, but the app boots once at process start instead of per request.

### 13.5 After `gonew` (manual edits you still do)

- Update `README.md` service name, ports, and deploy notes.
- Set image naming defaults in `docker-compose.prod.yaml` (e.g. `IMAGE`, `TAG`) for your registry.
- Confirm optional infra behavior:
- If `DB_HOST` or `DB_NAME` is empty, DB wiring should no-op or return a clear disabled error at call sites (but must not block startup).
  - If `REDIS_HOST` is empty, Redis wiring should be disabled (but must not block startup).
- If you enable Inngest, confirm missing keys gate requests with `501/503` rather than failing startup.

### 13.6 Private template notes (if applicable)

If the template repo is private, ensure the developer has Git credentials configured so `gonew` can fetch it, and prefer pinning to a tagged version for repeatability:

```bash
gonew github.com/<org>/go-vps-fx-template@v0.1.0 github.com/<org>/<service>
```
