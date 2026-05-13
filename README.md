# wallet-service

A learning project exploring **DDD**, **event sourcing**, **CQRS**, and **gRPC** in Go.

## Stack

| Concern | Technology |
|---------|-----------|
| Transport | gRPC (protobuf) |
| Database | PostgreSQL (pgx v5) |
| Queries | sqlc (generated) |
| Migrations | goose |
| Config | viper |

## Architecture

```
cmd/server/          → entry point
internal/
  domain/            → aggregates, events, errors (no dependencies)
  application/
    command/         → one file per command + handler
    query/           → one file per query + handler
  projection/        → WalletView DTO + WalletStore interface
  infrastructure/
    postgres/        → repository, projector, codec, view store
    fake/            → in-memory test doubles
  grpc/              → server, interceptors, error mapping
  bootstrap/         → wires everything together
  config/            → viper config
api/proto/v1/        → .proto source files
gen/proto/v1/        → generated Go (gitignored — run `make generate`)
migrations/          → goose SQL migrations
```

## Quick Start

```bash
# Start Postgres and apply migrations
make docker-up

# Run the server
make run

# Check health
grpcurl -plaintext localhost:50051 grpc.health.v1.Health/Check
```

## Development

```bash
make test             # run all tests
make generate         # buf generate + sqlc generate
make migrate          # apply pending migrations
make migrate-status   # show migration state
make migrate-down     # roll back last migration
```

## Roadmap

- [ ] Structured logging in `main.go`
- [ ] Prometheus metrics (gRPC interceptor + `/metrics` endpoint)
- [ ] OpenTelemetry tracing
- [ ] JWT auth interceptor
- [ ] Integration tests (testcontainers)
- [ ] DB ping retry on startup
