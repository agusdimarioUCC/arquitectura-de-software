package main

import (
	"log"

	"clientes/controllers"
	"clientes/repositories"
	"clientes/services"

	"github.com/gin-gonic/gin"
)

func main() {
	repositorio := repositories.NuevoClientesMemoria()
	servicio := services.NuevoClienteService(repositorio)
	controlador := controllers.ClienteController{Service: servicio}

	router := gin.Default()
	controlador.Registrar(router)

	log.Println("Servicio clientes escuchando en http://localhost:8081")
	if err := router.Run(":8081"); err != nil {
		log.Fatal(err)
	}
}
