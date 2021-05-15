package reposfiles

import (
	"errors"
	"io/fs"
	"os"

	"github.com/tomascpmarques/PAP/backend/robinservicodocumentacao/resolvedschema"
)

// Caminho da dir repo, a partir da pasta reposfiles
const homePath = "../../repo"

// CriarRepositorio_repo Cria a pasta para o repositorio dados nos parametros
func CriarRepositorio_repo(repo *resolvedschema.Repositorio) error {
	// Mudar para a diretoria dos repos
	os.Chdir(homePath)

	workingDir, _ := os.Getwd()
	if !VerificarDirBase(workingDir) {
		return errors.New("a dir base está errada")
	}

	// Verificar se o repo já existe
	root, _ := os.Getwd()
	if dirExiste := VerificarDirExiste(root, repo.Nome); dirExiste {
		return errors.New("o repo já existe")
	}

	// Verifica se foi possivél criar a dir
	if criarDirerr := os.Mkdir(repo.Nome, fs.FileMode(6640)); criarDirerr != nil {
		return criarDirerr
	}

	return nil
}
