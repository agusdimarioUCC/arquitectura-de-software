package repositories

import (
	"testing"
	"time"

	"pedidos/models"
)

type productosRepoSimulado struct {
	llamadasListar int
	datos          []models.Producto
}

func (simulado *productosRepoSimulado) Listar() ([]models.Producto, error) {
	simulado.llamadasListar++
	copia := make([]models.Producto, len(simulado.datos))
	copy(copia, simulado.datos)
	return copia, nil
}

func (simulado *productosRepoSimulado) BuscarPorID(identificador string) (models.Producto, error) {
	for _, producto := range simulado.datos {
		if producto.ID == identificador {
			return producto, nil
		}
	}
	return models.Producto{}, ErrProductoNoEncontrado
}

func TestProductosCache_SegundaLlamadaNoTocaElRepo(t *testing.T) {
	simulado := &productosRepoSimulado{datos: []models.Producto{{ID: "P-1", Nombre: "Auriculares"}}}
	cache := NuevoProductosCache(simulado, time.Minute)

	if _, err := cache.Listar(); err != nil {
		t.Fatalf("error: %v", err)
	}
	if _, err := cache.Listar(); err != nil {
		t.Fatalf("error: %v", err)
	}
	if simulado.llamadasListar != 1 {
		t.Errorf("llamadasListar = %d; se esperaba 1 (segunda desde caché)", simulado.llamadasListar)
	}
}

func TestProductosCache_RecargaAlExpirarElTTL(t *testing.T) {
	simulado := &productosRepoSimulado{datos: []models.Producto{{ID: "P-1"}}}
	cache := NuevoProductosCache(simulado, 20*time.Millisecond)

	_, _ = cache.Listar()
	time.Sleep(40 * time.Millisecond)
	_, _ = cache.Listar()

	if simulado.llamadasListar != 2 {
		t.Errorf("llamadasListar = %d; se esperaba 2 (recarga tras expirar)", simulado.llamadasListar)
	}
}

func TestProductosCache_BuscarPorIDUsaLaCache(t *testing.T) {
	simulado := &productosRepoSimulado{datos: []models.Producto{{ID: "P-1", Nombre: "Auriculares"}}}
	cache := NuevoProductosCache(simulado, time.Minute)

	_, _ = cache.Listar()
	producto, err := cache.BuscarPorID("P-1")
	if err != nil || producto.Nombre != "Auriculares" {
		t.Fatalf("producto = %+v, err = %v", producto, err)
	}
	if simulado.llamadasListar != 1 {
		t.Errorf("llamadasListar = %d; se esperaba 1", simulado.llamadasListar)
	}
}
