package endpointfuncs

import (
	"reflect"
	"strings"
)

func ExtractValuesFromJSON(querys []interface{}, obj map[string]interface{}, final map[string]interface{}) {
	for _, k := range querys {
		var temp = strings.Split(k.(string), ".")
		if val := GetMapValueRecurssive(temp, obj); val != nil {
			final[k.(string)] = val
		}
	}
}

func GetMapValueRecurssive(keys []string, obj map[string]interface{}) interface{} {
	// Verifica se o valor é válido
	if CheckValueIsValid(obj[keys[0]]) {
		// Verifica se estamos no ultimo elemento, se sim,
		// devolve logo valor para a chave keys[0] (ultimo valor passado no query)
		if len(keys) <= 1 {
			return obj[keys[0]]
		}
		// Verifica se o valor é um map, se sim
		// Devolve esse mapa e o query atualizado
		if CheckValMapStrInter(obj[keys[0]]) {
			return GetMapValueRecurssive(keys[1:], obj[keys[0]].(map[string]interface{}))
		}
	}
	return nil
}

func CheckValueIsValid(val interface{}) bool {
	return reflect.ValueOf(val).IsValid()
}

func CheckValMapStrInter(val interface{}) bool {
	return reflect.TypeOf(val).String() == "map[string]interface {}"
}
