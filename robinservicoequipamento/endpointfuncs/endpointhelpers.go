package endpointfuncs

import (
	"context"
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
)

// ExtractValuesFromJSON Itera sobre os querys a serem feitos ao registo passado (obj)
func ExtractValuesFromJSON(querys []interface{}, obj map[string]interface{}, final map[string]interface{}) {
	for _, k := range querys {
		var temp = strings.Split(k.(string), ".")
		if val := GetMapValueRecurssive(temp, obj); val != nil {
			final[k.(string)] = val
		}
	}
}

// GetMapValueRecurssive Busca o valor especifico (lastIndex do keys[]), através de recursão
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

// CheckValueIsValid Verifica se o valor é válido
func CheckValueIsValid(val interface{}) bool {
	return reflect.ValueOf(val).IsValid()
}

// CheckValMapStrInter Verifica se val é um mapa
func CheckValMapStrInter(val interface{}) bool {
	return reflect.TypeOf(val).String() == "map[string]interface {}"
}

// GetColecaoFromDB Devolve a coleção especificada, de um documento especifico
func GetColecaoFromDB(dbCollPar map[string]interface{}) *mongo.Collection {
	return MongoClient.Database(dbCollPar["db"].(string)).Collection(dbCollPar["cl"].(string))
}

// GetRegistosDaColecao Busca todos os resultados de uma coleção, que igualam ao filtro
func GetRegistosDaColecao(filter interface{}, colecao *mongo.Collection) ([]map[string]interface{}, error) {
	// Busca todos os resultados que coincidêm com o filtro passado
	cursor, err := colecao.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	// Descodifica todos os valores encontrados para um map[string]interface{}
	var res []map[string]interface{}
	if err := cursor.All(context.TODO(), &res); err != nil {
		return nil, err
	}

	// Se tudo correu bem, devolve os resultados encontrados
	return res, nil
}
