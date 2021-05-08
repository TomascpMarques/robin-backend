package ficheiros

import (
	"context"
	"reflect"
	"time"

	"github.com/tomascpmarques/PAP/backend/robinservicodocumentacao/endpointfuncs"
	"github.com/tomascpmarques/PAP/backend/robinservicodocumentacao/loggers"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CriarFicheiro Cria a meta data de um ficheiro, para prepara o upload de conteúdo
func CriarFicheiro(ficheiroMetaData map[string]interface{}, token string) (retorno map[string]interface{}) {
	retorno = make(map[string]interface{})

	// if endpointfuncs.VerificarTokenUser(token) != "OK" {
	// 	loggers.OperacoesBDLogger.Println("Erro: A token fornecida é inválida ou expirou")
	// 	retorno["erro"] = "A token fornecida é inválida ou expirou"
	// 	return
	// }

	metaHash, err := CriarMetaHash(ficheiroMetaData)
	if err != nil {
		loggers.OperacoesBDLogger.Println("Erro ao criar hash para meta data: ", err)
		retorno["erro"] = "Erro ao criar hash para meta data fornecida"
		return
	}
	ficheiroMetaData["hash"] = metaHash
	if err := MetaDataBaseValida(ficheiroMetaData); err != nil {
		loggers.OperacoesBDLogger.Println(err.Error())
		retorno["erro"] = err.Error()
		return
	}
	ficheiroMetaData["criacao"] = time.Now().Unix()

	// Get the mongo colection
	mongoCollection := endpointfuncs.MongoClient.Database("documentacao").Collection("files-meta-data")
	cntx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	insserido, err := mongoCollection.InsertOne(cntx, ficheiroMetaData, options.InsertOne())
	defer cancel()
	if err != nil || !reflect.ValueOf(insserido.InsertedID).IsValid() {
		loggers.OperacoesBDLogger.Println("Erro ao insserir o registo na BD: ", err)
		retorno["erro"] = "Erro ao insserir o registo na BD"
		return
	}

	loggers.OperacoesBDLogger.Println("Meta Data insserida com sucesso")
	retorno["sucesso"] = true
	return
}
