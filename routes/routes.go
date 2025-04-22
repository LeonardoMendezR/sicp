package routes

import (
	"github.com/gin-gonic/gin"
	"gobierno-inscripcion/controllers"
)

func RegisterRoutes(r *gin.Engine) {
	// 📌 Verificación de persona (consulta externa por CUIL)
	r.POST("/verificar-persona", controllers.BuscarPersonaHandler)

	// 📝 Gestión de inscripciones manuales
	r.POST("/inscribir-persona", controllers.InscribirPersona)
	r.GET("/inscriptos", controllers.ObtenerInscritos)
	
	// Este endpoint puede ser otro handler diferente si querés buscar inscripto por CUIL localmente
	// Por ahora lo comento o podrías hacer uno específico
	// r.GET("/buscar-inscripto/:cuil", controllers.BuscarInscripto)

	r.POST("/resetear-inscriptos", controllers.ResetearInscripciones)
	r.POST("/guardar-inscriptos", controllers.GuardarInscripcionesEnJSON)
	r.POST("/cargar-inscriptos", controllers.CargarInscriptosDesdeJSON)

	// 📚 Cursos disponibles
	r.GET("/cursos", controllers.ObtenerCursos)

	// 🛠️ Futuro: generación de QR, carga masiva, filtros, etc.
	// r.GET("/generar-qr", controllers.GenerarQR)
	// r.POST("/cargar-inscriptos-desde-json", controllers.CargarDesdeBackup)
}
