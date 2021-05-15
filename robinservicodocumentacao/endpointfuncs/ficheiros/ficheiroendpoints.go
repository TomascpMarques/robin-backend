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
		loggers.OperacoesBDLogger.Println("Erro ao criar hash para meta data: ", err)
		retorno["erro"] = "Erro ao criar hash para meta data fornecida"
		return
	}
	ficheiroMetaData["hash"] = metaHash
	// Verificar a validade da meta-info fornecida
	if err := MetaDataBaseValida(ficheiroMetaData); err != nil {
		loggers.OperacoesBDLogger.Println(err.Error())
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
		loggers.OperacoesBDLogger.Println("Erro ao insserir o registo na BD: ", err)
		retorno["erro"] = "Erro ao insserir o registo na BD"
		return
	}

	// Insere o nome e o path do novo ficheiro, no repo onde a meta data do fiche. especificado
	err = RepoInserirMetaFileInfo(ficheiro.RepoNome, &ficheiro)
	if err != nil {
		loggers.OperacoesBDLogger.Println("Erro: ", err)
		retorno["erro"] = err
		return
	}

	// Adiciona o ficheiro ás contribuições do user no serviço user-info
	if err := ModificarContrbFileInRepoUsrInfo("add", ficheiro.Autor, ficheiroMetaData["reponome"].(string), ficheiro.Nome, token); err != nil {
		loggers.OperacoesBDLogger.Println("Erro: ", err)
		retorno["erro"] = err
		return
	}

	// Cria o ficheiro em local-storage após a criação da meta-data correspondente
	if err := reposfiles.CriarFicheiro_repo(&ficheiro); err != nil {
		loggers.OperacoesBDLogger.Println("Erro: ", err)
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
		loggers.OperacoesBDLogger.Println("Erro: Sem meta data para esse ficheiro")
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
		loggers.OperacoesBDLogger.Println("Erro ao criar hash para meta data: ", err)
		retorno["erro"] = "Erro ao criar hash para meta data fornecida"
		return
	}

	// Apaga o ficheiro que contêm o campo "hash" igual ao fornecido
	err = ApagarMetaDataFicheiro(metaHash)
	if err != nil {
		loggers.OperacoesBDLogger.Println("Erro: Não foi possivél apagar este ficheiro: ", err)
		retorno["erro"] = "Não foi possivél apagar este ficheiro"
		return
	}

	// Apaga o ficheiro que contêm o campo "hash" igual ao fornecido, no repositório indicado
	err = ApagarFicheiroMetaRepo(metaHash, campos["autor"].(string))
	if err != nil {
		loggers.OperacoesBDLogger.Println("Não foi possivél apagar um ficheiro devido ao erro: ", err)
		retorno["erro"] = "Não foi possivél apagar este ficheiro"
		return
	}

	// Remove o ficheiro das contribuições do user no sistema user-info
	err = ModificarContrbFileInRepoUsrInfo("rmv", campos["autor"].(string), campos["reponome"].(string), campos["nome"].(string), token)
	if err != nil {
		loggers.OperacoesBDLogger.Println("Erro: ", err)
		retorno["erro"] = err
		return
	}

	// Busca a meta data que corresponde aos campos dados
	// De um só registo
	ficheiroMetaData := GetMetaDataFicheiro(campos)
	if reflect.ValueOf(ficheiroMetaData).IsZero() {
		loggers.OperacoesBDLogger.Println("Erro: Sem meta data para esse ficheiro")
		retorno["erro"] = "Sem meta data para esse ficheiro"
		return
	}

	// Apaga o ficheiro de local storage
	err = reposfiles.ApagarFicheiro_repo(&ficheiroMetaData)
	if err != nil {
		loggers.OperacoesBDLogger.Println("Erro: Sem meta data para esse ficheiro")
		retorno["erro"] = "Sem meta data para esse ficheiro"
		return
	}

	retorno["sucesso"] = true
	return
}

func InserirConteudoFicheiro(contntMeta map[string]interface{}, token string) (retorno map[string]interface{}) {
	retorno = make(map[string]interface{})

	// Verificação de igualdade entre request user, e file autor
	// if endpointfuncs.VerificarTokenUserSpecif(token, campos["autor"].(string)) != "OK" || endpointfuncs.VerificarTokenAdmin(token) != "OK" {
	// 	loggers.ServerErrorLogger.Println("Erro: Este utilizador não têm permissões para esta operação, ou token expirada")
	// 	retorno["erro"] = "Este utilizador não têm permissões para esta operação, ou token expirada"
	// 	return
	// }

	ficheiroStruct := resolvedschema.FicheiroConteudoParaStruct(&contntMeta)
	if ficheiroStruct.Nome != ficheiroStruct.Path[len(ficheiroStruct.Path)-1] {
		loggers.OperacoesBDLogger.Println("Erro: o nome do ficheiro não corresponde ao nome no path")
		retorno["erro"] = "O nome do ficheiro não corresponde ao nome no path"
		return
	}

	// Verificação da check sum do ficheiro
	err := ConteudoRecebidoCheckSum(&ficheiroStruct, ficheiroStruct.Hash)
	if err != nil {
		loggers.OperacoesBDLogger.Println("Erro: ", err.Error())
		retorno["erro"] = err.Error()
		return
	}

	// Inserção do conteudo de ficheiro recebido, no ficheiro pré-criado correspondente
	if err := reposfiles.AdicionarConteudoFicheiro_file(&ficheiroStruct); err != nil {
		loggers.OperacoesBDLogger.Println("Erro: ", err.Error())
		retorno["erro"] = err.Error()
		return
	}

	loggers.OperacoesBDLogger.Println("Conteudo adicionado com sucesso")
	retorno["sucesso"] = true
	return
}

// AtualizarFicheiroMetaData Busca um ficheiro pela sua hash e atualiza a meta-data através das atualizações fornecidas
// TODO Hennnnnn mais ou menos (fazes se tiveres tempo, :))
