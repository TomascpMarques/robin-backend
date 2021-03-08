package mongodbhandle

import (
	"context"
	"time"

	"github.com/tomascpmarques/PAP/backend/robinservicoequipamento/loggers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InsserirUmRegisto -
func InsserirUmRegisto(registo interface{}, colecao *mongo.Collection, inssertTimeOut int) (*mongo.InsertOneResult, error) {
	cntx, cancel := context.WithTimeout(context.Background(), time.Duration(inssertTimeOut)*time.Second)

	index, err := colecao.InsertOne(cntx, registo, options.InsertOne())
	defer cancel()

	if err != nil {
		loggers.MongoDBLogger.Println("Erro: ", err)
		return nil, err
	}

	return index, err
}
