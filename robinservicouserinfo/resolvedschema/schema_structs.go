package resolvedschema

// Utilizador  -
type Utilizador struct {
	Nome           string                `json:"nome"`
	User           string                `json:"user"`
	Status         string                `json:"status"`
	Email          string                `json:"email"`
	StatusMss      string                `json:"statusmss"`
	Contribuicoes  []map[string][]string `json:"contribuicoes"`
	Especialidades []string              `json:"especialidades"`
}
