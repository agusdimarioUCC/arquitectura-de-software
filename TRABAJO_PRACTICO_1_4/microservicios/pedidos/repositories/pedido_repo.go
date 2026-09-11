package repositories

import (
	"fmt"
	"sync"

	"pedidos/models"
)

// PedidosRepo define el contrato de persistencia de pedidos.
type PedidosRepo interface {
	Guardar(pedido models.Pedido) (models.Pedido, error)
}

// PedidosMemoria guarda los pedidos confirmados en RAM.
type PedidosMemoria struct {
	mutex     sync.Mutex
	datos     map[string]models.Pedido
	secuencia int
}

// NuevoPedidosMemoria crea un repositorio de pedidos vacío.
func NuevoPedidosMemoria() *PedidosMemoria {
	return &PedidosMemoria{datos: make(map[string]models.Pedido)}
}

func (repositorio *PedidosMemoria) Guardar(pedido models.Pedido) (models.Pedido, error) {
	repositorio.mutex.Lock()
	defer repositorio.mutex.Unlock()
	if pedido.ID == "" {
		repositorio.secuencia++
		pedido.ID = fmt.Sprintf("PED-%d", repositorio.secuencia)
	}
	repositorio.datos[pedido.ID] = pedido
	return pedido, nil
}
