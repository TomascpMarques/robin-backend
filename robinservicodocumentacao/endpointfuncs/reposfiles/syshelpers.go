package reposfiles

import (
	"errors"
	"os"
	"path/filepath"
)

// VerificarDirExiste Verifica se a pasta base para o repo pedido já existe
func VerificarDirExiste(dir string, repoNome string) bool {
	// Verificar se o repo já existe
	dirExiste := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && info.Name() == repoNome {
			return errors.New("repo já existe")
		}
		return nil
	})
	return dirExiste != nil
}

// Verifica se a dir base "home" é a dir "repo"
func VerificarDirBase(dir string) bool {
	// Verificar que estamos na diretoria certa
	return filepath.Base(dir) == "repo"
}
