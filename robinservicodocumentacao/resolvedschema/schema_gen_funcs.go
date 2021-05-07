package resolvedschema

import (
	"encoding/json"
)

func RepositorioParaStruct(repo *map[string]interface{}) Repositorio {
	var returnStruct Repositorio
	temp, err := json.Marshal(repo)
	if err != nil {
		return Repositorio{}
	}
	err = json.Unmarshal(temp, &returnStruct)
	if err != nil {
		return Repositorio{}
	}
	return returnStruct
}
