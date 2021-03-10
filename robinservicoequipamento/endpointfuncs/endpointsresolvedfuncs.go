package endpointfuncs

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/tomascpmarques/PAP/backend/robinservicoequipamento/loggers"
	"github.com/tomascpmarques/PAP/backend/robinservicoequipamento/mongodbhandle"
	"github.com/tomascpmarques/PAP/backend/robinservicoequipamento/structextract"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoParams = mongodbhandle.MongoConexaoParams{
	URI: "mongodb://0.0.0.0:27018",
}

// MongoClient -
var MongoClient = mongodbhandle.CriarConexaoMongoDB(mongoParams)

// AdicionarRegisto Adiciona um registo numa base de dados e coleção especifícada
func AdicionarRegisto(tipoDeIndex string, dbCollPar map[string]interface{}, item map[string]interface{}, token string) map[string]interface{} {
	result := make(map[string]interface{}, 0)

	// if VerificarTokenUser(token) != "OK" {
	// 	fmt.Println("Erro: A token fornecida é inválida ou expirou")
	// 	result["erro"] = "A token fornecida é inválida ou expirou"
	// 	return result
	// }

	// Declara o tipo de registo para outras funções terem o tipo de dados necessários
	// Para apontarem para structs compativeis
	item["tipo_de_registo"] = tipoDeIndex

	// Get the mongo colection
	mongoCollection := MongoClient.Database(dbCollPar["db"].(string)).Collection(dbCollPar["cl"].(string))

	// Insser um registo na coleção e base de dados especificada
	record, err := mongodbhandle.InsserirUmRegisto(item, mongoCollection, 10)

	if err != nil {
		loggers.ServerErrorLogger.Println("Error: ", err)
		result["Error"] = err
		return nil
	}
	result["resultado"] = record

	loggers.MongoDBLogger.Println("Registo ensserido!")
	return result
}

// BuscarRegistoPorObjID Busca um registo na base de dados pelo ID especificado
func BuscarRegistoPorObjID(dbCollPar map[string]interface{}, id string, token string) map[string]interface{} {
	result := make(map[string]interface{}, 0)

	// if VerificarTokenUser(token) != "OK" {
	// 	fmt.Println("Erro: A token fornecida é inválida ou expirou")
	// 	result["erro"] = "A token fornecida é inválida ou expirou"
	// 	return result
	// }

	// Converte o ID de uma String para um ObjectID
	idOBJ, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		loggers.ServerErrorLogger.Println()
		fmt.Println("Error: ID de registo não pode ser convertido")
		result["Erro"] = err
		return result
	}

	// Setup do filtro e coleção a usar
	filter := bson.M{"_id": idOBJ}
	collection := MongoClient.Database(dbCollPar["db"].(string)).Collection(dbCollPar["cl"].(string))

	// Var temporária para guardar o valor recebido da base de dados
	var target map[string]interface{}
	cntx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	// Procura na coleção o registo com o ID igual
	res := collection.FindOne(cntx, filter).Decode(&target)
	defer cancel()
	if res != nil {
		loggers.ServerErrorLogger.Println("Error: Registo não encontrado")
		result["Erro"] = "Registo não encontrado!"
		return result
	}

	// Converte o registo de um map[string]interface{} para a struct adequada
	registoStruct := mongodbhandle.ParseTipoDeRegisto(target)
	if registoStruct == nil {
		loggers.ServerErrorLogger.Println("Error: Ao cinverter o registo")
		result["Erro"] = "Impossível de converter!"
		return result
	}

	loggers.DbFuncsLogger.Println("Registo Encontrado, pronto a enviar...")
	result["resultado"] = registoStruct
	result["meta_data"] = map[string]interface{}{
		"id":              target["_id"],
		"tipo_de_registo": target["tipo_de_registo"],
	}

	return result
}

