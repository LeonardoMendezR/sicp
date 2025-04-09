package main

import (
    "log"
    "os"

    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"

    "gobierno-inscripcion/routes"
    "gobierno-inscripcion/services"
)

func main() {
    // Cargar variables de entorno desde .env
    if err := godotenv.Load(); err != nil {
        log.Println("No se pudo cargar el archivo .env, usando valores por defecto")
    }

    // Cargar cursos desde CSV
    err := services.CargarCursosDesdeCSV("Cursos.csv")
    if err != nil {
        log.Fatalf("Error al cargar cursos: %v", err)
    }

    // Inicializar router Gin
    router := gin.Default()
    
    router.Static("/static", "./static")

    // Registrar todas las rutas de la app
    routes.RegisterRoutes(router)

    // Leer el puerto desde .env o usar 8080 por defecto
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    // Iniciar servidor HTTP
    log.Println("Servidor corriendo en puerto " + port)
    if err := router.Run(":" + port); err != nil {
        log.Fatal("No se pudo iniciar el servidor:", err)
    }
}
