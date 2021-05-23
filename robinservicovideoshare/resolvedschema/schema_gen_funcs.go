package resolvedschema

import (
	"encoding/json"
)

func VideoParaStruct(param1 *map[string]interface{}) Video {
	var returnStruct Video
	temp, err := json.Marshal(param1)
	if err != nil {
		return Video{}
	}
	err = json.Unmarshal(temp, &returnStruct)
	if err != nil {
		return Video{}
	}
	return returnStruct
}

func SearchParamsParaStruct(param1 *map[string]interface{}) VideoSearchParams {
	var returnStruct VideoSearchParams
	temp, err := json.Marshal(param1)
	if err != nil {
		return VideoSearchParams{}
	}
	err = json.Unmarshal(temp, &returnStruct)
	if err != nil {
		return VideoSearchParams{}
	}
	return returnStruct
}
