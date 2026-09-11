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
func (controlador ProductoController) Registrar(router *gin.Engine) {
	router.GET("/productos", controlador.Listar)
}

// Listar maneja GET /productos.
func (controlador ProductoController) Listar(contexto *gin.Context) {
	productos, err := controlador.Service.Listar()
	if err != nil {
		contexto.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	contexto.JSON(http.StatusOK, gin.H{"productos": productos})
}
