package endpointfuncs

import (
	"context"
	"fmt"
	"reflect"

	"github.com/tomascpmarques/PAP/backend/robinservicouserinfo/loggers"
	"github.com/tomascpmarques/PAP/backend/robinservicouserinfo/mongodbhandle"
	"github.com/tomascpmarques/PAP/backend/robinservicouserinfo/resolvedschema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
func GetInfoUtilizador(usrNome string, token string) (retorno map[string]interface{}) {
	retorno = make(map[string]interface{})

	// Se a token não for de um admin, ou do user em sí, não se atualiza o user
	if VerificarTokenUser(token) != "OK" /*|| VerificarTokenUserSpecif(token, usrNome) != "OK"*/ {
		loggers.ServerErrorLogger.Println("A token não têm permissões")
		retorno["error"] = "A token não têm permissões"
		return
	}

	// Defenição do filter a usar nas pesquisas da bd
	// filter := bson.M{"user": usrNome}
	// Conexão à bd e coleção a usar
	//operacoesColl := MongoClient.Database("users_data").Collection("account_info")

	operacoesColl := SetupColecao("users_data", "account_info")
	operacoesColl.Filter = bson.M{"user": usrNome}

	// defenições de utilitários
	var registoUser resolvedschema.Utilizador

	// Procura por 1 registo que iguale às opções
	err := operacoesColl.Colecao.FindOne(operacoesColl.Cntxt, operacoesColl.Filter, options.FindOne()).Decode(&registoUser)
	defer operacoesColl.CancelFunc()
	if err != nil {
		loggers.OperacoesBDLogger.Println("Erro ao procurar pelo utilizador: ", usrNome, err)
		retorno["erro"] = "Erro ao procurar pelo utilizador: " + usrNome
		return
	}

	retorno["user"] = registoUser
	return
}

// UpdateInfoUtilizador Atualiza todos os dados especificádos, nos parametros da func, de um utilizador.
func UpdateInfoUtilizador(usrNome string, params map[string]interface{}, token string) (retorno map[string]interface{}) {
	retorno = make(map[string]interface{})

	// Se a token não for de um admin, ou do user em sí, não se atualiza o user
	if VerificarTokenAdmin(token) != "OK" || VerificarTokenUserSpecif(token, usrNome) != "OK" {
		loggers.ServerErrorLogger.Println("A token não têm permissões")
		retorno["error"] = "A token não têm permissões"
		return
	}

	// Defenição do filter a usar nas pesquisas da bd
	filter := bson.M{"user": usrNome}
	operacoesColl := SetupColecao("users_data", "account_info")

	// atualização do registo e retorno da operação
	registosUpdt, err := operacoesColl.Colecao.UpdateOne(operacoesColl.Cntxt, filter, bson.M{"$set": params}, options.MergeUpdateOptions())
	defer operacoesColl.CancelFunc()
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

	// Se a token não for de um admin, não se regista um user novo
	if VerificarTokenAdmin(token) != "OK" {
		loggers.ServerErrorLogger.Println("A token não têm permissões")
		retorno["error"] = "A token não têm permissões"
		return
	}

	operacoesColl := SetupColecao("users_data", "account_info")

	// Verifica se a info do user que queremos inserir já existe
	exists := operacoesColl.Colecao.FindOne(context.Background(), bson.M{"user": userInfo["user"]})
	if exists != nil && exists.Err() != mongo.ErrNoDocuments {
		loggers.ServerErrorLogger.Println("Já existe informação para esse user")
		retorno["error"] = "Já existe informação para esse user"
		return

	}

	// Conver-te a info de user, para uma struct
	info := resolvedschema.UtilizadorParaStruct(&userInfo)
	if reflect.ValueOf(info).IsZero() {
		loggers.ServerErrorLogger.Println("Erro ao converter os dados para a struct correta")
		retorno["error"] = "Erro com o tipo de dados e sua conversão"
		return
	}

	// Insere o registo da info do user na BD
	inserted, err := mongodbhandle.InserirUmRegisto(info, operacoesColl.Colecao, 10)
	if err != nil {
		loggers.ServerErrorLogger.Println("Error: ", err)
		retorno["Error"] = err
		return
	}

	retorno["inserido"] = inserted
	return retorno
}

// ModificarContribuicoes Modifica o valor do array que contêm as contribuições, pode adicionar ou retirar desse mesmo array
func ModificarContribuicoes(operacaoConfig string, repoUpdate map[string]interface{}, token string) (retorno map[string]interface{}) {
	retorno = make(map[string]interface{})

	// Se a token não for de um utilizador, não executa a função
	// if VerificarTokenAdmin(token) != "OK" /*|| VerificarTokenUserSpecif(token, usrNome) != "OK"*/ {
	// 	loggers.ServerErrorLogger.Println("A token não têm permissões")
	// 	retorno["error"] = "A token não têm permissões"
	// 	return
	// }

	// Defenições a serem usadas para executar as operações na BD
	operacoesColl := SetupColecao("users_data", "account_info")
	operacoesColl.Filter = bson.M{"user": repoUpdate["user"], "contribuicoes.reponome": repoUpdate["repo"].(string)}

	// Avalia o tipo de operação pedido, add para adicionar contribuição, rmv para remover contribuição
	switch operacaoConfig {
	// Operação que adiciona um ficheiro ás contribuições do user
	case "add":
		err := operacoesColl.AdicionarContribuicao(repoUpdate["repo"].(string), repoUpdate["file"].(string))
		if err != nil {
			loggers.ServerErrorLogger.Println("Error: ", err)
			retorno["Error"] = err
			return
		}
		// Operação que remove um ficheiro ás contribuições do user
	case "rmv":
		err := operacoesColl.RemoverContribuicaoFile(repoUpdate["repo"].(string), repoUpdate["file"].(string))
		if err != nil {
			loggers.ServerErrorLogger.Println("Error: ", err)
			retorno["Error"] = err
			return
		}
		// Se a operação pedida não estiver implementada corre este código
	default:
		loggers.ServerErrorLogger.Println("Error: Tipo de operação não reconhecido, <'add' ou 'rmv'>")
		retorno["Error"] = "Tipo de operação não reconhecido, <'adicionar' ou 'remover'>"
		return
	}

	loggers.ResolverLogger.Println("Operação acabou com sucesso")
	retorno["sucesso"] = true
	return
}

// AdicionarContrbRepo Adiciona o repo de contribuições à informação do user
func AdicionarContrbRepo(usrNome string, repoNome string, token string) (retorno map[string]interface{}) {
	retorno = make(map[string]interface{})

	// Se a token não for de um utilizador, não executa a função
	// if VerificarTokenAdmin(token) != "OK" /*|| VerificarTokenUserSpecif(token, usrNome) != "OK"*/ {
	// 	loggers.ServerErrorLogger.Println("A token não têm permissões")
	// 	retorno["error"] = "A token não têm permissões"
	// 	return
	// }

	// Defenições a serem usadas para executar as operações na BD
	operacoesColl := SetupColecao("users_data", "account_info")
	operacoesColl.Filter = bson.M{"user": usrNome}

	// Avalia se a operação de criar o repo teve sucesso
	if err := operacoesColl.CriarRepoContribuicoes(repoNome); err != nil {
		loggers.ServerErrorLogger.Println("Erro ao criar repo nas contribuições do user")
		retorno["erro"] = "Erro ao criar repo nas contribuições do user"
		return
	}

	loggers.OperacoesBDLogger.Println("Repo criado nas contribuições do user")
	retorno["sucesso"] = true
	return
}

// RemoverRepoContributo Remove o repo de contribuições
func RemoverRepoContributo(repoinfo map[string]interface{}, token string) (retorno map[string]interface{}) {
	retorno = make(map[string]interface{})

	// Se a token não for de um utilizador, não executa a função
	// if VerificarTokenAdmin(token) != "OK" /*|| VerificarTokenUserSpecif(token, usrNome) != "OK"*/ {
	// 	loggers.ServerErrorLogger.Println("A token não têm permissões")
	// 	retorno["error"] = "A token não têm permissões"
	// 	return
	// }

	// Defenições a serem usadas para executar as operações na BD
	operacoesColl := SetupColecao("users_data", "account_info")
	operacoesColl.Filter = bson.M{"user": repoinfo["user"], "contribuicoes.reponome": repoinfo["repo"].(string)}

	// Remove o repo das contribuições do utilizador
	if err := operacoesColl.RemoverRepoContribuicao(repoinfo["repo"].(string)); err != nil {
		loggers.MongoDBLogger.Println(err)
		loggers.ServerErrorLogger.Println("Erro ao largar o repo nas contribuições do user pedido")
		retorno["erro"] = "Erro ao largar o repo nas contribuições do user pedido"
		return
	}

	loggers.ServerErrorLogger.Println("Repo largado das contribuições do utilizador com sucesso")
	retorno["sucesso"] = true
	return
}
