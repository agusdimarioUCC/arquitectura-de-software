package services

import (
	"pedidos/models"
	"pedidos/repositories"
)

// ProductoService expone el catálogo (a través de la caché).
type ProductoService struct {
	Repo repositories.ProductosRepo
}

// NuevoProductoService cablea el servicio con su repositorio.
func NuevoProductoService(r repositories.ProductosRepo) ProductoService {
	return ProductoService{Repo: r}
}

// Listar devuelve todos los productos disponibles.
func (s ProductoService) Listar() ([]models.Producto, error) {
	return s.Repo.Listar()
}
