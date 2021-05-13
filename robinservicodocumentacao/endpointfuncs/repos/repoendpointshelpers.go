package repos

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"

	"github.com/tomascpmarques/PAP/backend/robinservicodocumentacao/endpointfuncs"
	"github.com/tomascpmarques/PAP/backend/robinservicodocumentacao/loggers"
	"github.com/tomascpmarques/PAP/backend/robinservicodocumentacao/resolvedschema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetRepoPorCampo Busca um repo e devolveo na struct resolvedschema.Repositorio
// Busca o repositório através de um campo e valor do mesmo, especificado na sua chamada
func GetRepoPorCampo(campo string, valor interface{}) (repo resolvedschema.Repositorio) {
	// Documento e repo onde procurar o repo
	collection := endpointfuncs.MongoClient.Database("documentacao").Collection("repos")
	cntx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	// Procura na BD do registo pedido
	err := collection.FindOne(cntx, map[string]interface{}{string(campo): valor}, options.FindOne()).Decode(&repo)
	defer cancel()
	if err != nil {
		// Devolve um repo vzaio se não se encontrar o pedido
		repo = resolvedschema.Repositorio{}
		return
	}

	// Devolve repo
	return
}

func DropRepoPorNome(repoNome string) (erro error) {
	// Define o filtro a usar na procura de informação na BD
	filter := bson.M{"nome": repoNome}
	// Documento e repo onde procurar o repo
	collection := endpointfuncs.MongoClient.Database("documentacao").Collection("repos")
	cntx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	// Depois de apagar o registo, a var err,
	// Vai ter o sucesso ou falhanço da operação como o seu valor
	_, erro = collection.DeleteOne(cntx, filter)
	defer cancel()

	return
}

func RepoDropFicheirosMeta(repoNome string) error {
	collection := endpointfuncs.MongoClient.Database("documentacao").Collection("files-meta-data")
	cntx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	// Apaga toda a meta info dos ficheiros existentes, que sejam do repo especificado
	_, err := collection.DeleteMany(cntx, bson.M{"reponome": repoNome})
	defer cancel()
	if err != nil {
		return err
	}
	return nil
}

func UpdateRepositorioPorNome(repoName string, mundancas map[string]interface{}) *mongo.UpdateResult {
	// Set-up do filtro
	filter := bson.M{"nome": repoName}

	// Get collection
	coll := endpointfuncs.MongoClient.Database("documentacao").Collection("repos")
	cntx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	// Atualiza o item através do map especificado nos params
	matchCount, err := coll.UpdateOne(cntx, filter, mundancas, options.MergeUpdateOptions())
	defer cancel()
	if err != nil {
		loggers.DbFuncsLogger.Println("Erro ao atualizar o registo")
		return nil
	}

	return matchCount
}

func InitRepoFichrContribCriacao(repo *resolvedschema.Repositorio) {
	// A inicialização destas estruturas de dadds,
	// evita bugs com o display de informação,
	// na front-end e evita possivéis erros de mudança dos mesmos.
	repo.Contribuidores = make([]string, 0)
	repo.Ficheiros = make([]resolvedschema.RepositorioMetaFileInfo, 0)

	// Data e hora da criação do repo, no server-side
	repo.Criacao = time.Now().Local().Format("2006/01/02 15:04:05")
}

func VerificarInfoBaseRepo(info map[string]interface{}) (err error) {
	err = nil
	// Keys obrigatorias que o a info deve conter
	keysObrg := []string{
		"nome",
		"autor",
		"tema",
	}
	// Itera sobre as keys
	for _, v := range keysObrg {
		if valor, existe := info[v]; !(reflect.ValueOf(valor).IsValid()) || !existe {
			err = errors.New("os dados fornecidos não cumpre os parametros minímos")
			break
		}
	}
	return
}

// Adiciona o repo no serviço user-info, após criação neste serviço
func AdicionarContrbRepoUsrInfo(repo *resolvedschema.Repositorio, token string) error {
	// Mongodb query para atualizar o status do user
	adicionarQuery := fmt.Sprintf("\"%s\",\n\"%s\",\n\"%s\",", repo.Autor, repo.Nome, token)
	// DynamicGoQuery body para conssumir o endpoint do serviço userinfo
	action := fmt.Sprintf("action:\nfuncs:\n\"AdicionarContrbRepo\":\n%s", adicionarQuery)

	// Utilização do endpoint UpdateInfoUtilizador, exposto em http://0.0.0.0:8001
	resp, err := http.Post("http://0.0.0.0:8001", "text/plain", bytes.NewBufferString(action))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bodyContentBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	loggers.ResolverLogger.Printf("AdicionarContrbRepo status: %v", string(bodyContentBytes))
	return nil
}

func RemoverContrbRepoFileUsrInfo(repo *resolvedschema.Repositorio, token string) error {
	// Mongodb query para atualizar o status do user
	rmvQueryoptions := fmt.Sprintf(`{"user": %s,"repo": %s, "file": %s}`, repo.Autor, repo.Nome, token)
	adicionarQuery := fmt.Sprintf("\"%s\",%s,\"%s\",\n", "rmv", rmvQueryoptions, token)
	// DynamicGoQuery body para conssumir o endpoint do serviço userinfo
	action := fmt.Sprintf("action:\nfuncs:\n\"ModificarContribuicoes\":\n%s", adicionarQuery)

	// Utilização do endpoint UpdateInfoUtilizador, exposto em http://0.0.0.0:8001
	resp, err := http.Post("http://0.0.0.0:8001", "text/plain", bytes.NewBufferString(action))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bodyContentBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	loggers.ResolverLogger.Printf("AdicionarContrbRepo status: %v", string(bodyContentBytes))

	return nil
}
