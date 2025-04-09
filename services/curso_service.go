package services

import (
    "encoding/csv"
    "fmt"
    "os"
    "gobierno-inscripcion/models"
)

var CursosDisponibles []models.Curso

func CargarCursosDesdeCSV(path string) error {
    file, err := os.Open(path)
    if err != nil {
        return fmt.Errorf("error al abrir el archivo: %v", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    reader.Comma = ';' // ¡IMPORTANTE! porque tu CSV usa ; en vez de ,
    
    records, err := reader.ReadAll()
    if err != nil {
        return fmt.Errorf("error al leer el CSV: %v", err)
    }

    var cursos []models.Curso
    for i, record := range records {
        if i == 0 {
            continue // Saltar encabezado
        }

        if len(record) < 2 {
            continue // línea incompleta
        }

        curso := models.Curso{
            ID:     record[0],
            Nombre: record[1],
        }
        cursos = append(cursos, curso)
    }

    CursosDisponibles = cursos
    fmt.Printf("Cursos cargados: %+v\n", CursosDisponibles) // Para debug
    return nil
}
