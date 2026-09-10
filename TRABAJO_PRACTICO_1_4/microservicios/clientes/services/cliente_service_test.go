package services

import (
	"errors"
	"testing"

	"clientes/repositories"
)

func TestClienteService_CrearValido(t *testing.T) {
	svc := NuevoClienteService(repositories.NuevoClientesMemoria())

	cliente, err := svc.Crear("Ana Perez", "ana@x.com")
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if cliente.ID == "" {
		t.Error("se esperaba un ID asignado")
	}

	recuperado, err := svc.BuscarPorID(cliente.ID)
	if err != nil || recuperado.Nombre != "Ana Perez" {
		t.Errorf("recuperado = %+v, err = %v", recuperado, err)
	}
}

func TestClienteService_CrearNombreVacio(t *testing.T) {
	svc := NuevoClienteService(repositories.NuevoClientesMemoria())
	_, err := svc.Crear("   ", "ana@x.com")
	if !errors.Is(err, ErrValidacion) {
		t.Errorf("err = %v; se esperaba ErrValidacion", err)
	}
}

func TestClienteService_BuscarInexistente(t *testing.T) {
	svc := NuevoClienteService(repositories.NuevoClientesMemoria())
	_, err := svc.BuscarPorID("C-1")
	if !errors.Is(err, repositories.ErrClienteNoEncontrado) {
		t.Errorf("err = %v; se esperaba ErrClienteNoEncontrado", err)
	}
}

func TestClienteService_BuscarIDVacio(t *testing.T) {
	svc := NuevoClienteService(repositories.NuevoClientesMemoria())
	_, err := svc.BuscarPorID("")
	if !errors.Is(err, ErrValidacion) {
		t.Errorf("err = %v; se esperaba ErrValidacion", err)
	}
}
