package resolvedschema

// Repositorio - Uma estrutura que guarda diretorias de ficheiros
//				 Contribuidores, o autore entre outros.

// Sem campo de ID na struct, os ficheiros ao serem guardados no mongo DB
// Recebem o ObjectID do novo registo, o ID em sí é o nome,
// Que será obrigatóriamente diferente para cada repo.
type Repositorio struct {
	Nome           string                `json:"nome,omitempty"`
	Tema           string                `json:"tema,omitempty"`
	Autor          string                `json:"autor,omitempty"`
	Contribuidores []string              `json:"contribuidores"`
	Ficheiros      []map[string][]string `json:"ficheiros,omitempty"` // i.e [{"file_name": "path/to/file.ext"},...]
}

// FicheiroMetaData - Contêm informações relativas ao ficheiro, não o conteudo em sí
type FicheiroMetaData struct {
	Nome     string   `json:"nome,omitempty"`
	Autor    string   `json:"autor,omitempty"`
	Criacao  int64    `json:"criacao,omitempty"`
	RepoNome string   `json:"repo_nome,omitempty"`
	Hash     string   `json:"hash,omitempty"` // A hash é gerada do formato json da struct, a partir dos campos: Nome, Autor, Path e RepoNome
	Path     []string `json:"path,omitempty"` // ["<repo_name>","folder1","folder2",...,"<file_name.extension>"]
}

type FicheiroConteudo struct {
	Nome     string `json:"nome,omitempty"`
	Conteudo string `json:"conteudo,omitempty"`
	Hash     string `json:"hash,omitempty"` // Hash do conteudo do ficheiro
}
