package repositories

import (
	"testing"
	"time"

	"pedidos/models"
)

type productosRepoFake struct {
	llamadasListar int
	datos          []models.Producto
}

func (f *productosRepoFake) Listar() ([]models.Producto, error) {
	f.llamadasListar++
	copia := make([]models.Producto, len(f.datos))
	copy(copia, f.datos)
	return copia, nil
}

func (f *productosRepoFake) BuscarPorID(id string) (models.Producto, error) {
	for _, p := range f.datos {
		if p.ID == id {
			return p, nil
		}
	}
	return models.Producto{}, ErrProductoNoEncontrado
}

func TestProductosCache_SegundaLlamadaNoTocaElRepo(t *testing.T) {
	fake := &productosRepoFake{datos: []models.Producto{{ID: "P-1", Nombre: "Auriculares"}}}
	cache := NuevoProductosCache(fake, time.Minute)

	if _, err := cache.Listar(); err != nil {
		t.Fatalf("error: %v", err)
	}
	if _, err := cache.Listar(); err != nil {
		t.Fatalf("error: %v", err)
	}
	if fake.llamadasListar != 1 {
		t.Errorf("llamadasListar = %d; se esperaba 1 (segunda desde caché)", fake.llamadasListar)
	}
}

func TestProductosCache_RecargaAlExpirarElTTL(t *testing.T) {
	fake := &productosRepoFake{datos: []models.Producto{{ID: "P-1"}}}
	cache := NuevoProductosCache(fake, 20*time.Millisecond)

	_, _ = cache.Listar()
	time.Sleep(40 * time.Millisecond)
	_, _ = cache.Listar()

	if fake.llamadasListar != 2 {
		t.Errorf("llamadasListar = %d; se esperaba 2 (recarga tras expirar)", fake.llamadasListar)
	}
}

func TestProductosCache_BuscarPorIDUsaLaCache(t *testing.T) {
	fake := &productosRepoFake{datos: []models.Producto{{ID: "P-1", Nombre: "Auriculares"}}}
	cache := NuevoProductosCache(fake, time.Minute)

	_, _ = cache.Listar()
	p, err := cache.BuscarPorID("P-1")
	if err != nil || p.Nombre != "Auriculares" {
		t.Fatalf("p = %+v, err = %v", p, err)
	}
	if fake.llamadasListar != 1 {
		t.Errorf("llamadasListar = %d; se esperaba 1", fake.llamadasListar)
	}
}
