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
var HomePath = "../../repo"

// func AtualizarHomeDir(newDir string) {
// 	HomePath = newDir
// }

// CriarRepositorio_repo Cria a pasta para o repositorio dados nos parametros
func CriarRepositorio_repo(repo *resolvedschema.Repositorio) error {
	// Mudar para a diretoria dos repos
	os.Chdir(HomePath)

	workingDir, _ := os.Getwd()
	if !VerificarDirBase(workingDir) {
		fmt.Println(workingDir)
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
	os.Chdir(HomePath)

	// Verificar se a dir base é a correta
	dirAtual, _ := os.Getwd()
	if !VerificarDirBase(dirAtual) {
		return errors.New("a dir base está errada")
	}

	// Walk path, cria as pastas necessárias, e muda de dir para essas mesmas
	for _, dir := range ficheiro.Path[1 : len(ficheiro.Path)-1] {
		if _, existe := ioutil.ReadDir("./" + dir); existe != nil {
			if err := os.Mkdir(("./" + dir), os.FileMode(6640)); err != nil {
				fmt.Println("Erro: ", err)
				return errors.New("não foi possivél criar a pasta")
			}
		}
		// Muda para a dir correspondente à que se encontra dentro de valor
		os.Chdir(dir)
	}

	// Cria o ficheiro Vazio
	novoFicheiro, err := os.Create(ficheiro.Nome)
	if err != nil {
		return err
	}
	defer novoFicheiro.Close()

	return nil
}
