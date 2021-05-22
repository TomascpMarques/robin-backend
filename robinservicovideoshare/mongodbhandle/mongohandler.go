package mongodbhandle

import (
	"context"
	"time"

	"github.com/tomascpmarques/PAP/backend/robinservicovideoshare/loggers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InserirUmRegisto :
//	 Insser um registo na coleção fornecida, a bd tem de ser defenida préviamente
func InserirUmRegisto(registo interface{}, colecao *mongo.Collection, inssertTimeOut int) (*mongo.InsertOneResult, error) {
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

// PesquisaComQueryCustom :
// 	Pesquisa e encontra todos os registos que satisfaçam as condições do query
func PesquisaComQueryCustom(collection *mongo.Collection, query map[string]interface{}) []map[string]interface{} {
	// Setup do filtro e atribuição
	bsonFilter := query

	// Busca e decoding do retorno
	var temp []map[string]interface{}
	cntx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	// collection.Find - Devolve um curssor com todos os valores encontrados
	curssor, err := collection.Find(cntx, bsonFilter, options.Find())
	defer cancelFunc()

	// Error handeling
	if err != nil {
		loggers.MongoDBLogger.Println("Erro: ", err)
		return nil
	}

	// Descodificar os registos para um []map[string]interface{}
	if err := curssor.All(context.TODO(), &temp); err != nil {
		loggers.MongoDBLogger.Println("Erro: ", err)
		return nil
	}

	return temp
}
