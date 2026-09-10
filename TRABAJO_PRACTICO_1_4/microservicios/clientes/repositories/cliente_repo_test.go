package repositories

import (
	"errors"
	"testing"

	"clientes/models"
)

func TestClientesMemoria_CrearAsignaIDYPersiste(t *testing.T) {
	repo := NuevoClientesMemoria()

	creado, err := repo.Crear(models.Cliente{Nombre: "Ana Perez", Email: "ana@x.com"})
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if creado.ID == "" {
		t.Fatal("se esperaba un ID asignado")
	}

	encontrado, err := repo.BuscarPorID(creado.ID)
	if err != nil {
		t.Fatalf("no se esperaba error al buscar: %v", err)
	}
	if encontrado.Nombre != "Ana Perez" {
		t.Errorf("nombre = %q, se esperaba %q", encontrado.Nombre, "Ana Perez")
	}
}

func TestClientesMemoria_IDsCorrelativos(t *testing.T) {
	repo := NuevoClientesMemoria()
	c1, _ := repo.Crear(models.Cliente{Nombre: "Uno"})
	c2, _ := repo.Crear(models.Cliente{Nombre: "Dos"})
	if c1.ID != "C-1" || c2.ID != "C-2" {
		t.Errorf("IDs = %q, %q; se esperaba C-1, C-2", c1.ID, c2.ID)
	}
}

func TestClientesMemoria_BuscarInexistente(t *testing.T) {
	repo := NuevoClientesMemoria()
	_, err := repo.BuscarPorID("C-99")
	if !errors.Is(err, ErrClienteNoEncontrado) {
		t.Errorf("err = %v; se esperaba ErrClienteNoEncontrado", err)
	}
}
