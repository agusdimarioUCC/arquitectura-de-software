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
func (controlador ClienteController) Registrar(router *gin.Engine) {
	router.POST("/clientes", controlador.Crear)
	router.GET("/clientes/:id", controlador.ObtenerPorID)
}

// Crear maneja POST /clientes.
func (controlador ClienteController) Crear(contexto *gin.Context) {
	var solicitud crearClienteRequest
	if err := contexto.ShouldBindJSON(&solicitud); err != nil {
		contexto.JSON(http.StatusBadRequest, gin.H{"error": "body inválido"})
		return
	}
	cliente, err := controlador.Service.Crear(solicitud.Nombre, solicitud.Email)
	if err != nil {
		contexto.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	contexto.JSON(http.StatusCreated, cliente)
}

// ObtenerPorID maneja GET /clientes/:id.
func (controlador ClienteController) ObtenerPorID(contexto *gin.Context) {
	cliente, err := controlador.Service.BuscarPorID(contexto.Param("id"))
	if err != nil {
		if errors.Is(err, repositories.ErrClienteNoEncontrado) {
			contexto.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		contexto.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	contexto.JSON(http.StatusOK, cliente)
}
