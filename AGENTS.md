# AGENTS.md

## Project overview

`simon-jp-api` — REST API backend in Go (chi router, bun ORM, Postgres, Redis/Valkey, RabbitMQ). Follows a layered architecture:

- `cmd/` — entrypoints: `api`, `worker`, `scheduler`, `migrate`
- `internal/httpapi/` — HTTP handlers, routing, response DTOs
- `internal/service/` — business logic, errors
- `internal/repository/` — data access (bun queries)
- `internal/models/` — bun table models
- `internal/db/` — connection + goose migrations (`migrations/*.sql` and `*.go`)
- `internal/mq/`, `internal/worker/`, `internal/scheduler/`, `internal/config/`

## Commands

```sh
make compose-up        # start docker services (postgres, valkey, rabbitmq, adminer)
make compose-clean     # stop + wipe all volumes (destructive)
make migrate-up        # apply goose migrations (go run ./cmd/migrate, no CLI needed)
make migrate-down      # rollback last migration
make migrate-create NAME=xxx   # needs goose CLI on PATH (see below)
make run-api           # run API server (requires compose-up running)
make build             # build all binaries
make build-prod        # build static binaries
```

Verify changes with `go build ./... && go vet ./...`.

`make migrate-create` shells out to the `goose` binary — install it once with `go install github.com/pressly/goose/v3/cmd/goose@latest` (goose is already a Go module dep; only the CLI is missing). `migrate-up`/`migrate-down` use `go run ./cmd/migrate` and need no CLI. Run the API only after `make compose-up` (it connects to postgres + valkey and auto-runs migrations on boot).

## Database conventions

**Never edit an already-applied migration.** Create a new one via `make migrate-create NAME=xxx`. Migration files can be **either SQL (`.sql` goose blocks) or Go (`.go` calling `goose.AddMigrationContext` in `init()`)**, mixed freely in `internal/db/migrations` — version = numeric file prefix (`00001_...`), all tracked in `goose_db_version`. Use Go migrations when logic is needed (e.g. seed admin: bcrypt password in `00002_seed_admin.go`). Go migrations must live in package `migrations` and are wired in via blank imports (`internal/db/migrate.go`, `cmd/migrate/main.go`).

**Every table** embeds the same audit + public-ID pattern (see `internal/db/migrations/00001_create_users.sql` and `internal/models/base.go`):

- `id BIGSERIAL PRIMARY KEY` — internal PK, btree index. **Never exposed to FE.**
- `uuid UUID NOT NULL DEFAULT gen_random_uuid()` — exposed to FE as `"uuid"`. No unique constraint (v4 is collision-safe). Hash index for point lookups: `CREATE INDEX ... USING HASH (uuid)`.
- `created_at`, `updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()` — UTC storage, never exposed.
- `deleted_at TIMESTAMPTZ` — soft delete, nullable, never exposed.
- `created_by`, `updated_by`, `deleted_by BIGINT REFERENCES users(id)` — audit actor, nullable, never exposed.

Rules:
- FK references **the internal auto-increment `id`** (btree), never uuid.
- In API responses, only `uuid` and business fields are sent. Map any related-row FK to its `*_uuid` in response DTOs. `id`, `_at`, `_by` must be `json:"-"`.
- Use `TIMESTAMPTZ` for all time columns. Never epoch/integer timestamps.
- Model structs embed `models.BaseID` + `models.Audit`; Go fields for `id`/`uuid`/audit stay `json:"-"` in the shared structs.

## Code conventions

- Go 1.26, module `simon-jp-api`. No comments unless requested.
- Handlers are thin; logic lives in `service`. Services return sentinel errors (e.g. `service.ErrInvalidCredentials`) matched with `errors.Is`.
- JSON responses via `httpapi.writeJSON` / `httpapi.writeError`.
- Config via env vars (`internal/config`), loaded with `github.com/caarlos0/env/v11`.
- Auth: bearer token sessions stored in Redis (`service.SessionStore`).
- New tables: migration + `models` struct (embed `BaseID` + `Audit`) + `repository` + `service` + `httpapi` handler/route.
