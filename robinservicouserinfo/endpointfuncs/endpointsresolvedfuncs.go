package endpointfuncs

import (
	"context"
	"fmt"
	"time"

	"github.com/tomascpmarques/PAP/backend/robinservicouserinfo/loggers"
	"github.com/tomascpmarques/PAP/backend/robinservicouserinfo/mongodbhandle"
	"github.com/tomascpmarques/PAP/backend/robinservicouserinfo/resolvedschema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoParams = mongodbhandle.MongoConexaoParams{
	URI: "mongodb://0.0.0.0:27019/",
}

// MongoClient cliente com a conexão à instancia mongo
var MongoClient = mongodbhandle.CriarConexaoMongoDB(mongoParams)

// PingServico responde que o serviço está online
func PingServico(name string) (result map[string]interface{}) {
	result = make(map[string]interface{})
	result["status"] = fmt.Sprintf("Hello %s, I'm alive and OK", name)
	return
}

// GetInfoUtilizador Busca toda a informação do utilizador
func GetInfoUtilizador(usrNome string) (result map[string]interface{}) {
	result = make(map[string]interface{})

	filter := bson.M{"nome": usrNome}
	colecao := MongoClient.Database("users_data").Collection("account_info")

	var registoUser resolvedschema.Utilizador
	context, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	err := colecao.FindOne(context, filter, options.FindOne()).Decode(&registoUser)
	defer cancel()

	if err != nil {
		loggers.OperacoesBDLogger.Println("Erro ao procurar pelo utilizador: ", usrNome)
		result["erro"] = "Erro ao procurar pelo utilizador: " + usrNome
		return
	}

	result["user"] = registoUser
	return
}

// UpdateInfoUrilizador Atualiza todos os dados especificádos, nos parametros da func, de um utilizador.
func UpdateInfoUrilizador(usrNome string, params ...map[string]interface{}) (result map[string]interface{}) {
	result = make(map[string]interface{})

	filter := bson.M{"nome": usrNome}
	colecao := MongoClient.Database("users_data").Collection("account_info")

	

	return
}
