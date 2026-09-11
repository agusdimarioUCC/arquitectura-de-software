package controllers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"pedidos/messaging"
	"pedidos/models"
	"pedidos/repositories"
	"pedidos/services"

	"github.com/gin-gonic/gin"
)

type publicadorSimulado struct{ errorDevuelto error }

func (simulado publicadorSimulado) PublicarPedidoConfirmado(models.EventoPedidoConfirmado) error {
	return simulado.errorDevuelto
}

func routerCon(publicador messaging.EventoPublisher) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	catalogo := repositories.NuevoProductosMemoria()
	servicioProductos := services.NuevoProductoService(catalogo)
	servicioPedidos := services.NuevoPedidoService(catalogo, repositories.NuevoPedidosMemoria(), publicador)
	ProductoController{Service: servicioProductos}.Registrar(router)
	PedidoController{Service: servicioPedidos}.Registrar(router)
	return router
}

func hacer(router *gin.Engine, metodo, ruta, body string) *httptest.ResponseRecorder {
	grabador := httptest.NewRecorder()
	solicitud := httptest.NewRequest(metodo, ruta, strings.NewReader(body))
	router.ServeHTTP(grabador, solicitud)
	return grabador
}

func TestGetProductos_Ok(t *testing.T) {
	grabador := hacer(routerCon(messaging.PublisherConsola{}), http.MethodGet, "/productos", "")
	if grabador.Code != http.StatusOK {
		t.Fatalf("status = %d; se esperaba 200", grabador.Code)
	}
	if !strings.Contains(grabador.Body.String(), `"productos"`) {
		t.Errorf("body sin clave productos: %s", grabador.Body.String())
	}
}

func TestPostPedidos_Creado(t *testing.T) {
	grabador := hacer(routerCon(messaging.PublisherConsola{}), http.MethodPost, "/pedidos",
		`{"cliente_id":"C-1","producto_id":"P-1","cantidad":2}`)
	if grabador.Code != http.StatusCreated {
		t.Fatalf("status = %d; se esperaba 201. body: %s", grabador.Code, grabador.Body.String())
	}
	if !strings.Contains(grabador.Body.String(), `"pedido_id"`) {
		t.Errorf("body sin pedido_id: %s", grabador.Body.String())
	}
}

func TestPostPedidos_ProductoInexistente(t *testing.T) {
	grabador := hacer(routerCon(messaging.PublisherConsola{}), http.MethodPost, "/pedidos",
		`{"cliente_id":"C-1","producto_id":"P-99","cantidad":1}`)
	if grabador.Code != http.StatusNotFound {
		t.Errorf("status = %d; se esperaba 404", grabador.Code)
	}
}

func TestPostPedidos_Validacion(t *testing.T) {
	grabador := hacer(routerCon(messaging.PublisherConsola{}), http.MethodPost, "/pedidos",
		`{"cliente_id":"C-1","producto_id":"P-1","cantidad":0}`)
	if grabador.Code != http.StatusBadRequest {
		t.Errorf("status = %d; se esperaba 400", grabador.Code)
	}
}

func TestPostPedidos_FalloPublicacion(t *testing.T) {
	grabador := hacer(routerCon(publicadorSimulado{errorDevuelto: errors.New("broker caído")}), http.MethodPost, "/pedidos",
		`{"cliente_id":"C-1","producto_id":"P-1","cantidad":1}`)
	if grabador.Code != http.StatusBadGateway {
		t.Errorf("status = %d; se esperaba 502", grabador.Code)
	}
}
