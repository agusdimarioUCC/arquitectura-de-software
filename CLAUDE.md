# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this repo is

Course material for **Arquitectura de Software** (Ingeniería en Informática, UCC). One folder per class (`CLASE_1`..`CLASE_4`) plus `TRABAJO_PRACTICO_1_4`. Each class folder holds a `DEMO/` (shown in class) and a `TP/` or `ACTIVIDADES/` (student work). Everything is written in Spanish, including commit messages, comments and identifiers — match that when editing.

The pedagogical arc: layered architecture + dependency injection (CLASE_1–2) → caching (CLASE_3) → async messaging with RabbitMQ (CLASE_4) → decompose a monolith into microservices + publish an event (TP).

## Module layout

There is **no workspace / root `go.mod`**. Every runnable unit is its own Go module (`go 1.27`). Always `cd` into the directory that contains the `go.mod` before running Go commands. Modules and their entrypoints:

| Path | Module | How to run |
| --- | --- | --- |
| `CLASE_1/DEMO` | `demo` | `go run ./1.Paso` … `./6.Paso` (progressive build-up of one design) |
| `CLASE_1/TP` | `practico` | `go run .` |
| `CLASE_2/DEMO` | `main` | `go run .` (needs MongoDB) |
| `CLASE_3/DEMO` | `main` | `go run .` (needs MongoDB; see caveat below) |
| `CLASE_4/DEMO/1`, `CLASE_4/DEMO/2` | — | `go run ./consumidor` and `go run ./productor` in separate terminals |
| `CLASE_4/DEMO/3` | `rabbitmq-demo-3` | `go run .` |
| `CLASE_4/ACTIVIDADES` | `rabbitmq-actividades` | one module for both actividades; `go run ./actividad-1/productor`, `go run ./actividad-2/logistica`, etc. |
| `CLASE_4/ACTIVIDADES_RESUELTAS` | — | solved versions of the above |
| `TRABAJO_PRACTICO_1_4/monolito` | `monolito-ecommerce` | `go run .` |
| `TRABAJO_PRACTICO_1_4/microservicios/{clientes,pedidos,eventos}` | — | **empty scaffolds** — student implements them |

After pulling or switching modules, run `go mod tidy` in that module if the build complains.

### CLASE_3/DEMO caveat

`main_ejercicio1.go` … `main_ejercicio4.go` all live in `package main`. Only **one** may define an active `func main()` at a time — the others must stay commented out, or `go run .` fails with `main redeclared in this block`. Switching exercises means commenting/uncommenting the corresponding `main_ejercicioN.go` (and its repo file under `repositories/items/`). The repo is often left mid-switch, so check this before assuming a build is broken.

## Tests

Only CLASE_2 and CLASE_3 have tests (`services/items/item_service_test.go`), which use the in-memory mock repo and need no infrastructure.

```bash
cd CLASE_2/DEMO && go test ./...
cd CLASE_3/DEMO && go test ./services/items -run TestGetItem_Success   # single test
```

## Infrastructure dependencies

Started manually — there is no compose file for MongoDB/Memcached.

- **MongoDB** (CLASE_2/DEMO, CLASE_3/DEMO): expected at `mongodb://root:root@localhost:27017`, database `items-api`, collection `items`. Seed it with the mongosh script at the repo root: `mongosh "mongodb://root:root@localhost:27017" db.js`.
- **Memcached** (CLASE_3 ejercicios 3–4): expected at `localhost:11211`.
- **RabbitMQ** (all of CLASE_4, and the TP `pedidos` service): `amqp://user:pass@localhost:5672`, management UI at `http://localhost:15672` (`user`/`pass`). Start it from `CLASE_4/` with `docker compose up -d` (`CLASE_4/compose.yaml`). Durable queues can't be redeclared with different args — if you change queue config mid-test, stop the programs and restart the container.

RabbitMQ programs run one-per-terminal and are stopped with `CTRL+C`. Several activities are parameterised by env vars, e.g. `CONSUMIDOR=lento DELAY_SEGUNDOS=5 go run ./consumidor`.

## Architecture conventions

The demos teach and consistently apply a 4-layer structure with constructor-style dependency injection wired in `main`:

```
controllers/  → services/  → repositories/  → models/
   (DTOs)         (domain      (interface +      (domain
                   logic)       impls)            structs)
```

- The **repository is an interface** (`ItemsRepo` in `repositories/items/items_repo.go`) with swappable implementations: `items_mock.go` (in-memory map, used by tests and offline dev), `items_mongo*.go` (real DB), and in CLASE_3 the cache decorators `items_ccache_ej2.go` / `items_memcached_ej3.go` / `items_memcached_ej4.go`.
- **Caching is a decorator repo**: a cache implementation wraps a `NextRepo ItemsRepo` and delegates on miss. The service is unaware caching exists.
- `main.go` does all wiring: build repo → inject into `ItemsService` → inject into `ItemsController` → register Gin routes → `router.Run(":8080")`.
- HTTP is **Gin** everywhere (`gin.Default()`); MongoDB uses the v1 driver `go.mongodb.org/mongo-driver` with `Database("items-api").Collection("items")`.

### CLASE_4 messaging

Plain `github.com/rabbitmq/amqp091-go`, deliberately minimal: text payloads (not JSON), **manual ACK** after simulated processing, one idea per demo/activity (fair dispatch/QoS in actividad-1, fanout exchange in actividad-2, DLQ in DEMO 2, connection backoff in DEMO 3). Each demo/activity folder has its own README with the exact run sequence.

### The TP (`TRABAJO_PRACTICO_1_4`)

Start from `monolito/` (Gin app on `:8080` with `crearCliente` / `obtenerCliente` / `listarProductos` / `confirmarPedido` in `controllers.go`). Deliverable is `microservicios/`: split into `clientes` (`:8081`) and `pedidos` (`:8082`), each layered (controller/service/repository), in-memory data, a simple product-list cache in `pedidos`, and `pedidos` publishes a `pedido.confirmado` event to the `pedidos-confirmados` queue after confirming an order. No auth/frontend/real DB/logística consumer. See `TRABAJO_PRACTICO_1_4/README.md`.

## Git

Single `main` branch. Commits follow Conventional Commits with Spanish subjects (`feat: add first TP`, `fix: edit class 4`).
