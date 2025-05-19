package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"gobierno-inscripcion/routes"
)

func main() {
	// Cargar variables de entorno desde .env
	if err := godotenv.Load(); err != nil {
		log.Println("No se pudo cargar el archivo .env, usando valores por defecto")
	}

	

	// Inicializar router Gin
	router := gin.Default()
	router.Static("/static", "./static")

// Ruta raíz que carga index.html explícitamente
	router.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})


	// Registrar rutas de la API
	routes.RegisterRoutes(router)

	// Leer el puerto desde el entorno o usar 8080 por defecto
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Iniciar servidor
	log.Println("Servidor corriendo en puerto " + port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("No se pudo iniciar el servidor:", err)
	}
}
