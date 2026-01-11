# vps-go-fx-template

VPS-friendly Go service skeleton using:
- Uber FX for dependency injection
- Chi for HTTP routing

## Run locally

```bash
make start
curl http://localhost:${APP_PORT:-8080}/health
```

No env vars are required; defaults are in `config/config.go`.

Optional Postgres wiring is enabled when `DB_HOST` and `DB_NAME` are set (with `DB_USER`/`DB_PASSWORD` as needed).

## Production (Docker)

```bash
cp .env.example .env
make build/prod
make start/prod
```

Image name/tag can be customized via `IMAGE`/`TAG` when running the compose commands.
