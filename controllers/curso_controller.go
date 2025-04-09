package controllers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "gobierno-inscripcion/services"
)

func ObtenerCursos(c *gin.Context) {
    c.JSON(http.StatusOK, services.CursosDisponibles)
}
