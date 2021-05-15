package reposfiles

import (
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"

	"github.com/tomascpmarques/PAP/backend/robinservicodocumentacao/loggers"
	"github.com/tomascpmarques/PAP/backend/robinservicodocumentacao/resolvedschema"
)

// Caminho da dir repo, a partir da pasta reposfiles
var HomePath = "./repo"

// CriarRepositorio_repo Cria a pasta para o repositorio dado nos parametros
func CriarRepositorio_repo(repo *resolvedschema.Repositorio) error {
	workingDir, _ := os.Getwd()
	if !VerificarDirBase(workingDir) {
		// Mudar para a diretoria dos repos
		// E verifica se hove algum erro no processo
		if err := os.Chdir(HomePath); err != nil {
			fmt.Println(err)
			return err
		}
	}

	// Verificar se o repo já existe
	root, _ := os.Getwd()
	if dirExiste := VerificarDirExiste(root, repo.Nome); dirExiste {
		loggers.DocsStorage.Println(dirExiste)
		return errors.New("o repo já existe")
	}

	// Verifica se foi possivél criar a dir
	if criarDirerr := os.Mkdir(repo.Nome, fs.FileMode(6640)); criarDirerr != nil {
		loggers.DocsStorage.Println(criarDirerr)
		return criarDirerr
	}

	return nil
}

// ApagarRepositorio_repo Apaga a pasta para o repositorio dado nos parametros
func ApagarRepositorio_repo(repo *resolvedschema.Repositorio) error {
	workingDir, _ := os.Getwd()
	if !VerificarDirBase(workingDir) {
		// Mudar para a diretoria dos repos
		// E verifica se hove algum erro no processo
		if err := os.Chdir(HomePath); err != nil {
			fmt.Println(err)
			return err
		}
	}

	// Remove a dir do repo especificada
	if err := os.Remove(repo.Nome); err != nil {
		loggers.DocsStorage.Println(err)
		return errors.New("não foi possivél aoagar o repo no storage")
	}

	return nil
}

// CriarFicheiro_repo Cria o ficheiro dentro das pastas corretas no repo(dir) especificado pelos params fornecidos
func CriarFicheiro_repo(ficheiro *resolvedschema.FicheiroMetaData) error {
	workingDir, _ := os.Getwd()
	if !VerificarDirBase(workingDir) {
		// Mudar para a diretoria dos repos
		// E verifica se hove algum erro no processo
		if err := os.Chdir(HomePath); err != nil {
			fmt.Println(err)
			return err
		}
	}

	// Walk path, cria as pastas necessárias, e muda de dir para essas mesmas
	for _, dir := range ficheiro.Path[1 : len(ficheiro.Path)-1] {
		if _, existe := ioutil.ReadDir("./" + dir); existe != nil {
			if err := os.Mkdir(("./" + dir), os.FileMode(6640)); err != nil {
				loggers.DocsStorage.Println(err)
				return errors.New("não foi possivél criar a pasta")
			}
		}
		// Muda para a dir correspondente à que se encontra dentro de valor
		os.Chdir(dir)
	}

	// Cria o ficheiro Vazio
	novoFicheiro, err := os.Create(ficheiro.Nome)
	if err != nil {
		loggers.DocsStorage.Println(err)
		return err
	}
	// Fehca SEMPRE O FICHEIRO, bery(== a very com o v trocado por um b) nice :)
	defer novoFicheiro.Close()

	return nil
}
