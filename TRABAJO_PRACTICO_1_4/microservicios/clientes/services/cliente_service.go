package services

import (
	"errors"
	"fmt"
	"strings"

	"clientes/models"
	"clientes/repositories"
)

// ErrValidacion indica que los datos de entrada no son válidos.
var ErrValidacion = errors.New("datos de cliente inválidos")

// ClienteService contiene la lógica de negocio de clientes.
type ClienteService struct {
	Repo repositories.ClientesRepo
}

// NuevoClienteService cablea el servicio con su repositorio.
func NuevoClienteService(r repositories.ClientesRepo) ClienteService {
	return ClienteService{Repo: r}
}

// Crear valida y da de alta un cliente. El email es opcional.
func (s ClienteService) Crear(nombre, email string) (models.Cliente, error) {
	if strings.TrimSpace(nombre) == "" {
		return models.Cliente{}, fmt.Errorf("%w: el nombre es obligatorio", ErrValidacion)
	}
	return s.Repo.Crear(models.Cliente{Nombre: nombre, Email: email})
}

// BuscarPorID recupera un cliente por su identificador.
func (s ClienteService) BuscarPorID(id string) (models.Cliente, error) {
	if strings.TrimSpace(id) == "" {
		return models.Cliente{}, fmt.Errorf("%w: el id es obligatorio", ErrValidacion)
	}
	return s.Repo.BuscarPorID(id)
}
