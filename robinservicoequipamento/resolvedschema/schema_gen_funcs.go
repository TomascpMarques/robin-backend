package resolvedschema

import (
	"encoding/json"
)

func QueryParaStruct(param1 *map[string]interface{}) Query {
	var returnStruct Query
	temp, err := json.Marshal(param1)
	if err != nil {
		return Query{}
	}
	err = json.Unmarshal(temp, &returnStruct)
	if err != nil {
		return Query{}
	}
	return returnStruct
}

func RegistoMetaParaStruct(param1 *map[string]interface{}) RegistoMeta {
	var returnStruct RegistoMeta
	temp, err := json.Marshal(param1)
	if err != nil {
		return RegistoMeta{}
	}
	err = json.Unmarshal(temp, &returnStruct)
	if err != nil {
		return RegistoMeta{}
	}
	return returnStruct
}
