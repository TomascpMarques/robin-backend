package endpointfuncs

import (
	"context"
	"fmt"
	"reflect"
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
func PingServico(name string) (retorno map[string]interface{}) {
	retorno = make(map[string]interface{})
	retorno["status"] = fmt.Sprintf("Hello %s, I'm alive and OK", name)
	return
}

// GetInfoUtilizador Busca toda a informação do utilizador especificado pelo id usrNome
func GetInfoUtilizador(usrNome string) (retorno map[string]interface{}) {
	retorno = make(map[string]interface{})

	// Defenição do filter a usar nas pesquisas da bd
	filter := bson.M{"user": usrNome}
	// Conexão à bd e coleção a usar
	colecao := MongoClient.Database("users_data").Collection("account_info")

	// defenições de utilitários
	var registoUser resolvedschema.Utilizador
	context, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	// Procura por 1 registo que iguale às opções
	err := colecao.FindOne(context, filter, options.FindOne()).Decode(&registoUser)
	defer cancel()
	if err != nil {
		loggers.OperacoesBDLogger.Println("Erro ao procurar pelo utilizador: ", usrNome)
		retorno["erro"] = "Erro ao procurar pelo utilizador: " + usrNome
		return
	}

	retorno["user"] = registoUser
	return
}

// UpdateInfoUtilizador Atualiza todos os dados especificádos, nos parametros da func, de um utilizador.
func UpdateInfoUtilizador(usrNome string, params map[string]interface{}, token string) (retorno map[string]interface{}) {
	retorno = make(map[string]interface{})

	// Se a token não for de um admin, ou de o user em sí, não se regista um user novo
	if VerificarTokenAdmin(token) != "OK" /*|| VerificarTokenUserSpecif(token, usrNome) != "OK"*/ {
		loggers.ServerErrorLogger.Println("A token não têm permissões")
		retorno["error"] = "A token não têm permissões"
		return
	}

	// Defenição do filter a usar nas pesquisas da bd
	filter := bson.M{"user": usrNome}
	colecao := MongoClient.Database("users_data").Collection("account_info")
	context, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	// atualização do registo e retorno da operação
	registosUpdt, err := colecao.UpdateOne(context, filter, params, options.MergeUpdateOptions())
	defer cancel()
	if err != nil {
		loggers.OperacoesBDLogger.Println("Erro ao atualizar a info do utilizador, erro: ", err)
		retorno["erro"] = "Erro ao atualizar a info do utilizador: " + usrNome
		return
	}
	// Se o numero de registos atualizados for diferente da quantidade pedida dá erro
	if registosUpdt.ModifiedCount < 1 {
		if registosUpdt.ModifiedCount == 0 {
			loggers.OperacoesBDLogger.Println("Sem atualizações, campos iguais")
			retorno["num_campos_updt"] = "Sem atualizações, campos iguais."
			return
		}
		loggers.OperacoesBDLogger.Println("Erro ao atualizar a info do utilizador, erro: ", err)
		retorno["erro"] = "Erro ao atualizar a info do utilizador pedido, verifica os parametros sff (sem repetições de valores)."
		return
	}

	retorno["num_campos_updt"] = registosUpdt.ModifiedCount
	return
}

// CriarRegistoUser cria um registo mongo db, com parametros nulos ou não, excepto o username (sempre !null)
func CriarRegistoUser(userInfo map[string]interface{}, token string) (retorno map[string]interface{}) {
	retorno = make(map[string]interface{})

	// Se a token não for de um admin, ou de o user em sí, não se regista um user novo
	if VerificarTokenAdmin(token) != "OK" {
		loggers.ServerErrorLogger.Println("A token não têm permissões")
		retorno["error"] = "A token não têm permissões"
		return
	}

	info := resolvedschema.UtilizadorParaStruct(&userInfo)
	if reflect.ValueOf(info).IsZero() {
		loggers.ServerErrorLogger.Println("Erro ao converter os dados para a struct correta")
		retorno["error"] = "Erro com o tipo de dados e sua conversão"
		return
	}

	colecao := MongoClient.Database("users_data").Collection("account_info")

	insserted, err := mongodbhandle.InsserirUmRegisto(info, colecao, 10)
	if err != nil {
		loggers.ServerErrorLogger.Println("Error: ", err)
		retorno["Error"] = err
		return nil
	}

	retorno["insserido"] = insserted
	return retorno
}
