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
	Crear(cliente models.Cliente) (models.Cliente, error)
	BuscarPorID(identificador string) (models.Cliente, error)
}

// ClientesMemoria es una implementación en RAM de ClientesRepo.
type ClientesMemoria struct {
	mutex     sync.Mutex
	datos     map[string]models.Cliente
	secuencia int
}

// NuevoClientesMemoria crea un repositorio de clientes vacío.
func NuevoClientesMemoria() *ClientesMemoria {
	return &ClientesMemoria{datos: make(map[string]models.Cliente)}
}

func (repositorio *ClientesMemoria) Crear(cliente models.Cliente) (models.Cliente, error) {
	repositorio.mutex.Lock()
	defer repositorio.mutex.Unlock()
	repositorio.secuencia++
	cliente.ID = fmt.Sprintf("C-%d", repositorio.secuencia)
	repositorio.datos[cliente.ID] = cliente
	return cliente, nil
}

func (repositorio *ClientesMemoria) BuscarPorID(identificador string) (models.Cliente, error) {
	repositorio.mutex.Lock()
	defer repositorio.mutex.Unlock()
	cliente, existe := repositorio.datos[identificador]
	if !existe {
		return models.Cliente{}, ErrClienteNoEncontrado
	}
	return cliente, nil
}
