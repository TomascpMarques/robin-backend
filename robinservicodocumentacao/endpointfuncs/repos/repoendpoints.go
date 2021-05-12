package repos

import (
	"reflect"

	"github.com/tomascpmarques/PAP/backend/robinservicodocumentacao/endpointfuncs"
	"github.com/tomascpmarques/PAP/backend/robinservicodocumentacao/loggers"
	"github.com/tomascpmarques/PAP/backend/robinservicodocumentacao/mongodbhandle"
	"github.com/tomascpmarques/PAP/backend/robinservicodocumentacao/resolvedschema"
	"go.mongodb.org/mongo-driver/bson"
)

// CriarRepositorio Cria um repo para guardar a informação relativa a um tema e/ou tarefa
func CriarRepositorio(repoInfo map[string]interface{}, token string) (retorno map[string]interface{}) {
	retorno = make(map[string]interface{})

	// if endpointfuncs.VerificarTokenUser(token) != "OK" {
	// 	loggers.OperacoesBDLogger.Println("Erro: A token fornecida é inválida ou expirou")
	// 	retorno["erro"] = "A token fornecida é inválida ou expirou"
	// 	return retorno
	// }

	// Get the mongo colection
	mongoCollection := endpointfuncs.MongoClient.Database("documentacao").Collection("repos")

	if repoExiste := GetRepoPorCampo("nome", repoInfo["nome"].(string)); !(reflect.ValueOf(repoExiste).IsZero()) {
		loggers.DbFuncsLogger.Println("Não foi possivél criar o repositório pedido: ", repoInfo["nome"], ".Já existe um com esse nome")
		retorno["erro"] = ("Não foi possivél criar o repositório pedido, devido ao nome ser igual a um existente")
		return
	}

	if err := VerificarInfoBaseRepo(repoInfo); err != nil {
		loggers.DbFuncsLogger.Println("Os dados estão incompletos para criar um repo")
		retorno["erro"] = "Os dados estão incompletos para criar um repo"
		return
	}
	// Transformação da informação de repo para uma struct do tipo Repo
	repo := resolvedschema.RepositorioParaStruct(&repoInfo)
	// Inicializa os arrays de contribuições e de ficheiros a zero, evita null erros no decoding
	InitRepoFichrContribCriacao(&repo)

	// Insser um registo na coleção e base de dados especificada
	registo, err := mongodbhandle.InsserirUmRegisto(repo, mongoCollection, 10)
	if err != nil {
		loggers.DbFuncsLogger.Println("Não foi possivél criar o repositório pedido: ", repoInfo["nome"])
		retorno["erro"] = ("Não foi possivél criar o repositório pedido: ." + repoInfo["nome"].(string))
		return
	}

	if AdicionarContrbRepoUsrInfo(&repo, token) != nil {
		loggers.DbFuncsLogger.Println("Não foi possivél inserir o repo na user-info")
		retorno["erro"] = ("Não foi possivél inserir o repo na user-info")
		return
	}

	loggers.OperacoesBDLogger.Println("Repo criado com sucesso! <", repoInfo["nome"], ">")
	retorno["resultado"] = registo
	return
}

// BuscarRepositorio Busca um repositório existente, e devolve a sua estrutura/conteúdos
func BuscarRepositorio(campos map[string]interface{}, token string) (retorno map[string]interface{}) {
	retorno = make(map[string]interface{})
	//fmt.Println("AND NOW THE TIME: ", time.Now().Local().Format("2006/01/02 15:04:05"))

	// if endpointfuncs.VerificarTokenUser(token) != "OK" {
	// 	loggers.OperacoesBDLogger.Println("Erro: A token fornecida é inválida ou expirou")
	// 	retorno["erro"] = "A token fornecida é inválida ou expirou"
	// 	return retorno
	// }

	// Busca o repositório por um campo especifico, e o valor esperado nesse campo
	repositorio := GetRepoPorCampo("nome", campos["nome"].(string))

	// Se o resultado da busca for nil, devolve umas menssagens de erro
	if reflect.ValueOf(repositorio).IsZero() {
		loggers.OperacoesBDLogger.Println("Não foi possivél encontrar o repositório pedido pelo campo: ", campos)
		retorno["erro"] = ("Não foi possivél encontrar o repositório pedido pelos campos pedidos")
		return
	}

	loggers.OperacoesBDLogger.Println("Rrepositório procurado pelo campo: ", campos, ", enviado.")
	retorno["repo"] = repositorio
	return
}

