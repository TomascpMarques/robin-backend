package resolvedschema

// Utilizador -
type Utilizador struct {
	Nome           string            `json:"nome"`
	Status         string            `json:"status"`
	Contribuicoes  map[string]string `json:"contribuicoes"`
	Especialidades []string          `json:"especialidades"`
}
