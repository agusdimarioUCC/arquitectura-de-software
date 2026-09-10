package repositories

import (
	"log"
	"sync"
	"time"

	"pedidos/models"
)

// ProductosCache decora a otro ProductosRepo con una caché en memoria
// de tiempo de vida acotado (TTL). Implementa ProductosRepo.
type ProductosCache struct {
	Next ProductosRepo
	TTL  time.Duration

	mu        sync.Mutex
	cache     []models.Producto
	cargadoEn time.Time
}

// NuevoProductosCache crea la caché sobre el repositorio indicado.
func NuevoProductosCache(next ProductosRepo, ttl time.Duration) *ProductosCache {
	return &ProductosCache{Next: next, TTL: ttl}
}

func (c *ProductosCache) Listar() ([]models.Producto, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.cache != nil && time.Since(c.cargadoEn) < c.TTL {
		log.Println("CACHE HIT (productos)")
		return c.copiaCache(), nil
	}

	log.Println("CACHE MISS (productos)")
	lista, err := c.Next.Listar()
	if err != nil {
		return nil, err
	}
	c.cache = lista
	c.cargadoEn = time.Now()
	return c.copiaCache(), nil
}

func (c *ProductosCache) copiaCache() []models.Producto {
	copia := make([]models.Producto, len(c.cache))
	copy(copia, c.cache)
	return copia
}

func (c *ProductosCache) BuscarPorID(id string) (models.Producto, error) {
	lista, err := c.Listar()
	if err != nil {
		return models.Producto{}, err
	}
	for _, p := range lista {
		if p.ID == id {
			return p, nil
		}
	}
	return models.Producto{}, ErrProductoNoEncontrado
}
