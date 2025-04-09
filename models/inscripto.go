package models

type Inscripto struct {
    CUIL     string `json:"cuil"`
    Nombre   string `json:"nombre"`
    Apellido string `json:"apellido"`
    CursoID  string `json:"curso_id"` // <-- nuevo campo
}
