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
func (pc PedidoController) Registrar(r *gin.Engine) {
	r.POST("/pedidos", pc.Confirmar)
}

// Confirmar maneja POST /pedidos.
func (pc PedidoController) Confirmar(ctx *gin.Context) {
	var req confirmarPedidoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "body inválido"})
		return
	}

	pedido, err := pc.Service.Confirmar(req.ClienteID, req.ProductoID, req.Cantidad)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrValidacion):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrProductoNoEncontrado):
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrPublicacion):
			ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error(), "pedido_id": pedido.ID})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"mensaje":   "pedido confirmado",
		"pedido_id": pedido.ID,
		"estado":    pedido.Estado,
	})
}
