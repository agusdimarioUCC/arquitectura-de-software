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
	PublicarPedidoConfirmado(evento models.EventoPedidoConfirmado) error
}

// PublisherConsola imprime el evento por stdout. Se usa como fallback
// cuando RabbitMQ no está disponible.
type PublisherConsola struct{}

func (PublisherConsola) PublicarPedidoConfirmado(evento models.EventoPedidoConfirmado) error {
	cuerpo, _ := json.Marshal(evento)
	log.Printf("[evento-consola] %s", cuerpo)
	return nil
}

// RabbitMQPublisher publica el evento en la cola pedidos-confirmados.
type RabbitMQPublisher struct {
	conexion *amqp.Connection
	canal    *amqp.Channel
}

// NuevoRabbitMQPublisher abre la conexión, el canal y declara la cola durable.
func NuevoRabbitMQPublisher(direccionConexion string) (*RabbitMQPublisher, error) {
	conexion, err := amqp.Dial(direccionConexion)
	if err != nil {
		return nil, err
	}
	canal, err := conexion.Channel()
	if err != nil {
		conexion.Close()
		return nil, err
	}
	if _, err := canal.QueueDeclare(NombreCola, true, false, false, false, nil); err != nil {
		canal.Close()
		conexion.Close()
		return nil, err
	}
	return &RabbitMQPublisher{conexion: conexion, canal: canal}, nil
}

func (publicador *RabbitMQPublisher) PublicarPedidoConfirmado(evento models.EventoPedidoConfirmado) error {
	cuerpo, err := json.Marshal(evento)
	if err != nil {
		return err
	}
	contexto, cancelar := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelar()
	return publicador.canal.PublishWithContext(contexto, "", NombreCola, false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Body:         cuerpo,
	})
}

// Cerrar libera el canal y la conexión.
func (publicador *RabbitMQPublisher) Cerrar() {
	if publicador.canal != nil {
		publicador.canal.Close()
	}
	if publicador.conexion != nil {
		publicador.conexion.Close()
	}
}
