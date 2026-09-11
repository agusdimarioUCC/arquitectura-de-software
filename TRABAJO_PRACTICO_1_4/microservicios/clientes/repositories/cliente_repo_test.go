package repositories

import (
	"errors"
	"testing"

	"clientes/models"
)

func TestClientesMemoria_CrearAsignaIDYPersiste(t *testing.T) {
	repositorio := NuevoClientesMemoria()

	creado, err := repositorio.Crear(models.Cliente{Nombre: "Ana Perez", Email: "ana@x.com"})
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if creado.ID == "" {
		t.Fatal("se esperaba un ID asignado")
	}

	encontrado, err := repositorio.BuscarPorID(creado.ID)
	if err != nil {
		t.Fatalf("no se esperaba error al buscar: %v", err)
	}
	if encontrado.Nombre != "Ana Perez" {
		t.Errorf("nombre = %q, se esperaba %q", encontrado.Nombre, "Ana Perez")
	}
}

func TestClientesMemoria_IDsCorrelativos(t *testing.T) {
	repositorio := NuevoClientesMemoria()
	primerCliente, _ := repositorio.Crear(models.Cliente{Nombre: "Uno"})
	segundoCliente, _ := repositorio.Crear(models.Cliente{Nombre: "Dos"})
	if primerCliente.ID != "C-1" || segundoCliente.ID != "C-2" {
		t.Errorf("IDs = %q, %q; se esperaba C-1, C-2", primerCliente.ID, segundoCliente.ID)
	}
}

func TestClientesMemoria_BuscarInexistente(t *testing.T) {
	repositorio := NuevoClientesMemoria()
	_, err := repositorio.BuscarPorID("C-99")
	if !errors.Is(err, ErrClienteNoEncontrado) {
		t.Errorf("err = %v; se esperaba ErrClienteNoEncontrado", err)
	}
}
