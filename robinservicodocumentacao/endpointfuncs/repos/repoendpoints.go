package repos

import (
	"reflect"

	"github.com/tomascpmarques/PAP/backend/robinservicodocumentacao/endpointfuncs"
	"github.com/tomascpmarques/PAP/backend/robinservicodocumentacao/loggers"
	"github.com/tomascpmarques/PAP/backend/robinservicodocumentacao/mongodbhandle"
)

// CriarRepositorio Cria um repo para guardar a informação relativa a um tema e/ou tarefa
func CriarRepositorio(repoInfo map[string]interface{}, token string) (retorno map[string]interface{}) {
	retorno = make(map[string]interface{})

	if endpointfuncs.VerificarTokenUser(token) != "OK" {
		loggers.OperacoesBDLogger.Println("Erro: A token fornecida é inválida ou expirou")
		retorno["erro"] = "A token fornecida é inválida ou expirou"
		return retorno
	}

	// Get the mongo colection
	mongoCollection := endpointfuncs.MongoClient.Database("documentacao").Collection("repos")

	if existe := GetRepoPorCampo("nome", repoInfo["nome"]); !(reflect.ValueOf(existe).IsZero()) {
		loggers.DbFuncsLogger.Println("Não foi possivél criar o repositório pedido: ", repoInfo["nome"], ".Já existe um com esse nome")
		retorno["erro"] = ("Não foi possivél criar o repositório pedido, devido ao nome ser igual a um existente")
		return
	}

	if err := VerificarInfoBaseRepo(repoInfo); err != nil {
		loggers.DbFuncsLogger.Println("Os dados estão incompletos para criar um repo")
		retorno["erro"] = "Os dados estão incompletos para criar um repo"
		return
	}

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

	// if endpointfuncs.VerificarTokenUser(token) != "OK" {
	// 	loggers.OperacoesBDLogger.Println("Erro: A token fornecida é inválida ou expirou")
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

// DropRepositorio Busca o repo especificado por nome e apaga o mesmo, se esse pedido for feito pelo autor do repo
func DropRepositorio(repoNome string, token string) (retorno map[string]interface{}) {
	retorno = make(map[string]interface{})

	if endpointfuncs.VerificarTokenUser(token) != "OK" {
		loggers.ServerErrorLogger.Println("Erro: A token fornecida é inválida ou expirou")
		retorno["erro"] = "A token fornecida é inválida ou expirou"
		return retorno
	}

	// Busca o repositório para se poder comparar o autor com o user que fez o pedido
	repositorio := GetRepoPorCampo("nome", repoNome)
	// Se o resultado da busca for nil, devolve umas menssagens de erro
	if reflect.ValueOf(repositorio).IsZero() {
		loggers.OperacoesBDLogger.Println("Não foi possivél encontrar o repositório pedido")
		retorno["erro"] = ("Não foi possivél encontrar o repositório pedido")
		return
	}

	// Verificação de igualdade entre request user, e repo autor
	if endpointfuncs.VerificarTokenUserSpecif(token, repositorio.Autor) != "OK" {
		loggers.ServerErrorLogger.Println("Erro: Este utilizador não têm permissões para esta operação")
		retorno["erro"] = "Este utilizador não têm permissões para esta operação"
		return retorno
	}

	// Drop do repo pedido
	if err := DropRepoPorNome(repoNome); err != nil {
		loggers.ServerErrorLogger.Println("Erro: Este utilizador não têm permissões para esta operação")
		retorno["erro"] = "Este utilizador não têm permissões para esta operação"
		return retorno
	}

	loggers.DbFuncsLogger.Println("Repositório apagado com sucesso")
	retorno["ok"] = true
	return
}

func UpdateRepositorio(repoNome string, updateQuery map[string]interface{}, token string) (retorno map[string]interface{}) {
	retorno = make(map[string]interface{})

	if endpointfuncs.VerificarTokenUser(token) != "OK" {
		loggers.ServerErrorLogger.Println("Erro: A token fornecida é inválida ou expirou")
		retorno["erro"] = "A token fornecida é inválida ou expirou"
		return
	}

	// Busca o repositório para se poder comparar o autor com o user que fez o pedido
	repositorio := GetRepoPorCampo("nome", repoNome)
	// Se o resultado da busca for nil, devolve umas menssagens de erro
	if reflect.ValueOf(repositorio).IsZero() {
		loggers.OperacoesBDLogger.Println("Não foi possivél encontrar o repositório pedido")
		retorno["erro"] = ("Não foi possivél encontrar o repositório pedido")
		return
	}

	// Verificação de igualdade entre request user, e repo autor
	if endpointfuncs.VerificarTokenUserSpecif(token, repositorio.Autor) != "OK" {
		loggers.ServerErrorLogger.Println("Erro: Este utilizador não têm permissões para esta operação")
		retorno["erro"] = "Este utilizador não têm permissões para esta operação"
		return
	}

	atualizacoes := UpdateRepositorioPorNome(repoNome, updateQuery) // i.e: {"$set":{"autor": "efefef"}},
	if atualizacoes == nil {
		loggers.ServerErrorLogger.Println("Erro ao atualizar os valores pedidos")
		retorno["erro"] = "Erro ao atualizar os valores pedidos"
		return
	}

	retorno["resultado"] = atualizacoes
	return
}
