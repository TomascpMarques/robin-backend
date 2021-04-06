package resolvedschema

// Item -
type Item struct {
	Nome       string `json:"nome,omitempty"`
	Tipo       string `json:"tipo,omitempty"`
	Quantidade int    `json:"quantidade,omitempty"`
}
