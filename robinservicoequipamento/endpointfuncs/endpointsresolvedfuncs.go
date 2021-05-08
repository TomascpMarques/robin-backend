package endpointfuncs

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/tomascpmarques/PAP/backend/robinservicoequipamento/loggers"
	"github.com/tomascpmarques/PAP/backend/robinservicoequipamento/mongodbhandle"
	"github.com/tomascpmarques/PAP/backend/robinservicoequipamento/structextract"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoParams = mongodbhandle.MongoConexaoParams{
	URI: "mongodb://0.0.0.0:27018/",
}

// MongoClient cliente com a conexão à instancia mongo
var MongoClient = mongodbhandle.CriarConexaoMongoDB(mongoParams)

// PingServico responde que o serviço está online
func PingServico(name string) map[string]interface{} {
	result := make(map[string]interface{})

	result["status"] = fmt.Sprintf("Hello %s, I'm alive and OK", name)
	return result
}

// AdicionarRegisto Adiciona um registo numa base de dados e coleção especifícada
func AdicionarRegisto(tipoDeIndex string, dbCollPar map[string]interface{}, item map[string]interface{}, token string) map[string]interface{} {
	result := make(map[string]interface{})

	if VerificarTokenUser(token) != "OK" {
		fmt.Println("Erro: A token fornecida é inválida ou expirou")
		result["erro"] = "A token fornecida é inválida ou expirou"
		return result
	}

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

// QueryRegistoObjID Busca um registo na base de dados pelo ID especificado
func QueryRegistoObjID(dbCollPar map[string]interface{}, id string, token string) map[string]interface{} {
	result := make(map[string]interface{})

	if VerificarTokenUser(token) != "OK" {
		fmt.Println("Erro: A token fornecida é inválida ou expirou")
		result["erro"] = "A token fornecida é inválida ou expirou"
		return result
	}

	// Converte o ID de uma String para um ObjectID
	idOBJ, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		loggers.ServerErrorLogger.Println()
		fmt.Println("Error: ID de registo não pode ser convertido")
		result["erro"] = err
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
		result["erro"] = "Registo não encontrado!"
		return result
	}

	// Converte o registo de um map[string]interface{} para a struct adequada
	registoStruct := mongodbhandle.ParseTipoDeRegisto(target)
	if registoStruct == nil {
		loggers.ServerErrorLogger.Println("Error: Ao cinverter o registo")
		result["erro"] = "Impossível de converter!"
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

// QueryRegistosCustom Toma um nome de uma bd e uma coleção como alvos do query.
// O query em sí é um map, que vai fornecer os valores ao filtro do tipo bson.M.
// Toma uma token para autorização
func QueryRegistosCustom(dbCollPar map[string]interface{}, query map[string]interface{}, token string) map[string]interface{} {
	result := make(map[string]interface{})

	if VerificarTokenUser(token) != "OK" {
		fmt.Println("Erro: A token fornecida é inválida ou expirou")
		result["erro"] = "A token fornecida é inválida ou expirou"
		return result
	}

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

// BuscarInfoItemQuery Busca um registo e devolve os campos especificados no query
func BuscarInfoItemQuery(dbCollPar map[string]interface{}, id string, query map[string]interface{}, token string) map[string]interface{} {
	result := make(map[string]interface{})

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
		result["erro"] = err
		return result
	}

	// Bson filter setup
	bsonFilter := bson.M{"_id": idOBJ}

	// Get collection da db fornecida
	coll := MongoClient.Database(dbCollPar["db"].(string)).Collection(dbCollPar["cl"].(string))

	// Resultado temporário, e contexto/cancel para uso no MongoClient
	var temp map[string]interface{}
	cntx, cancel := mongodbhandle.MongoCtxMaker("bg", time.Duration(10))

	// Procura um registo
	err = coll.FindOne(cntx, bsonFilter, options.FindOne()).Decode(&temp)
	defer cancel()
	if err != nil {
		loggers.ServerErrorLogger.Println("Error: Registo não encontrado para id fornecido")
		result["erro"] = "Registo não encontrado para itemID: " + id
		return result
	}

	// Traduz o registo de um mpa para a struct equivalente
	retStruct := mongodbhandle.ParseTipoDeRegisto(temp)
	if retStruct == nil {
		loggers.ServerErrorLogger.Println("Error: Não foi possível converter o registo")
		result["erro"] = "Não foi possível converter o registo"
		return result
	}

	/* Iníco da conversão de um map[string]interface{} para um map[string][]string */
	mapConvertido := make(map[string][]string)
	// Tradução do query para json
	res, err := json.Marshal(query)
	if err != nil {
		loggers.ServerErrorLogger.Println("Error: ", err)
		result["erro"] = "err"
		return result
	}

	// Tradução do valor em bytes do query para o mapa que se quer
	err = json.Unmarshal(res, &mapConvertido)
	if err != nil {
		loggers.ServerErrorLogger.Println("Error: ", err)
		result["erro"] = "err"
		return result
	}
	/* Fim da conversão de mapas */

	// Extráção dos campos
	temp = structextract.ExtrairCamposEspecificosStruct(retStruct, mapConvertido)
	if temp == nil {
		loggers.ServerErrorLogger.Println("Error: Não foi possível extrair os valores pedidos")
		result["erro"] = "Não foi possível extrair os valores pedidos"
		return result
	}

	result["registo"] = temp
	return result
}

// BuscarInfoItems Busca vários registos, com campos que satisfazem o queryBD e devolve os campos especificados no query
func BuscarInfoItems(dbCollPar map[string]interface{}, queryDB map[string]interface{}, query map[string]interface{}, token string) (result map[string]interface{}) {
	result = make(map[string]interface{})

	if VerificarTokenUser(token) != "OK" {
		fmt.Println("Erro: A token fornecida é inválida ou expirou")
		result["erro"] = "A token fornecida é inválida ou expirou"
		return result
	}

	// Get collection da db fornecida
	coll := MongoClient.Database(dbCollPar["db"].(string)).Collection(dbCollPar["cl"].(string))

	// Resultado temporário, e contexto/cancel para uso no MongoClient
	var temp []map[string]interface{}
	cntx, cancel := mongodbhandle.MongoCtxMaker("bg", time.Duration(10))

	// Procura por todos os registos que validam o query
	cursor, err := coll.Find(cntx, queryDB, options.Find())
	defer cancel()
	if err != nil {
		loggers.DbFuncsLogger.Println("Erro ao procurar os registos")
		result["erro"] = "Erro ao procurar os registos"
		return
	}

	// Transfere os valores do curssor mongoDB para um array de map[string]interface{}
	if err := cursor.All(cntx, &temp); err != nil {
		loggers.DbFuncsLogger.Println("erro a decodificar valores bson para interface{}")
		result["erro"] = "Erro a decodificar valores bson para interface{}"
		return
	}

	/* Iníco da conversão de um map[string]interface{} para um map[string][]string */
	mapConvertido := make(map[string][]string)
	// Tradução do query para json
	res, err := json.Marshal(query)
	if err != nil {
		loggers.ServerErrorLogger.Println("Error: ", err)
		result["erro"] = "err"
		return result
	}

	// Tradução do valor em bytes do query para o mapa que se quer
	err = json.Unmarshal(res, &mapConvertido)
	if err != nil {
		loggers.ServerErrorLogger.Println("Error: ", err)
		result["erro"] = "err"
		return result
	}
	/* Fim da conversão de mapas */

	// Resultados da extração de campos do query
	resultadoQuery := make([]map[string]interface{}, len(temp))

	// Itera sobre todos os mapas retirados da BD
	for k, v := range temp {
		// Traduz o registo de um mapa para a struct equivalente
		registoEmEstrutura := mongodbhandle.ParseTipoDeRegisto(v)

		if registoEmEstrutura == nil {
			loggers.ServerErrorLogger.Println("Error: Não foi possível converter o registo")
			result["erro"] = "Não foi possível converter o registo"
			return result
		}

		// Cria o mapa para a estrutura currente e adiciona o id da estrutura
		resultadoQuery[k] = map[string]interface{}{
			"id": v["_id"],
		}

		// Extração dos campos através de uma lista custom
		temp := structextract.ExtrairCamposEspecificosStruct(registoEmEstrutura, mapConvertido)
		for _, value := range reflect.ValueOf(temp).MapKeys() {
			resultadoQuery[k][value.String()] = temp[value.String()]
			fmt.Println("<<", value, temp[value.String()])
		}
	}

	result["registos"] = resultadoQuery
	return
}

// ApagarRegistoPorID Apaga um registo pelo seu ObjectID, na bd e coleção fornecida
func ApagarRegistoPorID(dbCollPar map[string]interface{}, idItem string, token string) map[string]interface{} {
	result := make(map[string]interface{})

	if VerificarTokenUser(token) != "OK" {
		fmt.Println("Erro: A token fornecida é inválida ou expirou")
		result["erro"] = "A token fornecida é inválida ou expirou"
		return result
	}

	// Converte o ID de uma String para um ObjectID
	idOBJ, err := primitive.ObjectIDFromHex(idItem)
	if err != nil {
		loggers.ServerErrorLogger.Println()
		fmt.Println("Error: ID de registo não pode ser convertido")
		result["erro"] = err
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
		loggers.ServerErrorLogger.Println("Erro: Não foi possível apagar o registo de id: ", idItem)
		result["erro"] = "Error: Não foi possível apagar o item de registo:" + idItem
		return result
	}

	result["registo_apagado"] = item
	return result
}

// AtualizarRegistoDeItem Na bd e coleção escolhida, o registo de id idItem
// vai ser atualizado para os valores especificados em item
func AtualizarRegistoDeItem(dbCollPar map[string]interface{}, idItem string, item map[string]interface{}, token string) map[string]interface{} {
	result := make(map[string]interface{})

	if VerificarTokenUser(token) != "OK" {
		fmt.Println("Erro: A token fornecida é inválida ou expirou")
		result["erro"] = "A token fornecida é inválida ou expirou"
		return result
	}

	// Converte o ID de uma String para um ObjectID
	idOBJ, err := primitive.ObjectIDFromHex(idItem)
	if err != nil {
		loggers.ServerErrorLogger.Println()
		fmt.Println("Error: ID de registo não pode ser convertido")
		result["erro"] = err
		return result
	}

	// Set-up do filtro
	filter := bson.M{"_id": idOBJ}

	// Get collection
	coll := MongoClient.Database(dbCollPar["db"].(string)).Collection(dbCollPar["cl"].(string))
	cntx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	// Atualiza o item através do map especificado nos params
	matchCount, err := coll.UpdateOne(cntx, filter, item, options.MergeUpdateOptions())
	defer cancel()
	if err != nil {
		loggers.ServerErrorLogger.Println("Erro: ", err, " | registo de id: ", idItem)
		result["erro"] = "Error: registo:" + idItem
		return result
	}

	result["atualizacoes"] = matchCount
	return result
}
