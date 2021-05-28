package resolvedschema

// Query -
type Query struct {
	Campos  map[string]interface{} `json:"campos,omitempty"`
	Extrair [][]interface{}        `json:"extrair,omitempty"`
}

type Registo struct {
	/* Meta do Regsito */
	TipoRegisto string
	Estado      string
	Quantidade  int64
	/* =============== */
	// Body criado pelo o user para o registo
	Body map[string]interface{}
}
