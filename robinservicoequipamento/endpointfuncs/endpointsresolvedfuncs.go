package endpointfuncs

import (
	"context"
	"fmt"
	"time"

	"github.com/tomascpmarques/PAP/backend/robinservicoequipamento/mongodbhandle"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoParams = mongodbhandle.MongoConexaoParams{
	Ctx:    context.Background(),
	Cancel: nil,
	URI:    "mongodb://0.0.0.0:27018",
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
func AdicionarRegisto(item map[string]interface{}, token string) map[string]interface{} {
	result := make(map[string]interface{}, 0)

	// if VerificarTokenUser(token) != "OK" {
	// 	fmt.Println("Erro: A token fornecida é inválida ou expirou")
	// 	result["erro"] = "A token fornecida é inválida ou expirou"
	// 	return result
	// }

	mongoCollection := mongodbhandle.GetMongoDatabase(MongoClient, "local").Collection("startup_log")

	record, err := mongodbhandle.InsserirUmRegisto(item, mongoCollection, 10)

	if err != nil {
		fmt.Println("Error: ", err)
		result["Error"] = err
		return nil
	}
	result["Result"] = record

	return result
}

// BuscarRegistoPorObjID -
func BuscarRegistoPorObjID(id string, token string) map[string]interface{} {
	result := make(map[string]interface{}, 0)

	// if VerificarTokenUser(token) != "OK" {
	// 	fmt.Println("Erro: A token fornecida é inválida ou expirou")
	// 	result["erro"] = "A token fornecida é inválida ou expirou"
	// 	return result
	// }

	idOBJ, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println("Error:  encontrado")
		result["Erro"] = err
		return result
	}

	filter := bson.M{"_id": idOBJ}
	collection := MongoClient.Database("local").Collection("startup_log")

	var target map[string]interface{}

	cntx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	res := collection.FindOne(cntx, filter).Decode(&target)
	defer cancel()

	if res != nil {
		fmt.Println("Error: Registo não encontrado")
		result["Erro"] = "Registo não encontrado!"
		return result
	}
	fmt.Println("Target: ", target)
	result["Result"] = target

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

// BuscarInfoDeItems -
func BuscarInfoDeItems() {}

// BuscarInfoDeItem -
func BuscarInfoDeItem() {}
