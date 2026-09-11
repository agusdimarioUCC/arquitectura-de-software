# Microservicios del TP

Reemplazan al monolito de `../monolito`. Dos servicios Go independientes,
cada uno con su `go.mod`, organizados en capas
`controllers → services → repositories` (+ `messaging` en `pedidos`).

| Servicio | Puerto | Endpoints |
| --- | --- | --- |
| `clientes` | 8081 | `POST /clientes`, `GET /clientes/:id` |
| `pedidos` | 8082 | `GET /productos`, `POST /pedidos` |

## Requisitos

- Go 1.25.
- Docker (para RabbitMQ). El servicio `pedidos` igual arranca sin RabbitMQ:
  si no puede conectarse, publica los eventos por consola.

## Puesta en marcha

```bash
# 1. RabbitMQ (desde la carpeta CLASE_4 del repo)
cd ../../CLASE_4 && docker compose up -d
#    panel: http://localhost:15672  (user / pass)

# 2. Servicio clientes (terminal A)
cd clientes && go mod tidy && go run .

# 3. Servicio pedidos (terminal B)
cd pedidos && go mod tidy && go run .
```

## Ejemplos

```bash
# Clientes
curl -s -X POST localhost:8081/clientes -d '{"nombre":"Ana Perez","email":"ana@x.com"}'
curl -s localhost:8081/clientes/C-1

# Productos: la 2da llamada dispara CACHE HIT en los logs de pedidos
curl -s localhost:8082/productos
curl -s localhost:8082/productos

# Confirmar pedido -> publica pedido.confirmado en la cola pedidos-confirmados
curl -s -X POST localhost:8082/pedidos \
  -d '{"cliente_id":"C-1","producto_id":"P-1","cantidad":2}'
```

Verificar el evento: en el panel de RabbitMQ, la cola `pedidos-confirmados`
muestra un mensaje "Ready" por cada `POST /pedidos` exitoso. Sin RabbitMQ, el
evento aparece en el log de `pedidos` con el prefijo `[evento-consola]`.

## Tests

```bash
cd clientes && go test ./...
cd pedidos  && go test ./...
```

## Contrato del evento

Ver `../microservicios/eventos/README.md`.
