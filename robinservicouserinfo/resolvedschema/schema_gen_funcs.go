package resolvedschema

import("encoding/json")

func UtilizadorParaStruct(param1 *map[string]interface{}) Utilizador {
	var returnStruct Utilizador
	temp, err := json.Marshal(param1)
	if err != nil {
		return Utilizador{} 
	}
	err = json.Unmarshal(temp, &returnStruct)
	if err != nil {
		return Utilizador{}
	}
	return returnStruct
}