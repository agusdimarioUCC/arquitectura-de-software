package services

import (
	"pedidos/models"
	"pedidos/repositories"
)

// ProductoService expone el catálogo (a través de la caché).
type ProductoService struct {
	Repositorio repositories.ProductosRepo
}

// NuevoProductoService cablea el servicio con su repositorio.
func NuevoProductoService(repositorio repositories.ProductosRepo) ProductoService {
	return ProductoService{Repositorio: repositorio}
}

// Listar devuelve todos los productos disponibles.
func (servicio ProductoService) Listar() ([]models.Producto, error) {
	return servicio.Repositorio.Listar()
}
