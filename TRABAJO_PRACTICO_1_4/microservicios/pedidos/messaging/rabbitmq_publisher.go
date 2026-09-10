package messaging

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"pedidos/models"

	amqp "github.com/rabbitmq/amqp091-go"
)

// NombreCola es la cola donde se publican los pedidos confirmados.
const NombreCola = "pedidos-confirmados"

// EventoPublisher abstrae la publicación del evento pedido.confirmado.
type EventoPublisher interface {
	PublicarPedidoConfirmado(e models.EventoPedidoConfirmado) error
}

// PublisherConsola imprime el evento por stdout. Se usa como fallback
// cuando RabbitMQ no está disponible.
type PublisherConsola struct{}

func (PublisherConsola) PublicarPedidoConfirmado(e models.EventoPedidoConfirmado) error {
	cuerpo, _ := json.Marshal(e)
	log.Printf("[evento-consola] %s", cuerpo)
	return nil
}

// RabbitMQPublisher publica el evento en la cola pedidos-confirmados.
type RabbitMQPublisher struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

// NuevoRabbitMQPublisher abre la conexión, el canal y declara la cola durable.
func NuevoRabbitMQPublisher(uri string) (*RabbitMQPublisher, error) {
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}
	if _, err := ch.QueueDeclare(NombreCola, true, false, false, false, nil); err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}
	return &RabbitMQPublisher{conn: conn, ch: ch}, nil
}

func (p *RabbitMQPublisher) PublicarPedidoConfirmado(e models.EventoPedidoConfirmado) error {
	cuerpo, err := json.Marshal(e)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return p.ch.PublishWithContext(ctx, "", NombreCola, false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Body:         cuerpo,
	})
}

// Cerrar libera el canal y la conexión.
func (p *RabbitMQPublisher) Cerrar() {
	if p.ch != nil {
		p.ch.Close()
	}
	if p.conn != nil {
		p.conn.Close()
	}
}
