package controllers

import (
	"net/http"

	"pedidos/services"

	"github.com/gin-gonic/gin"
)

// ProductoController expone el catálogo por HTTP.
type ProductoController struct {
	Service services.ProductoService
}

// Registrar monta GET /productos.
func (pc ProductoController) Registrar(r *gin.Engine) {
	r.GET("/productos", pc.Listar)
}

// Listar maneja GET /productos.
func (pc ProductoController) Listar(ctx *gin.Context) {
	productos, err := pc.Service.Listar()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"productos": productos})
}
