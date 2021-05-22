package endpointfuncs

import (
	"fmt"
	"reflect"

	"github.com/tomascpmarques/PAP/backend/robinservicovideoshare/loggers"
	"github.com/tomascpmarques/PAP/backend/robinservicovideoshare/mongodbhandle"
	"github.com/tomascpmarques/PAP/backend/robinservicovideoshare/resolvedschema"
)

var mongoParams = mongodbhandle.MongoConexaoParams{
	URI: "mongodb://0.0.0.0:27022/",
}

// MongoClient cliente com a conexão à instancia mongo
var MongoClient = mongodbhandle.CriarConexaoMongoDB(mongoParams)

// PingServico responde que o serviço está online
func PingServico(name string) (retorno map[string]interface{}) {
	retorno = make(map[string]interface{})
	retorno["status"] = fmt.Sprintf("Hello %s, I'm alive and OK", name)
	return
}

func CriarVideoShare(videoMetaData map[string]interface{}, token string) (retorno map[string]interface{}) {
	retorno = make(map[string]interface{})

	if VerificarTokenUser(token) != "OK" {
		loggers.ResolverLogger.Println("Token inválida.")
		retorno["err"] = "Token inválida ou expirada"
		return
	}

	videoSahre := resolvedschema.VideoParaStruct(&videoMetaData)
	if reflect.ValueOf(videoSahre).IsZero() {
		loggers.ResolverLogger.Println("Erro ao converter as informações dadas para uma struct.")
		retorno["err"] = "Erro ao converter as informações dadas para uma struct"
		return
	}

	return
}
