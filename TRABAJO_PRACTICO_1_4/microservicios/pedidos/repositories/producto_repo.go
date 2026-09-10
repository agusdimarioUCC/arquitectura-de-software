package repositories

import (
	"errors"
	"time"

	"pedidos/models"
)

// ErrProductoNoEncontrado se devuelve cuando un producto no existe.
var ErrProductoNoEncontrado = errors.New("producto no encontrado")

// ProductosRepo define el contrato de acceso al catálogo.
type ProductosRepo interface {
	Listar() ([]models.Producto, error)
	BuscarPorID(id string) (models.Producto, error)
}

// ProductosMemoria es un catálogo fijo en RAM.
type ProductosMemoria struct {
	datos []models.Producto
}

// NuevoProductosMemoria crea el catálogo sembrado con datos de ejemplo.
func NuevoProductosMemoria() *ProductosMemoria {
	return &ProductosMemoria{datos: []models.Producto{
		{ID: "P-1", Nombre: "Auriculares", Precio: 25000, Stock: 10},
		{ID: "P-2", Nombre: "Teclado", Precio: 40000, Stock: 8},
	}}
}

func (r *ProductosMemoria) Listar() ([]models.Producto, error) {
	// Simula la latencia de una base de datos real para que el efecto
	// de la caché sea observable.
	time.Sleep(300 * time.Millisecond)
	copia := make([]models.Producto, len(r.datos))
	copy(copia, r.datos)
	return copia, nil
}

func (r *ProductosMemoria) BuscarPorID(id string) (models.Producto, error) {
	for _, p := range r.datos {
		if p.ID == id {
			return p, nil
		}
	}
	return models.Producto{}, ErrProductoNoEncontrado
}
