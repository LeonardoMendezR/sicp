package controllers

import (
	"encoding/json"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"gobierno-inscripcion/services"
)

type InscripcionRequest struct {
	Cuil    string `json:"cuil" binding:"required"`
	CursoID string `json:"curso_id" binding:"required"`
}

type Inscripto struct {
	Cuil         string           `json:"cuil"`
	CursoID      string           `json:"curso_id"`
	DatosPersona services.Persona `json:"datos_persona"`
}


var inscritos []Inscripto
var mu sync.Mutex

// POST /inscribir-persona
func InscribirPersona(c *gin.Context) {
	var request InscripcionRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "CUIL y CursoID son requeridos"})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for _, p := range inscritos {
		if p.Cuil == request.Cuil && p.CursoID == request.CursoID {
			c.JSON(http.StatusConflict, gin.H{"error": "La persona ya está inscrita en este curso"})
			return
		}
	}

	// Validar curso
	cursoValido := false
	for _, curso := range services.CursosDisponibles {
		if curso.ID == request.CursoID {
			cursoValido = true
			break
		}
	}
	if !cursoValido {
		c.JSON(http.StatusBadRequest, gin.H{"error": "CursoID no válido"})
		return
	}

	persona, err := services.ConsultarPersonaPorCUIL(request.Cuil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

    inscripto := Inscripto{
        Cuil:         request.Cuil,
        CursoID:      request.CursoID,
        DatosPersona: *persona,
    }
    
	inscritos = append(inscritos, inscripto)

	c.JSON(http.StatusOK, gin.H{
		"message": "Persona inscrita correctamente",
		"persona": inscripto,
	})
}

// GET /inscriptos
func ObtenerInscritos(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()
	c.JSON(http.StatusOK, inscritos)
}

// POST /resetear-inscriptos
func ResetearInscripciones(c *gin.Context) {
	mu.Lock()
	inscritos = []Inscripto{}
	mu.Unlock()
	c.JSON(http.StatusOK, gin.H{"message": "Lista de inscriptos reseteada"})
}

// GET /buscar-inscripto?cuil=XXXXXXXXXXX
func BuscarInscriptoPorCUIL(c *gin.Context) {
	cuil := c.Query("cuil")
	if cuil == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "CUIL es requerido"})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for _, p := range inscritos {
		if p.Cuil == cuil {
			var nombreCurso string
			for _, curso := range services.CursosDisponibles {
				if curso.ID == p.CursoID {
					nombreCurso = curso.Nombre
					break
				}
			}

			c.JSON(http.StatusOK, gin.H{
				"persona":      p,
				"nombre_curso": nombreCurso,
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Persona no encontrada"})
}

// POST /guardar-inscriptos
func GuardarInscripcionesEnJSON(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()

	data, err := json.MarshalIndent(inscritos, "", "  ")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al serializar datos"})
		return
	}

	if err := os.WriteFile("inscriptos_backup.json", data, 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar archivo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Inscriptos guardados en inscriptos_backup.json"})
}

// POST /cargar-inscriptos
func CargarInscriptosDesdeJSON(c *gin.Context) {
	data, err := os.ReadFile("inscriptos_backup.json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo leer el archivo"})
		return
	}

	var backup []Inscripto
	if err := json.Unmarshal(data, &backup); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al parsear JSON"})
		return
	}

	mu.Lock()
	inscritos = backup
	mu.Unlock()

	c.JSON(http.StatusOK, gin.H{
		"message":  "Backup cargado correctamente",
		"cantidad": len(inscritos),
	})
}
