package controllers

import (
	"errors"
	"net/http"

	"clientes/repositories"
	"clientes/services"

	"github.com/gin-gonic/gin"
)

// ClienteController expone los endpoints HTTP de clientes.
type ClienteController struct {
	Service services.ClienteService
}

type crearClienteRequest struct {
	Nombre string `json:"nombre"`
	Email  string `json:"email"`
}

// Registrar monta las rutas en el router.
func (cc ClienteController) Registrar(r *gin.Engine) {
	r.POST("/clientes", cc.Crear)
	r.GET("/clientes/:id", cc.ObtenerPorID)
}

// Crear maneja POST /clientes.
func (cc ClienteController) Crear(ctx *gin.Context) {
	var req crearClienteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "body inválido"})
		return
	}
	cliente, err := cc.Service.Crear(req.Nombre, req.Email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, cliente)
}

// ObtenerPorID maneja GET /clientes/:id.
func (cc ClienteController) ObtenerPorID(ctx *gin.Context) {
	cliente, err := cc.Service.BuscarPorID(ctx.Param("id"))
	if err != nil {
		if errors.Is(err, repositories.ErrClienteNoEncontrado) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, cliente)
}
