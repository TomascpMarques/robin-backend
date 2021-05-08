package endpointfuncs

import (
	"fmt"
	"reflect"

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

// CriarRepositorio Cria um repo para guardar a informação relativa a um tema e/ou tarefa
func CriarRepositorio(repoInfo map[string]interface{}, token string) (retorno map[string]interface{}) {
	retorno = make(map[string]interface{})

	// if VerificarTokenUser(token) != "OK" {
	// 	fmt.Println("Erro: A token fornecida é inválida ou expirou")
	// 	retorno["erro"] = "A token fornecida é inválida ou expirou"
	// 	return retorno
	// }

	// Get the mongo colection
	mongoCollection := MongoClient.Database("documentacao").Collection("repos")

	// Insser um registo na coleção e base de dados especificada
	registo, err := mongodbhandle.InsserirUmRegisto(repoInfo, mongoCollection, 10)
	if err != nil {
		loggers.DbFuncsLogger.Println("Não foi possivél criar o repositório pedido: ", repoInfo["nome"])
		retorno["erro"] = ("Não foi possivél criar o repositório pedido: ." + repoInfo["nome"].(string))
		return
	}

	loggers.OperacoesBDLogger.Println("Repo criado com sucesso! <", repoInfo["nome"], ">")
	retorno["resultado"] = registo
	return
}

// BuscarRepositorio Busca um repositório existente, e devolve a sua estrutura/conteúdos
func BuscarRepositorio(repoCampo string, campoValor interface{}, token string) (retorno map[string]interface{}) {
	retorno = make(map[string]interface{})

	// if VerificarTokenUser(token) != "OK" {
	// 	fmt.Println("Erro: A token fornecida é inválida ou expirou")
	// 	retorno["erro"] = "A token fornecida é inválida ou expirou"
	// 	return retorno
	// }

	// Busca o repositório por um campo especifico, e o valor esperado nesse campo
	repositorio := GetRepoPorCampo(repoCampo, campoValor)

	// Se o resultado da busca for nil, devolve umas menssagens de erro
	if reflect.ValueOf(repositorio).IsZero() {
		loggers.OperacoesBDLogger.Println("Não foi possivél encontrar o repositório pedido pelo campo: ", repoCampo)
		retorno["erro"] = ("Não foi possivél encontrar o repositório pedido pelo campo: <" + repoCampo + ">")
		return
	}

	loggers.OperacoesBDLogger.Println("Rrepositório procurado pelo campo: ", repoCampo, ", enviado.")
	retorno["repo"] = repositorio
	return
}

//func DropRepositorio() (retorno map[string]interface{}) {}
