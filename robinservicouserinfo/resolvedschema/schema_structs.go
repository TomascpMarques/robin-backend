package resolvedschema

// Utilizador  -
type Utilizador struct {
	Nome           string                `json:"nome"`
	User           string                `json:"user"`
	Status         string                `json:"status"`
	StatusMss      string                `json:"status_mss"`
	Contribuicoes  []map[string][]string `json:"contribuicoes"`
	Especialidades []string              `json:"especialidades"`
}
