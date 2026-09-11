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
	router := gin.New()
	controlador := ClienteController{Service: services.NuevoClienteService(repositories.NuevoClientesMemoria())}
	controlador.Registrar(router)
	return router
}

func TestPostClientes_Creado(t *testing.T) {
	router := nuevoRouter()
	grabador := httptest.NewRecorder()
	solicitud := httptest.NewRequest(http.MethodPost, "/clientes",
		strings.NewReader(`{"nombre":"Ana Perez","email":"ana@x.com"}`))
	router.ServeHTTP(grabador, solicitud)

	if grabador.Code != http.StatusCreated {
		t.Fatalf("status = %d; se esperaba 201. body: %s", grabador.Code, grabador.Body.String())
	}
	var cuerpo map[string]any
	_ = json.Unmarshal(grabador.Body.Bytes(), &cuerpo)
	if cuerpo["id"] == nil || cuerpo["id"] == "" || cuerpo["nombre"] != "Ana Perez" {
		t.Errorf("body inesperado: %v", cuerpo)
	}
}

func TestPostClientes_NombreVacio(t *testing.T) {
	router := nuevoRouter()
	grabador := httptest.NewRecorder()
	solicitud := httptest.NewRequest(http.MethodPost, "/clientes", strings.NewReader(`{"nombre":""}`))
	router.ServeHTTP(grabador, solicitud)
	if grabador.Code != http.StatusBadRequest {
		t.Errorf("status = %d; se esperaba 400", grabador.Code)
	}
}

func TestGetCliente_NoEncontrado(t *testing.T) {
	router := nuevoRouter()
	grabador := httptest.NewRecorder()
	solicitud := httptest.NewRequest(http.MethodGet, "/clientes/C-99", nil)
	router.ServeHTTP(grabador, solicitud)
	if grabador.Code != http.StatusNotFound {
		t.Errorf("status = %d; se esperaba 404", grabador.Code)
	}
}

func TestGetCliente_Encontrado(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	servicio := services.NuevoClienteService(repositories.NuevoClientesMemoria())
	creado, _ := servicio.Crear("Ana Perez", "ana@x.com")
	ClienteController{Service: servicio}.Registrar(router)

	grabador := httptest.NewRecorder()
	solicitud := httptest.NewRequest(http.MethodGet, "/clientes/"+creado.ID, nil)
	router.ServeHTTP(grabador, solicitud)
	if grabador.Code != http.StatusOK {
		t.Errorf("status = %d; se esperaba 200", grabador.Code)
	}
}
