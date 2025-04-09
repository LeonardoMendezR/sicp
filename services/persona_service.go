package services

import (
    "fmt"
)

// PersonaGovResponse representa los datos que devuelve el servicio del gobierno
type PersonaGovResponse struct {
    Nombre    string `json:"nombre"`
    Apellido  string `json:"apellido"`
    Telefono  string `json:"telefono"`
    Correo    string `json:"correo"`
    Cuil      string `json:"cuil"`
}

// personasMock simula una base de datos de personas (mock)
var personasMock = map[string]PersonaGovResponse{
    "20123456789": {Nombre: "Juan", Apellido: "Pérez", Telefono: "123456789", Correo: "juan@example.com", Cuil: "20123456789"},
    "20234567890": {Nombre: "Ana", Apellido: "Gómez", Telefono: "987654321", Correo: "ana@example.com", Cuil: "20234567890"},
    "20345678901": {Nombre: "Carlos", Apellido: "López", Telefono: "1122334455", Correo: "carlos@example.com", Cuil: "20345678901"},
    "20456789012": {Nombre: "María", Apellido: "Rodríguez", Telefono: "2211334455", Correo: "maria@example.com", Cuil: "20456789012"},
    "20567890123": {Nombre: "Pedro", Apellido: "Martínez", Telefono: "3311223344", Correo: "pedro@example.com", Cuil: "20567890123"},
    // Podés agregar más personas mock acá
}

// ConsultarPersonaPorCUIL consulta el servicio SOAP del gobierno (actualmente mock)
func ConsultarPersonaPorCUIL(cuil string) (*PersonaGovResponse, error) {
    if persona, ok := personasMock[cuil]; ok {
        return &persona, nil
    }

    return nil, fmt.Errorf("persona con CUIL %s no encontrada en modo mock", cuil)
}
