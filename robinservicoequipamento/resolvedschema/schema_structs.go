package resolvedschema

// Query -
type Query struct {
	Campos  map[string]interface{} `json:"campos,omitempty"`
	Extrair [][]interface{}        `json:"extrair,omitempty"`
}
