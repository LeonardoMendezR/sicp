package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gobierno-inscripcion/services"
)

type VerificacionRequest struct {
	Cuil string `json:"cuil"`
}

func BuscarPersonaHandler(c *gin.Context) {
	var req VerificacionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato inválido de JSON"})
		return
	}

	persona, err := services.ConsultarPersonaPorCUIL(req.Cuil)
	if err != nil {
		if err == services.ErrPersonaNoEncontrada {
			c.JSON(http.StatusNotFound, gin.H{"error": "Persona no encontrada"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, persona)
}
