package services

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"pedidos/messaging"
	"pedidos/models"
	"pedidos/repositories"
)

var (
	// ErrValidacion indica datos de entrada inválidos.
	ErrValidacion = errors.New("datos de pedido inválidos")
	// ErrProductoNoEncontrado se reexporta para que el controlador no
	// dependa del paquete repositories.
	ErrProductoNoEncontrado = repositories.ErrProductoNoEncontrado
	// ErrPublicacion indica que el pedido se confirmó pero el evento no pudo publicarse.
	ErrPublicacion = errors.New("no se pudo publicar el evento pedido.confirmado")
)

// PedidoService confirma pedidos y publica el evento correspondiente.
type PedidoService struct {
	Productos repositories.ProductosRepo
	Pedidos   repositories.PedidosRepo
	Publisher messaging.EventoPublisher
}

// NuevoPedidoService cablea el servicio con sus dependencias.
func NuevoPedidoService(prod repositories.ProductosRepo, ped repositories.PedidosRepo, pub messaging.EventoPublisher) PedidoService {
	return PedidoService{Productos: prod, Pedidos: ped, Publisher: pub}
}

// Confirmar valida el pedido, lo persiste y publica pedido.confirmado.
// Si la publicación falla, el pedido ya quedó guardado y se devuelve
// junto a un error que envuelve ErrPublicacion.
func (s PedidoService) Confirmar(clienteID, productoID string, cantidad int) (models.Pedido, error) {
	if strings.TrimSpace(clienteID) == "" || strings.TrimSpace(productoID) == "" {
		return models.Pedido{}, fmt.Errorf("%w: cliente_id y producto_id son obligatorios", ErrValidacion)
	}
	if cantidad <= 0 {
		return models.Pedido{}, fmt.Errorf("%w: la cantidad debe ser mayor a 0", ErrValidacion)
	}

	if _, err := s.Productos.BuscarPorID(productoID); err != nil {
		return models.Pedido{}, err
	}

	pedido, err := s.Pedidos.Guardar(models.Pedido{
		ClienteID:  clienteID,
		ProductoID: productoID,
		Cantidad:   cantidad,
		Estado:     "confirmado",
	})
	if err != nil {
		return models.Pedido{}, err
	}

	evento := models.EventoPedidoConfirmado{
		Tipo:       models.TipoPedidoConfirmado,
		PedidoID:   pedido.ID,
		ClienteID:  pedido.ClienteID,
		ProductoID: pedido.ProductoID,
	}
	if err := s.Publisher.PublicarPedidoConfirmado(evento); err != nil {
		log.Printf("ERROR publicando evento de %s: %v", pedido.ID, err)
		return pedido, fmt.Errorf("%w: %v", ErrPublicacion, err)
	}

	return pedido, nil
}
