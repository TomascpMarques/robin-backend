package resolvedschema

// Query -
type Query struct {
	Campos  map[string]interface{} `json:"campos,omitempty"`
	Extrair [][]interface{}        `json:"extrair,omitempty"`
}

type Registo struct {
	Meta *RegistoMeta           `json:"meta,omitempty"`
	Body map[string]interface{} `json:"body,omitempty"`
}

type RegistoMeta struct {
	Tipo       string `json:"tipo,omitempty"`
	Estado     string `json:"estado,omitempty"`
	Quantidade int64  `json:"quantidade,omitempty"`
}
