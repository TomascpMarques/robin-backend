package reposfiles

import (
	"errors"
	"os"
	"path/filepath"
)

// VerificarDirExiste Verifica se a pasta base para o repo pedido já existe
func VerificarDirExiste(path string, dirNome string) bool {
	// Verificar se o repo já existe
	dirExiste := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && info.Name() == dirNome {
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
