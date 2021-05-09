package ficheiros

import (
	"context"
	"reflect"
	"time"

	"github.com/tomascpmarques/PAP/backend/robinservicodocumentacao/endpointfuncs"
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

	if !VerificarRepoExiste(ficheiroMetaData["reponome"].(string)) {
		loggers.OperacoesBDLogger.Println("O repo fornecido não existe, não se pode criar o ficheiro")
		retorno["erro"] = "O repo fornecido não existe, não se pode criar o ficheiro"
		return
	}

	metaHash, err := CriarMetaHash(ficheiroMetaData)
	if err != nil {
		loggers.OperacoesBDLogger.Println("Erro ao criar hash para meta data: ", err)
		retorno["erro"] = "Erro ao criar hash para meta data fornecida"
		return
	}
	ficheiroMetaData["hash"] = metaHash
	if err := MetaDataBaseValida(ficheiroMetaData); err != nil {
		loggers.OperacoesBDLogger.Println(err.Error())
		retorno["erro"] = err.Error()
		return
	}
	ficheiroMetaData["criacao"] = time.Now().Unix()
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

	// Insere o nome e o path do novo ficheiro, no repo onde a meta data do fiche. especifica
	err = RepoInserirMetaFileInfo(ficheiro.RepoNome, &ficheiro)
	if err != nil {
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

	if endpointfuncs.VerificarTokenUser(token) != "OK" {
		loggers.OperacoesBDLogger.Println("Erro: A token fornecida é inválida ou expirou")
		retorno["erro"] = "A token fornecida é inválida ou expirou"
		return
	}

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
	if endpointfuncs.VerificarTokenUserSpecif(token, campos["autor"].(string)) != "OK" || endpointfuncs.VerificarTokenAdmin(token) != "OK" {
		loggers.ServerErrorLogger.Println("Erro: Este utilizador não têm permissões para esta operação")
		retorno["erro"] = "Este utilizador não têm permissões para esta operação"
		return
	}

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

	retorno["sucesso"] = true
	return
}

// AtualizarFicheiroMetaData Busca um ficheiro pela sua hash e atualiza a meta-data através das atuali. fornecidas
// TODO
