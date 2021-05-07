package mongodbhandle

import (
	"context"
	"time"

	"github.com/tomascpmarques/PAP/backend/robinservicodocumentacao/loggers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoConexaoParams - Parametros base para uma conexão á mongo bd
type MongoConexaoParams struct {
	URI string
}

// MongoCtxMaker -
func MongoCtxMaker(ctxTipo string, duracao time.Duration) (context.Context, context.CancelFunc) {
	if ctxTipo == "bg" {
		return context.WithTimeout(context.Background(), duracao*time.Second)
	}
	return context.WithTimeout(context.TODO(), duracao*time.Second)
}

// CriarConexaoMongoDB -
func CriarConexaoMongoDB(params MongoConexaoParams) *mongo.Client {
	// Verifica para valores default
	if params.URI == "" {
		params.URI = "mongodb://0.0.0.0:27020"
	}

	// Necessidades da conexão, i.e.: operation-timeout, operation context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	// Liga à instância mongo apontada pelos parametros
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(params.URI))
	defer cancel()
	if err != nil {
		panic(err)
	}
	loggers.MongoDBLogger.Println("Cliente MongoDB criado!")

	// Verifica a conexão ao mongoDB, antes de devolver o cliente
	err = CheckConexaoMongo(ctx, client, cancel)
	if err != nil {
		loggers.MongoDBLogger.Println("Erro: conexão não establecida à instância MongoDB, a saír. . .")
		panic(err)
	}
	loggers.MongoDBLogger.Println("Ping com sucesso, DB está UP")

	return client
}

// CheckConexaoMongo - Verifica a conexão à instância mongodb especificada, e se está alive
func CheckConexaoMongo(ctx context.Context, client *mongo.Client, cancelFunc context.CancelFunc) error {
	err := client.Ping(ctx, readpref.Primary())
	defer cancelFunc()
	return err
}
