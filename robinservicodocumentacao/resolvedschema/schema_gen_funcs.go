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

func FicheiroMetaDataParaStruct(ficheiroMeta *map[string]interface{}) FicheiroMetaData {
	var returnStruct FicheiroMetaData
	temp, err := json.Marshal(ficheiroMeta)
	if err != nil {
		return FicheiroMetaData{}
	}
	err = json.Unmarshal(temp, &returnStruct)
	if err != nil {
		return FicheiroMetaData{}
	}
	return returnStruct
}

func FicheiroConteudoParaStruct(ficheiroConteudo *map[string]interface{}) FicheiroConteudo {
	var returnStruct FicheiroConteudo
	temp, err := json.Marshal(ficheiroConteudo)
	if err != nil {
		return FicheiroConteudo{}
	}
	err = json.Unmarshal(temp, &returnStruct)
	if err != nil {
		return FicheiroConteudo{}
	}
	return returnStruct
}
