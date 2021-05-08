package endpointfuncs

import (
	"context"
	"time"

	"github.com/tomascpmarques/PAP/backend/robinservicodocumentacao/resolvedschema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetRepoPorCampo Busca um repo e devolveo na struct resolvedschema.Repositorio
// Busca o repositório através de um campo e valor do mesmo, especificado na sua chamada
func GetRepoPorCampo(campo string, valor interface{}) (repo resolvedschema.Repositorio) {
	// Define o filtro a usar na procura de informação na BD
	filter := bson.M{campo: valor}
	// Documento e repo onde procurar o repo
	collection := MongoClient.Database("documentacao").Collection("repos")
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
