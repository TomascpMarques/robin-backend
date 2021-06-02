package endpointfuncs

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/tomascpmarques/PAP/backend/robinservicoequipamento/resolvedschema"
	"go.mongodb.org/mongo-driver/mongo"
)

// ExtractValuesFromJSON Itera sobre os querys a serem feitos ao registo passado (obj)
func ExtractValuesFromJSON(querys []interface{}, obj map[string]interface{}, final map[string]interface{}) {
	for _, k := range querys {
		var temp = strings.Split(k.(string), ".")
		fmt.Println("Objeto: ", obj, "\nQuery: ", temp)
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
		if CheckValMapStrInterface(obj[keys[0]]) {
			return GetMapValueRecurssive(keys[1:], obj[keys[0]].(map[string]interface{}))
		}
	}
	return nil
}

// CheckValueIsValid Verifica se o valor é válido
func CheckValueIsValid(val interface{}) bool {
	return reflect.ValueOf(val).IsValid()
}

// CheckValMapStrInterface Verifica se val é um mapa
func CheckValMapStrInterface(val interface{}) bool {
	return reflect.TypeOf(val).String() == "map[string]interface {}"
}

// GetColecaoFromDB Devolve a coleção especificada, de um documento especifico
func GetColecaoFromDB(colecao string) *mongo.Collection {
	return MongoClient.Database("recursos").Collection(colecao)
}

// ReturnDB Devolve a base-de-dados dos recursos
func ReturnDB() *mongo.Database {
	return MongoClient.Database("recursos")
}

// GetRegistosDaColecao Busca todos os resultados de uma coleção, que igualam ao filtro
func GetRegistosDaColecao(filter interface{}, colecao *mongo.Collection) ([]resolvedschema.Registo, error) {
	// Busca todos os resultados que coincidêm com o filtro passado
	cursor, err := colecao.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	// Descodifica todos os valores encontrados para um map[string]interface{}
	var res []resolvedschema.Registo
	if err := cursor.All(context.TODO(), &res); err != nil {
		return nil, err
	}

	// Se tudo correu bem, devolve os resultados encontrados
	return res, nil
}

// VerificarCamposMapa Verifica se um mapa contêm as keys e valores  válidos/necessários
func VerificarCamposMapa(campos []string, mapa map[string]interface{}) error {
	for _, campo := range campos {
		if valor, existe := mapa[campo]; !(reflect.ValueOf(valor).IsValid()) || !existe {
			return errors.New("os dados fornecidos não cumpre os parametros minímos")
		}
	}
	return nil
}

// VerificarCamposMetaRegisto Verifica a validade básica dos campos do novo registo
func VerificarCamposMetaRegisto(meta map[string]interface{}) error {
	// Verificar os campos base da info do registo
	campos := []string{
		"tipo",
		"estado",
		"quantidade",
	}
	// Verifica os campos base
	if err := VerificarCamposMapa(campos, meta); err != nil {
		return err
	}

	// Verifica se foi fornecido a quantidade em inventário do item
	if reflect.ValueOf(meta[campos[2]]).IsZero() {
		return errors.New("o registo deve fornecer uma quantidade miníma de items existentes")
	}

	// Define o estado do item em Upper, evita estados iguais escritos de maneiras diferentes
	meta[campos[1]] = strings.ToUpper(meta[campos[1]].(string))

	// Define o tipo do item em Upper, evita tipos iguais escritos de maneiras diferentes
	meta[campos[0]] = strings.ToUpper(meta[campos[0]].(string))

	// Define o valor minimo do item , evita quantidades inferiores a zero
	if meta[campos[2]].(float64) < 0 {
		meta[campos[2]] = 0
	}

	return nil
}

func RunQuerysOnRecords(querys resolvedschema.Query, registos []resolvedschema.Registo) []map[string]interface{} {
	var records = make([]map[string]interface{}, 0)
	fmt.Println("registos:", registos)

	if len(registos) > len(querys.Extrair) {
		for _, registo := range registos {
			// Junta todos os valores do registo para poderem ser "queryed"
			registoCompleto := JuntarPropriedadesDeRegisto(registo)

			// Mapa temporário a ser usado para extrair os valores
			mapTemp := make(map[string]interface{})
			// Query a executar na busca atual
			queryCurrente := querys.Extrair[0]

			// Verifica se o query está vazio, se sim
			// A extrasão dos elementos do registo, irá devolver todos os valores
			if len(queryCurrente) < 1 {
				for k := range registoCompleto {
					queryCurrente = append(queryCurrente, k)
				}
			}
			// Extrai os valores de JSON existentes na struct
			ExtractValuesFromJSON(queryCurrente, registoCompleto, mapTemp)

			// Se chegarmos ao ultimo query na lista, aplica esse mesmo a todos os registo restantes
			if len(querys.Extrair) >= 2 {
				querys.Extrair = querys.Extrair[1:]
			}
			// Adiciona os registos ao retorno
			records = append(records, mapTemp)
		}
		fmt.Println(records)
		return records
	}

	for _, registo := range registos {
		// Mapa temporário a ser usado para extrair os valores
		mapTemp := make(map[string]interface{})

		// Query a executar na busca atual
		queryCurrente := querys.Extrair[0]

		// Junta todos os valores do registo para poderem ser "queryed"
		registoCompleto := JuntarPropriedadesDeRegisto(registo)

		// Verifica se o query está vazio, se sim
		// A extrasão dos elementos do registo, irá devolver todos os valores
		if len(queryCurrente) <= 1 {
			for k := range registoCompleto {
				queryCurrente = append(queryCurrente, k)
			}
		}

		ExtractValuesFromJSON(queryCurrente, registo.Body, mapTemp)
		records = append(records, mapTemp)
	}
	return records
}

func JuntarPropriedadesDeRegisto(registo resolvedschema.Registo) map[string]interface{} {
	// Adiciona as propriedades estáticas da struct ao mapa
	registoCompleto := map[string]interface{}{
		"quantidade": registo.Meta.Quantidade,
		"estado":     registo.Meta.Estado,
		"tipo":       registo.Meta.Tipo,
	}

	// Itera sobre todos os valores do body (não defenidos/dinamicos)
	// E adiciona os mesmos ao mapa a retornar
	for chave, valor := range registo.Body {
		registoCompleto[chave] = valor
	}
	return registoCompleto
}
