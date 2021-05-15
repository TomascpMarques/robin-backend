package reposfiles

import (
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
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

// CriarFicheiro_repo Cria o ficheiro dentro das pastas corretas no repo(dir) especificado pelos params fornecidos
func CriarFicheiro_repo(ficheiro *resolvedschema.FicheiroMetaData) error {
	// Repo directorie
	os.Chdir(homePath)

	// Verificar se a dir base é a correta
	workingDir, _ := os.Getwd()
	if !VerificarDirBase(workingDir) {
		return errors.New("a dir base está errada")
	}

	// Walk path, cria as pastas necessárias, e muda de dir para essas
	for _, valor := range ficheiro.Path[1 : len(ficheiro.Path)-1] {
		if _, existe := ioutil.ReadDir("./" + valor); existe != nil {
			if err := os.Mkdir(("./" + valor), os.FileMode(6640)); err != nil {
				fmt.Println("Erro: ", err)
				return errors.New("não foi possivél criar a pasta")
			}
		}
		// Muda para a dir correspondente à que se encontra dentro de valor
		os.Chdir(valor)
	}

	// Cria o ficheiro Vazio
	file, err := os.Create(ficheiro.Nome)
	if err != nil {
		return err
	}
	defer file.Close()

	return nil
}
