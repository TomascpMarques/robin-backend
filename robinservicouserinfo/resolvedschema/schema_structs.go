package resolvedschema

// Utilizador  -
type Utilizador struct {
	Nome           string        `json:"nome"`
	User           string        `json:"user"`
	Status         string        `json:"status"`
	Email          string        `json:"email"`
	StatusMss      string        `json:"statusmss"`
	Contribuicoes  Contribuicoes `json:"contribuicoes"`
	Especialidades []string      `json:"especialidades"`
}

type Contribuicoes struct {
	RepoNome  string   `json:"reponome,omitempty"`
	Ficheiros []string `json:"ficheiros,omitempty"`
}
