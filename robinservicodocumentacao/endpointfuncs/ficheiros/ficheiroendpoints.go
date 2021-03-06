package ficheiros

import (
	"context"
	"reflect"
	"time"

	"github.com/tomascpmarques/PAP/backend/robinservicodocumentacao/endpointfuncs"
	"github.com/tomascpmarques/PAP/backend/robinservicodocumentacao/endpointfuncs/reposfiles"
	"github.com/tomascpmarques/PAP/backend/robinservicodocumentacao/loggers"
	"github.com/tomascpmarques/PAP/backend/robinservicodocumentacao/resolvedschema"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//  CriarFicheiroMetaData Cria a meta data de um ficheiro, para prepara o upload de conteúdo
func CriarFicheiroMetaData(ficheiroMetaData map[string]interface{}, token string) (retorno map[string]interface{}) {
	retorno = make(map[string]interface{})

	// if endpointfuncs.VerificarTokenUser(token) != "OK" {
	// 	loggers.OperacoesBDLogger.Println("Erro: A token fornecida é inválida ou expirou")
	// 	retorno["erro"] = "A token fornecida é inválida ou expirou"
	// 	return
	// }

	// Verificar se o repo a inserir a meta-info existe
	if !VerificarRepoExiste(ficheiroMetaData["reponome"].(string)) {
		loggers.OperacoesBDLogger.Println("O repo fornecido não existe, não se pode criar o ficheiro")
		retorno["erro"] = "O repo fornecido não existe, não se pode criar o ficheiro"
		return
	}

	// Criação da hash para a meta info do ficheiro
	metaHash, err := CriarMetaHash(ficheiroMetaData)
	if err != nil {
		loggers.ServerErrorLogger.Println("Erro ao criar hash para meta data: ", err)
		retorno["erro"] = "Erro ao criar hash para meta data fornecida"
		return
	}
	ficheiroMetaData["hash"] = metaHash
	// Verificar a validade da meta-info fornecida
	if err := MetaDataBaseValida(ficheiroMetaData); err != nil {
		loggers.ServerErrorLogger.Println(err.Error())
		retorno["erro"] = err.Error()
		return
	}

	// Atribuição da data de criação
	ficheiroMetaData["criacao"] = time.Now().Local().Format("2006/01/02 15:04:05")
	ficheiro := resolvedschema.FicheiroMetaDataParaStruct(&ficheiroMetaData)

	// Get the mongo colection
	mongoCollection := endpointfuncs.MongoClient.Database("documentacao").Collection("files-meta-data")
	cntx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	// Inser a meta data do file na bd respetiva para esses dados i.e: files-meta-data
	insserido, err := mongoCollection.InsertOne(cntx, ficheiro, options.InsertOne())
	defer cancel()
	if err != nil || !reflect.ValueOf(insserido.InsertedID).IsValid() {
		loggers.ServerErrorLogger.Println("Erro ao insserir o registo na BD: ", err)
		retorno["erro"] = "Erro ao insserir o registo na BD"
		return
	}

	// Insere o nome e o path do novo ficheiro, no repo onde a meta data do fiche. especificado
	err = RepoInserirMetaFileInfo(ficheiro.RepoNome, &ficheiro)
	if err != nil {
		loggers.ServerErrorLogger.Println("Erro: ", err)
		retorno["erro"] = err
		return
	}

	// Adiciona o ficheiro ás contribuições do user no serviço user-info
	if err := ModificarContrbFileInRepoUsrInfo("add", ficheiro.Autor, ficheiroMetaData["reponome"].(string), ficheiro.Nome, token); err != nil {
		loggers.ServerErrorLogger.Println("Erro: ", err)
		retorno["erro"] = err
		return
	}

	// Cria o ficheiro em local-storage após a criação da meta-data correspondente
	if err := reposfiles.CriarFicheiro_repo(&ficheiro); err != nil {
		loggers.ServerErrorLogger.Println("Erro: ", err)
		retorno["erro"] = err
		return
	}

	loggers.OperacoesBDLogger.Println("Meta Data insserida com sucesso")
	retorno["sucesso"] = true
	return
}

func BuscarMetaData(campos map[string]interface{}, token string) (retorno map[string]interface{}) {
	retorno = make(map[string]interface{})

	// if endpointfuncs.VerificarTokenUser(token) != "OK" {
	// 	loggers.OperacoesBDLogger.Println("Erro: A token fornecida é inválida ou expirou")
	// 	retorno["erro"] = "A token fornecida é inválida ou expirou"
	// 	return
	// }

	// Busca a meta data que corresponde aos campos dados
	// De um só registo
	metaData := GetMetaDataFicheiro(campos)
	if reflect.ValueOf(metaData).IsZero() {
		loggers.ServerErrorLogger.Println("Erro: Sem meta data para esse ficheiro")
		retorno["erro"] = "Sem meta data para esse ficheiro"
		return
	}

	loggers.OperacoesBDLogger.Println("Meta Data encontrada com sucesso")
	retorno["meta_data"] = metaData
	return
}

// ApagarFicheiroMetaData Apaga a meta data referente a um ficheiro
func ApagarFicheiroMetaData(campos map[string]interface{}, token string) (retorno map[string]interface{}) {
	retorno = make(map[string]interface{})

	// Verificação de igualdade entre request user, e file autor
	// if endpointfuncs.VerificarTokenUserSpecif(token, campos["autor"].(string)) != "OK" || endpointfuncs.VerificarTokenAdmin(token) != "OK" {
	// 	loggers.ServerErrorLogger.Println("Erro: Este utilizador não têm permissões para esta operação, ou token expirada")
	// 	retorno["erro"] = "Este utilizador não têm permissões para esta operação, ou token expirada"
	// 	return
	// }

	// Cria a hash dos campos fornecidos para procurar a meta data respetiva
	metaHash, err := CriarMetaHash(campos)
	if err != nil {
		loggers.ServerErrorLogger.Println("Erro ao criar hash para meta data: ", err)
		retorno["erro"] = "Erro ao criar hash para meta data fornecida"
		return
	}

	// Busca a meta data que corresponde aos campos dados, de um só registo
	ficheiroMetaData := GetMetaDataFicheiro(campos)
	if reflect.ValueOf(ficheiroMetaData).IsZero() {
		loggers.ServerErrorLogger.Println("Erro: Sem meta data para esse ficheiro")
		retorno["erro"] = "Sem meta data para esse ficheiro"
		return
	}

	// Apaga o ficheiro de local storage
	err = reposfiles.ApagarFicheiro_repo(&ficheiroMetaData)
	if err != nil {
		loggers.ServerErrorLogger.Println("Erro: Sem meta data para esse ficheiro para podêlo apagar do repo")
		retorno["erro"] = "Sem meta data para esse ficheiro para podêlo apagar do repo"
		return
	}

	// Apaga o ficheiro que contêm o campo "hash" igual ao fornecido
	// Na bd da meta data
	err = ApagarMetaDataFicheiro(metaHash)
	if err != nil {
		loggers.ServerErrorLogger.Println("Erro: Não foi possivél apagar este ficheiro: ", err)
		retorno["erro"] = "Não foi possivél apagar este ficheiro"
		return
	}

	// Apaga o ficheiro que contêm o campo "hash" igual ao fornecido, no repositório indicado no mongoDB
	err = ApagarFicheiroMetaRepo(metaHash, campos["autor"].(string))
	if err != nil {
		loggers.ServerErrorLogger.Println("Não foi possivél apagar um ficheiro devido ao erro: ", err)
		retorno["erro"] = "Não foi possivél apagar este ficheiro"
		return
	}

	// Remove o ficheiro das contribuições do user no sistema user-info
	err = ModificarContrbFileInRepoUsrInfo("rmv", campos["autor"].(string), campos["reponome"].(string), campos["nome"].(string), token)
	if err != nil {
		loggers.ServerErrorLogger.Println("Erro: ", err)
		retorno["erro"] = err
		return
	}

	retorno["sucesso"] = true
	return
}

func InserirConteudoFicheiro(contntMeta map[string]interface{}, token string) (retorno map[string]interface{}) {
	retorno = make(map[string]interface{})

	// Verificação de igualdade entre request user, e file autor
	if endpointfuncs.VerificarTokenUser(token) != "OK" {
		loggers.ServerErrorLogger.Println("Erro: Este utilizador não têm permissões para esta operação, ou token expirada")
		retorno["erro"] = "Este utilizador não têm permissões para esta operação, ou token expirada"
		return
	}

	ficheiroStruct := resolvedschema.FicheiroConteudoParaStruct(&contntMeta)
	if ficheiroStruct.Nome != ficheiroStruct.Path[len(ficheiroStruct.Path)-1] {
		loggers.ServerErrorLogger.Println("Erro: o nome do ficheiro não corresponde ao nome no path")
		retorno["erro"] = "O nome do ficheiro não corresponde ao nome no path"
		return
	}

	// Verificação da check sum do ficheiro
	err := ConteudoRecebidoCheckSum(&ficheiroStruct)
	if err != nil {
		loggers.ServerErrorLogger.Println("Erro: ", err.Error())
		retorno["erro"] = err.Error()
		return
	}

	// Get user from token, para evitar registo que includam o user
	tokenClaims := endpointfuncs.DevolveTokenClaims(token)
	usr := reflect.ValueOf(tokenClaims["user"]).String()

	if err := AdicionarContribuicaoRepo(&ficheiroStruct, usr); err != nil {
		loggers.ServerErrorLogger.Println("Erro: ", err.Error())
		retorno["erro"] = err.Error()
		return
	}

	// Inserção do conteudo de ficheiro recebido, no ficheiro pré-criado correspondente
	if err := reposfiles.AdicionarConteudoFicheiro_file(&ficheiroStruct); err != nil {
		loggers.ServerErrorLogger.Println("Erro: ", err.Error())
		retorno["erro"] = err.Error()
		return
	}

	loggers.OperacoesBDLogger.Println("Conteudo adicionado com sucesso")
	retorno["sucesso"] = true
	return
}

// BuscarConteudoFicheiro Busca um ficheiro lê o seu conteudo e devolve oa user
func BuscarConteudoFicheiro(campos map[string]interface{}, token string) (retorno map[string]interface{}) {
	retorno = make(map[string]interface{})
	var err error

	// if endpointfuncs.VerificarTokenUser(token) != "OK" {
	// 	loggers.OperacoesBDLogger.Println("Erro: A token fornecida é inválida ou expirou")
	// 	retorno["erro"] = "A token fornecida é inválida ou expirou"
	// 	return
	// }

	// Converte o query para um ficheiro meta info
	ficheiroContMeta := resolvedschema.FicheiroMetaDataParaStruct(&campos)
	if !reflect.ValueOf(ficheiroContMeta).IsValid() {
		loggers.ServerErrorLogger.Println("Erro: Não foi possível converter o query para o ficheiro, numa struct")
		retorno["erro"] = "Não fomos capazaes de concluir o request"
		return
	}

	// Coloca o conteudo, hash, etc, na response
	retorno["conteudo"], err = reposfiles.GetConteudoFicheiro_file(&ficheiroContMeta)
	if err != nil {
		loggers.ServerErrorLogger.Println("Erro: Não foi possível ler o ficheiro: ", err)
		retorno["erro"] = "Não fomos capazaes de concluir o request " + err.Error()
		return
	}

	loggers.OperacoesBDLogger.Println("Conteudo do ficheiro lido e enviado: ", ficheiroContMeta.Nome)
	retorno["sucesso"] = true
	return
}

func VerificarFicheiroExiste(params map[string]interface{}, token string) (retorno map[string]interface{}) {
	retorno = make(map[string]interface{})

	// if endpointfuncs.VerificarTokenUser(token) != "OK" {
	// 	loggers.OperacoesBDLogger.Println("Erro: A token fornecida é inválida ou expirou")
	// 	retorno["erro"] = "A token fornecida é inválida ou expirou"
	// 	return
	// }

	file := resolvedschema.FicheiroMetaDataParaStruct(&params)
	if !(reflect.ValueOf(file).IsValid()) {
		loggers.DocsStorage.Println("Erro: Erro ao converter o path fornecido para meta info")
		retorno["erro"] = "Erro ao converter o path fornecido para meta info"
		return
	}

	existe, err := reposfiles.VerificarFileExiste(&file)
	if err != nil {
		loggers.DocsStorage.Println("Erro: ", err)
		retorno["erro"] = err.Error()
		return
	}

	loggers.DocsStorage.Println("Procura com sucesso")
	retorno["existe"] = existe
	return
}

// AtualizarFicheiroMetaData Busca um ficheiro pela sua hash e atualiza a meta-data através das atualizações fornecidas
// TODO Hennnnnn mais ou menos (fazes se tiveres tempo, :))
