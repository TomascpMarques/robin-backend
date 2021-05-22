package resolvedschema

type Video struct {
	Titulo    string `json:"titulo,omitempty"`
	Descricao string `json:"descricao,omitempty"`
	URL       string `json:"url,omitempty"`
	Criador   string `json:"criador,omitempty"`
}
