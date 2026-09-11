package controllers

import (
	"errors"
	"net/http"

	"pedidos/services"

	"github.com/gin-gonic/gin"
)

// PedidoController expone la confirmación de pedidos por HTTP.
type PedidoController struct {
	Service services.PedidoService
}

type confirmarPedidoRequest struct {
	ClienteID  string `json:"cliente_id"`
	ProductoID string `json:"producto_id"`
	Cantidad   int    `json:"cantidad"`
}

// Registrar monta POST /pedidos.
func (controlador PedidoController) Registrar(router *gin.Engine) {
	router.POST("/pedidos", controlador.Confirmar)
}

// Confirmar maneja POST /pedidos.
func (controlador PedidoController) Confirmar(contexto *gin.Context) {
	var solicitud confirmarPedidoRequest
	if err := contexto.ShouldBindJSON(&solicitud); err != nil {
		contexto.JSON(http.StatusBadRequest, gin.H{"error": "body inválido"})
		return
	}

	pedido, err := controlador.Service.Confirmar(solicitud.ClienteID, solicitud.ProductoID, solicitud.Cantidad)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrValidacion):
			contexto.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrProductoNoEncontrado):
			contexto.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrPublicacion):
			contexto.JSON(http.StatusBadGateway, gin.H{"error": err.Error(), "pedido_id": pedido.ID})
		default:
			contexto.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	contexto.JSON(http.StatusCreated, gin.H{
		"mensaje":   "pedido confirmado",
		"pedido_id": pedido.ID,
		"estado":    pedido.Estado,
	})
}
