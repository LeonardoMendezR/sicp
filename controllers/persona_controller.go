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

    // Validar estructura del request
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "CUIL requerido"})
        return
    }

    // Llamar al servicio
    persona, err := services.ConsultarPersonaPorCUIL(request.CUIL)
    if err != nil {
        // Manejo de errores más preciso
        switch err.Error() {
        case "respuesta incompleta del servicio":
            c.JSON(http.StatusNotFound, gin.H{"error": "La persona no fue encontrada o los datos están incompletos"})
        case "respuesta inesperada del servicio":
            c.JSON(http.StatusInternalServerError, gin.H{"error": "El servicio SOAP respondió de manera inesperada"})
        default:
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }

    // Éxito
    c.JSON(http.StatusOK, gin.H{
        "message": "Persona encontrada",
        "persona": persona,
    })
}
