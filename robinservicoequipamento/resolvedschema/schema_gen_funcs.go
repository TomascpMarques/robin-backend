package resolvedschema

import (
	"encoding/json"
)

func ItemParaStruct(param1 *map[string]interface{}) Item {
	var returnStruct Item
	temp, err := json.Marshal(param1)
	if err != nil {
		return Item{}
	}
	err = json.Unmarshal(temp, &returnStruct)
	if err != nil {
		return Item{}
	}
	return returnStruct
}

