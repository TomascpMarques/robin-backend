package resolvedschema

// Repositorio - Uma estrutura que guarda diretorias de ficheiros
//				 Contribuidores, o autore entre outros.

// Sem campo de ID na struct, os ficheiros ao serem guardados no mongo DB
// Recebem o ObjectID do novo registo, o ID em sí é o nome,
// Que será obrigatóriamente diferente para cada repo.
type Repositorio struct {
	Nome           string   `json:"nome,omitempty"`
	Tema           string   `json:"tema,omitempty"`
	Autor          string   `json:"autor,omitempty"`
	Contribuidores []string `json:"contribuidores,omitempty"`
	Ficheiros      []string `json:"ficheiros,omitempty"`
}

// FicheiroMetaData - Contêm informações relativas ao ficheiro, não o conteudo em sí
type FicheiroMetaData struct {
	Nome     string   `json:"nome,omitempty"`
	Autor    string   `json:"autor,omitempty"`
	Criacao  string   `json:"criacao,omitempty"`
	RepoNome string   `json:"repo_nome,omitempty"`
	Hash     string   `json:"hash,omitempty"`
	Path     []string `json:"path,omitempty"` // ["<repo_name>","folder1","folder2",...]
}

type FicheiroConteudo struct {
	Nome     string `json:"nome,omitempty"`
	Conteudo string `json:"conteudo,omitempty"`
	Hash     string `json:"hash,omitempty"`
}
