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

type publisherFake struct{ err error }

func (f publisherFake) PublicarPedidoConfirmado(models.EventoPedidoConfirmado) error { return f.err }

func routerCon(pub messaging.EventoPublisher) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	catalogo := repositories.NuevoProductosMemoria()
	prodSvc := services.NuevoProductoService(catalogo)
	pedSvc := services.NuevoPedidoService(catalogo, repositories.NuevoPedidosMemoria(), pub)
	ProductoController{Service: prodSvc}.Registrar(r)
	PedidoController{Service: pedSvc}.Registrar(r)
	return r
}

func hacer(r *gin.Engine, metodo, ruta, body string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(metodo, ruta, strings.NewReader(body))
	r.ServeHTTP(w, req)
	return w
}

func TestGetProductos_Ok(t *testing.T) {
	w := hacer(routerCon(messaging.PublisherConsola{}), http.MethodGet, "/productos", "")
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d; se esperaba 200", w.Code)
	}
	if !strings.Contains(w.Body.String(), `"productos"`) {
		t.Errorf("body sin clave productos: %s", w.Body.String())
	}
}

func TestPostPedidos_Creado(t *testing.T) {
	w := hacer(routerCon(messaging.PublisherConsola{}), http.MethodPost, "/pedidos",
		`{"cliente_id":"C-1","producto_id":"P-1","cantidad":2}`)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d; se esperaba 201. body: %s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), `"pedido_id"`) {
		t.Errorf("body sin pedido_id: %s", w.Body.String())
	}
}

func TestPostPedidos_ProductoInexistente(t *testing.T) {
	w := hacer(routerCon(messaging.PublisherConsola{}), http.MethodPost, "/pedidos",
		`{"cliente_id":"C-1","producto_id":"P-99","cantidad":1}`)
	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d; se esperaba 404", w.Code)
	}
}

func TestPostPedidos_Validacion(t *testing.T) {
	w := hacer(routerCon(messaging.PublisherConsola{}), http.MethodPost, "/pedidos",
		`{"cliente_id":"C-1","producto_id":"P-1","cantidad":0}`)
	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d; se esperaba 400", w.Code)
	}
}

func TestPostPedidos_FalloPublicacion(t *testing.T) {
	w := hacer(routerCon(publisherFake{err: errors.New("broker caído")}), http.MethodPost, "/pedidos",
		`{"cliente_id":"C-1","producto_id":"P-1","cantidad":1}`)
	if w.Code != http.StatusBadGateway {
		t.Errorf("status = %d; se esperaba 502", w.Code)
	}
}
