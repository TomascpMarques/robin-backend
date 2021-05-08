package resolvedschema

// Repositorio - Uma estrutura que gaurda diretorias de ficheiros
//				 Contribuidores, o autore entre outros.

// Sem campo de ID na struct, os ficheiros ao serem guardados no mongo DB
// Recebem o ObjectID do novo registo, o ID em sí é o nome,
// Que será obrigatóriamente diferente para cada repo.
type Repositorio struct {
	Nome            string              `json:"nome,omitempty"`
	Tema            string              `json:"tema,omitempty"`
	Autor           string              `json:"autor,omitempty"`
	Contribuidores  []string            `json:"contribuidores,omitempty"`
	FicheirosMeta   map[string]string   `json:"ficheiros_meta,omitempty"`
	FicheirosPastas map[string][]string `json:"ficheiros_pastas,omitempty"`
}

type FicheiroMetaData struct {
	Nome     string `json:"nome,omitempty"`
	Autor    string `json:"autor,omitempty"`
	Criacao  string `json:"criacao,omitempty"`
	RepoNome string `json:"repo_nome,omitempty"`
}

type FicheiroConteudo struct {
	Nome     string `json:"nome,omitempty"`
	Conteudo string `json:"conteudo,omitempty"`
}
