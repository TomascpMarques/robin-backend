package resolvedschema

// Utilizador  -
type Utilizador struct {
	Nome           string            `json:"nome"`
	UsrNome        string            `json:"usr_nome"`
	Status         string            `json:"status"`
	Contribuicoes  map[string]string `json:"contribuicoes"`
	Especialidades []string          `json:"especialidades"`
}
