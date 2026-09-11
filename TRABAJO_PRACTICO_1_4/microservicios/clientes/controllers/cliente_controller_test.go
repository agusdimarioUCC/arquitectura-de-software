package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"clientes/repositories"
	"clientes/services"

	"github.com/gin-gonic/gin"
)

func nuevoRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	ctrl := ClienteController{Service: services.NuevoClienteService(repositories.NuevoClientesMemoria())}
	ctrl.Registrar(r)
	return r
}

func TestPostClientes_Creado(t *testing.T) {
	r := nuevoRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/clientes",
		strings.NewReader(`{"nombre":"Ana Perez","email":"ana@x.com"}`))
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d; se esperaba 201. body: %s", w.Code, w.Body.String())
	}
	var cuerpo map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &cuerpo)
	if cuerpo["id"] == nil || cuerpo["id"] == "" || cuerpo["nombre"] != "Ana Perez" {
		t.Errorf("body inesperado: %v", cuerpo)
	}
}

func TestPostClientes_NombreVacio(t *testing.T) {
	r := nuevoRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/clientes", strings.NewReader(`{"nombre":""}`))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d; se esperaba 400", w.Code)
	}
}

func TestGetCliente_NoEncontrado(t *testing.T) {
	r := nuevoRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/clientes/C-99", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d; se esperaba 404", w.Code)
	}
}

func TestGetCliente_Encontrado(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	svc := services.NuevoClienteService(repositories.NuevoClientesMemoria())
	creado, _ := svc.Crear("Ana Perez", "ana@x.com")
	ClienteController{Service: svc}.Registrar(r)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/clientes/"+creado.ID, nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("status = %d; se esperaba 200", w.Code)
	}
}
