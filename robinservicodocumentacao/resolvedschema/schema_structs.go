package resolvedschema

// Repositorio - Uma estrutura que gaurda diretorias de ficheiros
//				 Contribuidores, o autore entre outros.

// Sem campo de ID na struct, os ficheiros ao serem guardados no mongo DB
// Recebem o ObjectID do novo registo, o ID em sí é o nome,
// Que será obrigatóriamente diferente para cada repo.
type Repositorio struct {
	Nome            string                      `json:"nome"`
	Tema            string                      `json:"tema"`
	Autor           string                      `json:"autor"`
	Contribuidores  []string                    `json:"contribuidores"`
	FicheirosMeta   map[string]FicheiroMetaData `json:"ficheiros_meta"`
	FicheirosPastas map[string][]string         `json:"ficheiros_pastas"`
}

type FicheiroMetaData struct {
	Nome     string `json:"nome"`
	Autor    string `json:"autor"`
	Criacao  string `json:"criacao"`
	RepoNome string `json:"repo_nome"`
}

type FicheiroConteudo struct {
	Nome     string `json:"nome"`
	Conteudo string `json:"conteudo"`
}
