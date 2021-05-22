package endpointfuncs

import (
	"context"
	"time"

	"github.com/tomascpmarques/PAP/backend/robinservicovideoshare/resolvedschema"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoDBOperation Struct com o setup minímo para fazer uma oepração na BDs
type MongoDBOperation struct {
	Colecao    *mongo.Collection
	Cntxt      context.Context
	CancelFunc context.CancelFunc
	Filter     interface{}
}

// Setup Evita mais lihas desnecessárias e repetitivas para poder-se usar a coleção necessaria
func SetupColecao(dbName, collName string) (defs MongoDBOperation) {
	defs.Colecao = MongoClient.Database(dbName).Collection(collName)
	defs.Cntxt, defs.CancelFunc = context.WithTimeout(context.Background(), time.Second*10)
	return
}

// Adiciona os dados da video-share à base de dados
func AdicionarVideoShareDB(videoShare resolvedschema.Video) error {
	//colecao := SetupColecao("videoshares", "videos")

	return nil
}
