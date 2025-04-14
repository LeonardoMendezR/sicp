	package controllers

	import (
		"net/http"

		"github.com/gin-gonic/gin"
		"gobierno-inscripcion/services"
	)

	type VerificacionRequest struct {
		CUIL string `json:"cuil" binding:"required"`
	}

	func VerificarPersona(c *gin.Context) {
		var request VerificacionRequest

		// Validar JSON de entrada
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "CUIL requerido en formato JSON"})
			return
		}

		// Llamar al servicio
		persona, err := services.ConsultarPersonaPorCUIL(request.CUIL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "No se pudo consultar la persona",
				"detalle": err.Error(),
			})
			return
		}

		// Respuesta exitosa
		c.JSON(http.StatusOK, gin.H{
			"message": "Persona encontrada exitosamente",
			"persona": persona,
		})
	}
