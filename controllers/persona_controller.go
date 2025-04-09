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
    cuil := c.Param("cuil") // ← Extrae el CUIL de la URL

    persona, err := services.ConsultarPersonaPorCUIL(cuil)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Persona encontrada",
        "persona": persona,
    })
}

