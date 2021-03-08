package endpointfuncs

import (
	"context"
	"fmt"

	"github.com/tomascpmarques/PAP/backend/robinservicoequipamento/mongodbhandle"
	"go.mongodb.org/mongo-driver/bson"
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
	_ = mongoDBCollection.FindOne(context.Background(), filter, options.FindOne()).Decode(&temp)
	fmt.Println(temp)

	res["Result"] = temp

	return res
}

// AdicionarRegisto -
func AdicionarRegisto(tipoItem string, item []byte, token string) map[string]interface{} {
	result := make(map[string]interface{}, 0)

	// if VerificarTokenUser(token) != "OK" {
	// 	fmt.Println("Erro: A token fornecida é inválida ou expirou")
	// 	result["erro"] = "A token fornecida é inválida ou expirou"
	// 	return result
	// }

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
