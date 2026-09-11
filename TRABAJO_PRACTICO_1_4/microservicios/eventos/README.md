# Contrato del evento `pedido.confirmado`

Publicado por el microservicio `pedidos` cuando un pedido queda confirmado.
El destinatario conceptual es logística (no se implementa).

## Transporte

- Broker: RabbitMQ (`amqp://user:pass@localhost:5672`).
- Exchange: default (`""`).
- Routing key / cola: `pedidos-confirmados` (durable).
- Propiedades del mensaje: `content_type: application/json`, `delivery_mode: 2` (persistente).

## Payload

```json
{
  "tipo": "pedido.confirmado",
  "pedido_id": "PED-1",
  "cliente_id": "C-1",
  "producto_id": "P-1"
}
```

El struct correspondiente vive en `pedidos/models/evento.go`
(`models.EventoPedidoConfirmado`). No se comparte código entre módulos: cada
microservicio es un módulo Go independiente.
