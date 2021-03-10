package endpointfuncs

import (
	"context"
	"fmt"
	"time"

	"github.com/tomascpmarques/PAP/backend/robinservicoequipamento/loggers"
	"github.com/tomascpmarques/PAP/backend/robinservicoequipamento/mongodbhandle"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoParams = mongodbhandle.MongoConexaoParams{
	URI: "mongodb://0.0.0.0:27018",
}

// MongoClient -
var MongoClient = mongodbhandle.CriarConexaoMongoDB(mongoParams)

type a struct {
	ID         string
	Nome       string
	Info       Info
	Manutencao Manutencao
}

// Info -
type Info struct {
	Sala  int
	Piso  int
	Notas string
}

// Manutencao -
type Manutencao struct {
	Status           string
	UltimaManutencao string
}

// Hello - HELLOs back
func Hello(str string) map[string]interface{} {
	res := make(map[string]interface{})

	mongoDB := mongodbhandle.GetMongoDatabase(MongoClient, "local")
	mongoDBCollection := mongodbhandle.GetMongoCollection(mongoDB, "startup_log")

	filter := bson.M{"id": "PC1"}
	var temp a
	err := mongoDBCollection.FindOne(context.Background(), filter, options.FindOne()).Decode(&temp)
	if err != nil {
		println("Error: ", err)
		res["Error"] = err
		return nil
	}
	fmt.Println(temp)

	res["Result"] = temp

	return res
}

// AdicionarRegisto -
func AdicionarRegisto(tipoDeIndex string, item map[string]interface{}, token string) map[string]interface{} {
	result := make(map[string]interface{}, 0)

	// if VerificarTokenUser(token) != "OK" {
	// 	fmt.Println("Erro: A token fornecida é inválida ou expirou")
	// 	result["erro"] = "A token fornecida é inválida ou expirou"
	// 	return result
	// }

	// Declara o tipo de registo para outras funções terem o tipo de dados necessários
	// Para apontarem para structs compativeis
	item["tipo_de_registo"] = tipoDeIndex

	mongoCollection := mongodbhandle.GetMongoDatabase(MongoClient, "testing").Collection("base_collection")

	record, err := mongodbhandle.InsserirUmRegisto(item, mongoCollection, 10)

	if err != nil {
		fmt.Println("Error: ", err)
		result["Error"] = err
		return nil
	}
	result["Result"] = record

	loggers.MongoDBLogger.Println("Registo ensserido!")
	return result
}

// BuscarRegistoPorObjID Busca um registo na base de dados pelo ID especificado
func BuscarRegistoPorObjID(id string, token string) map[string]interface{} {
	result := make(map[string]interface{}, 0)

	// if VerificarTokenUser(token) != "OK" {
	// 	fmt.Println("Erro: A token fornecida é inválida ou expirou")
	// 	result["erro"] = "A token fornecida é inválida ou expirou"
	// 	return result
	// }

	// Converte o ID de uma String para um ObjectID
	idOBJ, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println("Error:  encontrado")
		result["Erro"] = err
		return result
	}

	// Setup do filtro e coleção a usar
	filter := bson.M{"_id": idOBJ}
	collection := MongoClient.Database("testing").Collection("base_collection")

	// Var temporária para guardar o valor recebido da base de dados
	var target map[string]interface{}
	cntx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	// Procura na coleção o registo com o ID igual
	res := collection.FindOne(cntx, filter).Decode(&target)
	defer cancel()

	// Converte o registo de um map[string]interface{} para a struct adequada
	registoStruct := mongodbhandle.ParseTipoDeRegisto(target)

	if res != nil {
		fmt.Println("Error: Registo não encontrado")
		result["Erro"] = "Registo não encontrado!"
		return result
	}

	loggers.DbFuncsLogger.Println("Registo Encontrado, pronto a enviar...")
	result["Result"] = registoStruct
	result["_id"] = target["_id"]
	result["tipo_de_registo"] = target["tipo_de_registo"]

	return result
}

// BuscarRegistosCustomQuery -
func BuscarRegistosCustomQuery(query map[string]interface{}, token string) map[string]interface{} {
	result := make(map[string]interface{}, 0)

	// if VerificarTokenUser(token) != "OK" {
	// 	fmt.Println("Erro: A token fornecida é inválida ou expirou")
	// 	result["erro"] = "A token fornecida é inválida ou expirou"
	// 	return result
	// }

	bsonFilter := make(bson.M, 0)
	bsonFilter = query

	collection := MongoClient.Database("testing").Collection("base_collection")

	var temp map[string]interface{}
	err := collection.FindOne(context.Background(), bsonFilter, options.FindOne()).Decode(&temp)
	if err != nil {
		result["Erro"] = err
		return result
	}

	result["Resultado"] = temp
	return result
}

// ApagarRegistoDeItem -
func ApagarRegistoDeItem(idItem string, token string) map[string]interface{} {
	result := make(map[string]interface{}, 0)

	// if VerificarTokenUser(token) != "OK" {
	// 	fmt.Println("Erro: A token fornecida é inválida ou expirou")
	// 	result["erro"] = "A token fornecida é inválida ou expirou"
	// 	return result
	// }

	return result
}

// AtualizararRegistoDeItem -
func AtualizararRegistoDeItem() {}

// BuscarInfoItems -
func BuscarInfoItems() {}

// BuscarInfoItem -
func BuscarInfoItem() {}
