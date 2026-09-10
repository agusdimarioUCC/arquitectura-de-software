package repositories

import (
	"errors"
	"fmt"
	"sync"

	"clientes/models"
)

// ErrClienteNoEncontrado se devuelve cuando un ID no existe en el repositorio.
var ErrClienteNoEncontrado = errors.New("cliente no encontrado")

// ClientesRepo define el contrato de persistencia de clientes.
type ClientesRepo interface {
	Crear(c models.Cliente) (models.Cliente, error)
	BuscarPorID(id string) (models.Cliente, error)
}

// ClientesMemoria es una implementación en RAM de ClientesRepo.
type ClientesMemoria struct {
	mu        sync.Mutex
	datos     map[string]models.Cliente
	secuencia int
}

// NuevoClientesMemoria crea un repositorio de clientes vacío.
func NuevoClientesMemoria() *ClientesMemoria {
	return &ClientesMemoria{datos: make(map[string]models.Cliente)}
}

func (r *ClientesMemoria) Crear(c models.Cliente) (models.Cliente, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.secuencia++
	c.ID = fmt.Sprintf("C-%d", r.secuencia)
	r.datos[c.ID] = c
	return c, nil
}

func (r *ClientesMemoria) BuscarPorID(id string) (models.Cliente, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	c, ok := r.datos[id]
	if !ok {
		return models.Cliente{}, ErrClienteNoEncontrado
	}
	return c, nil
}
