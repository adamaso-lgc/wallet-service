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
| Integration tests | testcontainers-go |

## Architecture

```
cmd/server/                        → entry point
config/                            → environment config files (viper)
api/proto/v1/                      → .proto source files
gen/proto/v1/                      → generated Go (gitignored — run `make generate`)
migrations/                        → goose SQL migrations
internal/
  domain/                          → aggregates, events, value objects, errors (no dependencies)
  service/                         → one file per operation; implements WalletServiceServer
  infrastructure/
    eventstore/                    → event-sourced wallet repository + codec
    projection/                    → WalletView read model, projector, view repository
    postgres/                      → sqlc-generated queries + type converters
  bootstrap/                       → wires everything together (server, middleware, logger, config)
```

### Key design decisions

- **Event sourcing**: wallet state is derived by replaying domain events stored in `wallet_events`. No mutable state rows.
- **Synchronous projection**: the `WalletProjector` updates `wallet_views` in the same transaction as the event write — reads are always consistent with writes.
- **Vertical slices in `service/`**: each operation (`create_wallet.go`, `deposit.go`, …) is self-contained. Shared helpers (`error.go`, `mapper.go`) are unexported within the package.
- **No application layer**: commands and queries flow directly from the gRPC handler into `service.Service`, which holds the repository and view store.

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
make test               # run all tests (quiet)
make test-v             # run all tests with verbose output
make test-integration   # run integration tests only (requires Docker)
make generate           # buf generate + sqlc generate
make migrate            # apply pending migrations
make migrate-status     # show migration state
make migrate-down       # roll back last migration
```

## Roadmap

- [ ] Structured logging in `main.go`
- [ ] Prometheus metrics (gRPC interceptor + `/metrics` endpoint)
- [ ] OpenTelemetry tracing
- [ ] JWT auth interceptor
- [ ] DB ping retry on startup
