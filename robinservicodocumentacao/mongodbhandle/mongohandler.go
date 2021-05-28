package mongodbhandle

import (
	"context"
	"time"

	"github.com/tomascpmarques/PAP/backend/robinservicodocumentacao/loggers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InsserirUmRegisto :
//	 Insser um registo na coleção fornecida, a bd tem de ser defenida préviamente
func InsserirUmRegisto(registo interface{}, colecao *mongo.Collection, inssertTimeOut int) (*mongo.InsertOneResult, error) {
	cntx, cancel := context.WithTimeout(context.Background(), time.Duration(inssertTimeOut)*time.Second)

	// Insser o registo na base de dados
	index, err := colecao.InsertOne(cntx, registo, options.InsertOne())
	defer cancel()

	// Error handeling
	if err != nil {
		loggers.MongoDBLogger.Println("Erro: ", err)
		return nil, err
	}

	return index, err
}
