package repositories

import (
	"fmt"
	"sync"

	"pedidos/models"
)

// PedidosRepo define el contrato de persistencia de pedidos.
type PedidosRepo interface {
	Guardar(p models.Pedido) (models.Pedido, error)
}

// PedidosMemoria guarda los pedidos confirmados en RAM.
type PedidosMemoria struct {
	mu        sync.Mutex
	datos     map[string]models.Pedido
	secuencia int
}

// NuevoPedidosMemoria crea un repositorio de pedidos vacío.
func NuevoPedidosMemoria() *PedidosMemoria {
	return &PedidosMemoria{datos: make(map[string]models.Pedido)}
}

func (r *PedidosMemoria) Guardar(p models.Pedido) (models.Pedido, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if p.ID == "" {
		r.secuencia++
		p.ID = fmt.Sprintf("PED-%d", r.secuencia)
	}
	r.datos[p.ID] = p
	return p, nil
}
