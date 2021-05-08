package repos

import (
	"context"
	"errors"
	"reflect"
	"time"

	"github.com/tomascpmarques/PAP/backend/robinservicodocumentacao/endpointfuncs"
	"github.com/tomascpmarques/PAP/backend/robinservicodocumentacao/loggers"
	"github.com/tomascpmarques/PAP/backend/robinservicodocumentacao/resolvedschema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetRepoPorCampo Busca um repo e devolveo na struct resolvedschema.Repositorio
// Busca o repositório através de um campo e valor do mesmo, especificado na sua chamada
func GetRepoPorCampo(campo string, valor interface{}) (repo resolvedschema.Repositorio) {
	// Define o filtro a usar na procura de informação na BD
	filter := bson.M{campo: valor}
	// Documento e repo onde procurar o repo
	collection := endpointfuncs.MongoClient.Database("documentacao").Collection("repos")
	cntx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	// Procura na BD do registo pedido
	err := collection.FindOne(cntx, filter, options.FindOne()).Decode(&repo)
	defer cancel()
	if err != nil {
		// Devolve um repo vzaio se não se encontrar o pedido
		repo = resolvedschema.Repositorio{}
		return
	}

	// Devolve repo
	return
}

func DropRepoPorNome(repoNome string) (erro error) {
	// Define o filtro a usar na procura de informação na BD
	filter := bson.M{"nome": repoNome}
	// Documento e repo onde procurar o repo
	collection := endpointfuncs.MongoClient.Database("documentacao").Collection("repos")
	cntx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	// Depois de apagar o registo, a var err,
	// Vai ter o sucesso ou falhanço da operação como o seu valor
	_, erro = collection.DeleteOne(cntx, filter)
	defer cancel()

	return
}

func UpdateRepositorioPorNome(repoName string, mundancas map[string]interface{}) *mongo.UpdateResult {
	// Set-up do filtro
	filter := bson.M{"nome": repoName}

	// Get collection
	coll := endpointfuncs.MongoClient.Database("documentacao").Collection("repos")
	cntx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	// Atualiza o item através do map especificado nos params
	matchCount, err := coll.UpdateOne(cntx, filter, mundancas, options.MergeUpdateOptions())
	defer cancel()
	if err != nil {
		loggers.DbFuncsLogger.Println("Erro ao atualizar o registo")
		return nil
	}

	return matchCount
}

func VerificarInfoBaseRepo(info map[string]interface{}) (err error) {
	err = nil
	keysObrg := []string{
		"nome",
		"autor",
		"tema",
	}
	for _, v := range keysObrg {
		if valor, existe := info[v]; !(reflect.ValueOf(valor).IsValid()) || !existe {
			err = errors.New("não cumpre os parametros")
			break
		}
	}
	return
}
