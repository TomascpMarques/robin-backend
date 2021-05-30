package endpointfuncs

import (
	"fmt"

	"github.com/tomascpmarques/PAP/backend/robinservicodocumentacao/loggers"
	"github.com/tomascpmarques/PAP/backend/robinservicodocumentacao/mongodbhandle"
)

var mongoParams = mongodbhandle.MongoConexaoParams{
	URI: "mongodb://0.0.0.0:27020/",
}

// MongoClient cliente com a conexão à instancia mongo
var MongoClient = mongodbhandle.CriarConexaoMongoDB(mongoParams)

// PingServico responde que o serviço está online
func PingServico(name string) (retorno map[string]interface{}) {
	retorno = make(map[string]interface{})
	retorno["status"] = fmt.Sprintf("Hello %s, I'm alive and OK", name)
	return
}

func TesteInvalidAuthTokenReload(token string, reloadToken string) (retorno map[string]interface{}) {
	retorno = make(map[string]interface{})

	if VerificarTokenUser(token) != "OK" {
		if err := VerificarTokenReAuth(reloadToken, token); err != "OK" {
			loggers.ServerErrorLogger.Println("Erro: A token de revalidação do user também está inválida")
			retorno["erro"] = "Por favor volta a iniciar sessão"
			return
		}
		retorno["sucesso"] = "Olá"
		return
	}

	return
}
