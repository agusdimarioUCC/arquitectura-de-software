package repositories

import (
	"log"
	"sync"
	"time"

	"pedidos/models"
)

// ProductosCache decora a otro ProductosRepo con una caché en memoria
// de tiempo de vida acotado. Implementa ProductosRepo.
type ProductosCache struct {
	Siguiente    ProductosRepo
	TiempoDeVida time.Duration

	mutex          sync.Mutex
	productosCache []models.Producto
	cargadoEn      time.Time
}

// NuevoProductosCache crea la caché sobre el repositorio indicado.
func NuevoProductosCache(siguiente ProductosRepo, tiempoDeVida time.Duration) *ProductosCache {
	return &ProductosCache{Siguiente: siguiente, TiempoDeVida: tiempoDeVida}
}

func (decorador *ProductosCache) Listar() ([]models.Producto, error) {
	decorador.mutex.Lock()
	defer decorador.mutex.Unlock()

	if decorador.productosCache != nil && time.Since(decorador.cargadoEn) < decorador.TiempoDeVida {
		log.Println("CACHE HIT (productos)")
		return decorador.copiaCache(), nil
	}

	log.Println("CACHE MISS (productos)")
	lista, err := decorador.Siguiente.Listar()
	if err != nil {
		return nil, err
	}
	decorador.productosCache = lista
	decorador.cargadoEn = time.Now()
	return decorador.copiaCache(), nil
}

func (decorador *ProductosCache) copiaCache() []models.Producto {
	copia := make([]models.Producto, len(decorador.productosCache))
	copy(copia, decorador.productosCache)
	return copia
}

func (decorador *ProductosCache) BuscarPorID(identificador string) (models.Producto, error) {
	lista, err := decorador.Listar()
	if err != nil {
		return models.Producto{}, err
	}
	for _, producto := range lista {
		if producto.ID == identificador {
			return producto, nil
		}
	}
	return models.Producto{}, ErrProductoNoEncontrado
}