// DropRepositorio Busca o repo especificado pelos campos passados (o nome é obrigatorio), e apaga o mesmo, se esse pedido for feito pelo autor do repo
func DropRepositorio(campos map[string]interface{}, token string) (retorno map[string]interface{}) {
	retorno = make(map[string]interface{})

	// if endpointfuncs.VerificarTokenUser(token) != "OK" {
	// 	loggers.ServerErrorLogger.Println("Erro: A token fornecida é inválida ou expirou")
	// 	retorno["erro"] = "A token fornecida é inválida ou expirou"
	// 	return
	// }

	// Busca o repositório para se poder comparar o autor com o user que fez o pedido
	repositorio := GetRepoPorCampo("nome", campos["nome"].(string))
	// Se o resultado da busca for nil, devolve umas menssagens de erro
	if reflect.ValueOf(repositorio).IsZero() {
		loggers.OperacoesBDLogger.Println("Não foi possivél encontrar o repositório pedido")
		retorno["erro"] = ("Não foi possivél encontrar o repositório pedido")
		return
	}

	// Verificação de igualdade entre request user, e repo autor
	// if endpointfuncs.VerificarTokenUserSpecif(token, repositorio.Autor) != "OK" {
	// 	loggers.ServerErrorLogger.Println("Erro: Este utilizador não têm permissões para esta operação")
	// 	retorno["erro"] = "Este utilizador não têm permissões para esta operação"
	// 	return
	// }

	// Drop do repo pedido
	if err := DropRepoPorNome(campos["nome"].(string)); err != nil {
		loggers.ServerErrorLogger.Println("Erro: Este utilizador não têm permissões para esta operação")
		retorno["erro"] = "Este utilizador não têm permissões para esta operação"
		return
	}

	if err := RepoDropFicheirosMeta(campos["nome"].(string)); err != nil {
		loggers.ServerErrorLogger.Println("Erro: Ou o repo não tinha ficheiros ou ouve complicações para apagar esses ficheiros")
		retorno["erro"] = "Ou o repo não tinha ficheiros ou ouve complicações para apagar esses ficheiros"
		return
	}

	loggers.DbFuncsLogger.Println("Repositório apagado com sucesso")
	retorno["ok"] = true
	return
}

func UpdateRepositorio(campos map[string]interface{}, updateQuery map[string]interface{}, token string) (retorno map[string]interface{}) {
	retorno = make(map[string]interface{})

	// if endpointfuncs.VerificarTokenUser(token) != "OK" {
	// 	loggers.ServerErrorLogger.Println("Erro: A token fornecida é inválida ou expirou")
	// 	retorno["erro"] = "A token fornecida é inválida ou expirou"
	// 	return
	// }

	// Busca o repositório para se poder comparar o autor com o user que fez o pedido
	repositorio := GetRepoPorCampo("nome", campos["nome"].(string))
	// Se o resultado da busca for nil, devolve umas menssagens de erro
	if reflect.ValueOf(repositorio).IsZero() {
		loggers.OperacoesBDLogger.Println("Não foi possivél encontrar o repositório pedido")
		retorno["erro"] = ("Não foi possivél encontrar o repositório pedido")
		return
	}

	// Verificação de igualdade entre request user, e repo autor
	// if endpointfuncs.VerificarTokenUserSpecif(token, repositorio.Autor) != "OK" {
	// 	loggers.ServerErrorLogger.Println("Erro: Este utilizador não têm permissões para esta operação")
	// 	retorno["erro"] = "Este utilizador não têm permissões para esta operação"
	// 	return
	// }

	atualizacoes := UpdateRepositorioPorNome(campos["nome"].(string), bson.M{"$set": updateQuery}) // i.e: {"$set":{"autor": "efefef"}},
	if atualizacoes == nil {
		loggers.ServerErrorLogger.Println("Erro ao atualizar os valores pedidos")
		retorno["erro"] = "Erro ao atualizar os valores pedidos"
		return
	}

	retorno["resultado"] = atualizacoes
	return
}
