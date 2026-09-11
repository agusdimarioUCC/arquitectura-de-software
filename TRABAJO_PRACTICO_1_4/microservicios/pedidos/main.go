package main

import (
	"log"
	"time"

	"pedidos/controllers"
	"pedidos/messaging"
	"pedidos/repositories"
	"pedidos/services"

	"github.com/gin-gonic/gin"
)

const (
	puerto   = ":8082"
	amqpURI  = "amqp://user:pass@localhost:5672"
	cacheTTL = 30 * time.Second
)

func main() {
	// Catálogo real envuelto en la caché con TTL.
	productos := repositories.NuevoProductosCache(repositories.NuevoProductosMemoria(), cacheTTL)
	pedidos := repositories.NuevoPedidosMemoria()

	// Publisher: RabbitMQ si está disponible, si no consola.
	var publisher messaging.EventoPublisher
	rabbit, err := messaging.NuevoRabbitMQPublisher(amqpURI)
	if err != nil {
		log.Printf("WARN: RabbitMQ no disponible (%v). Usando PublisherConsola.", err)
		publisher = messaging.PublisherConsola{}
	} else {
		log.Printf("RabbitMQ conectado, publicando en la cola %q", messaging.NombreCola)
		defer rabbit.Cerrar()
		publisher = rabbit
	}

	// Cableado de capas.
	productoCtrl := controllers.ProductoController{Service: services.NuevoProductoService(productos)}
	pedidoCtrl := controllers.PedidoController{Service: services.NuevoPedidoService(productos, pedidos, publisher)}

	router := gin.Default()
	productoCtrl.Registrar(router)
	pedidoCtrl.Registrar(router)

	log.Printf("Servicio pedidos escuchando en http://localhost%s", puerto)
	if err := router.Run(puerto); err != nil {
		log.Fatal(err)
	}
}
