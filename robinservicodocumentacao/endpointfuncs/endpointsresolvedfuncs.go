package endpointfuncs

import (
	"fmt"

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

// CriarRepositorio Cria um repo para guardar a informação relativa a um tema e/ou tarefa
func CriarRepositorio() (retorno map[string]interface{}) {
	retorno = make(map[string]interface{})

	return
}