// BuscarRegistosCamposCustom :
//	Toma um nome de uma bd e uma coleção como alvos do query.
// 	O query em sí é um map, que vai fornecer os valores ao filtro do tipo bson.M.
//	Toma uma token para autorização
func BuscarRegistosCamposCustom(dbCollPar map[string]interface{}, query map[string]interface{}, token string) map[string]interface{} {
	result := make(map[string]interface{}, 0)

	// if VerificarTokenUser(token) != "OK" {
	// 	fmt.Println("Erro: A token fornecida é inválida ou expirou")
	// 	result["erro"] = "A token fornecida é inválida ou expirou"
	// 	return result
	// }

	// Get collection da db fornecida
	coll := MongoClient.Database(dbCollPar["db"].(string)).Collection(dbCollPar["cl"].(string))

	// Chama a função que procura todos os registos, validados pelo query
	temp := mongodbhandle.PesquisaComQueryCustom(coll, query)
	if temp == nil {
		result["erro"] = "Valor não encontrado para parametros:"
		result["parametros"] = query
		return result
	}

	// O retorno são esses valores traduzidos para as suas estruturas correspondentes
	result["resultado"] = mongodbhandle.MongoRecordsParssedArrays(temp)
	return result
}

// BuscarInfoItem -
func BuscarInfoItem(dbCollPar map[string]interface{}, id string, query map[string]interface{}, token string) map[string]interface{} {
	result := make(map[string]interface{}, 0)

	// if VerificarTokenUser(token) != "OK" {
	// 	fmt.Println("Erro: A token fornecida é inválida ou expirou")
	// 	result["erro"] = "A token fornecida é inválida ou expirou"
	// 	return result
	// }

	// Converte o ID de uma String para um ObjectID
	idOBJ, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		loggers.ServerErrorLogger.Println()
		fmt.Println("Error: ID de registo não pode ser convertido")
		result["Erro"] = err
		return result
	}

	bsonFilter := bson.M{"_id": idOBJ}

	// Get collection da db fornecida
	coll := MongoClient.Database(dbCollPar["db"].(string)).Collection(dbCollPar["cl"].(string))

	var temp map[string]interface{}
	cntx, cancel := mongodbhandle.MongoCtxMaker("bg", time.Duration(10))

	err = coll.FindOne(cntx, bsonFilter, options.FindOne()).Decode(&temp)
	defer cancel()
	if err != nil {
		loggers.ServerErrorLogger.Println("Error: Registo não encontrado para id fornecido")
		result["Erro"] = "Registo não encontrado para itemID: " + id
		return result
	}

	retStruct := mongodbhandle.ParseTipoDeRegisto(temp)
	if retStruct == nil {
		loggers.ServerErrorLogger.Println("Error: Não foi possível converter o registo")
		result["Erro"] = "Não foi possível converter o registo"
		return result
	}

	xx := make(map[string][]string)
	fmt.Println(query)
	res, _ := json.Marshal(query)
	_ = json.Unmarshal(res, &xx)

	x := structextract.ExtrairCamposEspecificosStruct(retStruct, xx)

	result["Registo"] = x
	return result
}

// ApagarRegistoDeItem :
// 	Apaga um registo pelo seu ObjectID, na bd e coleção fornecida
func ApagarRegistoDeItem(dbCollPar map[string]interface{}, idItem string, token string) map[string]interface{} {
	result := make(map[string]interface{}, 0)

	// if VerificarTokenUser(token) != "OK" {
	// 	fmt.Println("Erro: A token fornecida é inválida ou expirou")
	// 	result["erro"] = "A token fornecida é inválida ou expirou"
	// 	return result
	// }

	// Converte o ID de uma String para um ObjectID
	idOBJ, err := primitive.ObjectIDFromHex(idItem)
	if err != nil {
		loggers.ServerErrorLogger.Println()
		fmt.Println("Error: ID de registo não pode ser convertido")
		result["Erro"] = err
		return result
	}

	// Set-up do filtro
	filter := bson.M{"_id": idOBJ}

	// Get collection
	coll := MongoClient.Database(dbCollPar["db"].(string)).Collection(dbCollPar["cl"].(string))
	cntx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	// Operação de delete de um só registo
	item, err := coll.DeleteOne(cntx, filter, options.Delete())
	defer cancel()
	if err != nil {
		loggers.ServerErrorLogger.Println()
		fmt.Println("Erro: Não foi possível apagar o registo de id: ", idItem)
		result["Erro"] = "Error: Não foi possível apagar o item de registo:" + idItem
		return result
	}

	result["registo_apagado"] = item
	return result
}

// AtualizararRegistoDeItem -
func AtualizararRegistoDeItem() {}

// BuscarInfoItems -
func BuscarInfoItems() {}
