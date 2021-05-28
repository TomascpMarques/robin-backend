package endpointfuncs

import (
	"context"
	"fmt"
	"time"

	"github.com/tomascpmarques/PAP/backend/robinservicoequipamento/loggers"
	"github.com/tomascpmarques/PAP/backend/robinservicoequipamento/mongodbhandle"
	"github.com/tomascpmarques/PAP/backend/robinservicoequipamento/resolvedschema"

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
	record, err := mongodbhandle.InserirUmRegisto(item, mongoCollection, 10)

	if err != nil {
		loggers.ServerErrorLogger.Println("Error: ", err)
		result["Error"] = err
		return nil
	}
	result["resultado"] = record

	loggers.MongoDBLogger.Println("Registo inserido!")
	return result
}

// QueryRegistoJSON Executa um query nos registos encontrados que satisfazêm o filtro de pesquisa
// Devolve só os campos pedidos dos registos encontrados, no formato { "key1.key2.key3.value1": result1 }
func QueryRegistoJSON(campos map[string]interface{}, dbCollPar map[string]interface{}, token string) (result map[string]interface{}) {
	result = make(map[string]interface{})

	// Define o query a usar nas buscas e a colecao alvo
	query := resolvedschema.QueryParaStruct(&campos)
	colecao := GetColecaoFromDB(dbCollPar)

	// Busca os registos da coleção, que igualêm os resultados
	registos, err := GetRegistosDaColecao(query.Campos, colecao)
	if err != nil {
		loggers.ServerErrorLogger.Println("Error: ", err)
		result["Error"] = err
		return
	}

	// Extrai os campos pedidos
	var records []map[string]interface{}
	for k, registo := range registos {
		// Mapa temporário a ser usado para extrair os valores
		mapTemp := make(map[string]interface{})
		ExtractValuesFromJSON(query.Extrair[k], registo, mapTemp)
		records = append(records, mapTemp)
	}

	loggers.ResolverLogger.Println("Sucesso, campos extraidos com sucesso!")
	result["queryRes"] = records
	return
}

// ApagarRegistoPorID Apaga um registo pelo seu ObjectID, na bd e coleção fornecida
func ApagarRegistoPorID(dbCollPar map[string]interface{}, idItem string, token string) map[string]interface{} {
	result := make(map[string]interface{})

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
		result["erro"] = err
		return result
	}

	// Set-up do filtro
	filter := bson.M{"_id": idOBJ}

	// Get collection
	coll := MongoClient.Database(dbCollPar["db"].(string)).Collection(dbCollPar["cl"].(string))
	cntx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	// Atualiza o item através do map especificado nos params
	matchCount, err := coll.UpdateOne(cntx, filter, bson.M{"$set": item}, options.MergeUpdateOptions())
	fmt.Println(matchCount)
	defer cancel()
	if err != nil {
		loggers.ServerErrorLogger.Println("Erro: ", err, " | registo de id: ", idItem)
		result["erro"] = "Error: registo:" + idItem
		return result
	}

	result["atualizacoes"] = matchCount
	return result
}
