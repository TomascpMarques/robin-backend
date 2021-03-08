package mongodbhandle

import (
	"context"
	"time"

	"github.com/tomascpmarques/PAP/backend/robinservicoequipamento/loggers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoConexaoParams -
type MongoConexaoParams struct {
	Ctx    context.Context
	Cancel context.CancelFunc
	URI    string
}

// mongoCtx       - contexto para a conexão ao serviço mongodb
// mongoCtxCancel - Contexto cancel para usar na conexão ao serviço mongodb
var mongoCtx, mongoCtxCancel = context.WithTimeout(context.Background(), 10*time.Second)

// CriarConexaoMongoDB -
func CriarConexaoMongoDB(params MongoConexaoParams) *mongo.Client {
	if params.Ctx == nil {
		params.Ctx = mongoCtx
	}
	if params.Cancel == nil {
		params.Cancel = mongoCtxCancel
	}
	if params.URI == "" {
		params.URI = "mongodb://localhost:27017"
	}

	client, err := mongo.Connect(params.Ctx, options.Client().ApplyURI(params.URI))
	if err != nil {
		panic(err)
	}
	loggers.MongoDBLogger.Println("Cliente MongoDB criado!")

	err = CheckConexaoMongo(params.Ctx, client, params.Cancel)
	if err != nil {
		panic(err)
	}
	loggers.MongoDBLogger.Println("Ping com sucesso, DB está UP")

	return client
}

// CheckConexaoMongo -
func CheckConexaoMongo(ctx context.Context, client *mongo.Client, cancelFunc context.CancelFunc) error {
	err := client.Ping(ctx, readpref.Primary())
	defer cancelFunc()
	return err
}

// GetMongoDatabase -
func GetMongoDatabase(cl *mongo.Client, dbName string) *mongo.Database {
	return cl.Database(dbName)
}

// GetMongoCollection -
func GetMongoCollection(db *mongo.Database, collName string) *mongo.Collection {
	return db.Collection(collName)
}
