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
	puerto                = ":8082"
	direccionRabbitMQ     = "amqp://user:pass@localhost:5672"
	tiempoDeVidaDeLaCache = 30 * time.Second
)

func main() {
	// Catálogo real envuelto en la caché con tiempo de vida acotado.
	productos := repositories.NuevoProductosCache(repositories.NuevoProductosMemoria(), tiempoDeVidaDeLaCache)
	pedidos := repositories.NuevoPedidosMemoria()

	// Publicador: RabbitMQ si está disponible, si no consola.
	var publicador messaging.EventoPublisher
	publicadorRabbit, err := messaging.NuevoRabbitMQPublisher(direccionRabbitMQ)
	if err != nil {
		log.Printf("WARN: RabbitMQ no disponible (%v). Usando PublisherConsola.", err)
		publicador = messaging.PublisherConsola{}
	} else {
		log.Printf("RabbitMQ conectado, publicando en la cola %q", messaging.NombreCola)
		defer publicadorRabbit.Cerrar()
		publicador = publicadorRabbit
	}

	// Cableado de capas.
	controladorProductos := controllers.ProductoController{Service: services.NuevoProductoService(productos)}
	controladorPedidos := controllers.PedidoController{Service: services.NuevoPedidoService(productos, pedidos, publicador)}

	router := gin.Default()
	controladorProductos.Registrar(router)
	controladorPedidos.Registrar(router)

	log.Printf("Servicio pedidos escuchando en http://localhost%s", puerto)
	if err := router.Run(puerto); err != nil {
		log.Fatal(err)
	}
}
