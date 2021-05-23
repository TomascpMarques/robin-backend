package resolvedschema

type Video struct {
	URL       string `json:"url,omitempty"`
	Tema      string `json:"tema,omitempty"`
	Titulo    string `json:"titulo,omitempty"`
	Criador   string `json:"criador,omitempty"`
	Descricao string `json:"descricao,omitempty"`
}
