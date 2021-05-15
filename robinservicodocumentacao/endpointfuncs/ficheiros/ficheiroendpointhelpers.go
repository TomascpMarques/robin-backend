package ficheiros

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/tomascpmarques/PAP/backend/robinservicodocumentacao/endpointfuncs"
	"github.com/tomascpmarques/PAP/backend/robinservicodocumentacao/endpointfuncs/repos"
	"github.com/tomascpmarques/PAP/backend/robinservicodocumentacao/loggers"
	"github.com/tomascpmarques/PAP/backend/robinservicodocumentacao/resolvedschema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// VerifCamposBaseMeta Valida os campos existentes com os campos que são obrigatórios
// A funcção valida os existentes, se eles conterem os obrigatórios
func VerifCamposBaseMeta(metaData map[string]interface{}, camposObrg []string) error {
	// Verifica se contêm os campos básicos
	for _, campo := range camposObrg {
		if valor, existe := metaData[campo]; !(reflect.ValueOf(valor).IsValid()) || !existe {
			return errors.New("a meta data fornecida não é sufeciente para a criação do ficheiro")
		}
	}

	return nil
}

// VerifMetaNomeELower Verifica se o nome do ficheir está todo em lowercase
func VerifMetaNomeELower(meta *resolvedschema.FicheiroMetaData) error {
	if meta.Nome != strings.ToLower(meta.Nome) {
		return errors.New("não foi possivél criar o ficheiro pedido, o nome do ficheiro deve ser todo em lower case")
	}
	return nil
}

func VerifPathMinLen(meta *resolvedschema.FicheiroMetaData) error {
	// O path têm de conter mais de 2 elementos
	if len(meta.Path) < 2 {
		return errors.New("não foi possivél criar o ficheiro pedido, erros no seu path")
	}
	return nil
}

func VerifPathValido(meta *resolvedschema.FicheiroMetaData) error {
	// Se o path não corresponder ao seguinte formato: "repo/<file_repo_name>/.../<file_name>"
	if meta.Path[1] != meta.RepoNome || meta.Path[0] != "repo" || meta.Path[len(meta.Path)-1] != meta.Nome {
		return errors.New("não foi possivél criar o ficheiro pedido, erros no armazenamento descrito na meta data")
	}
	return nil
}

func MetaDataBaseValida(metaData map[string]interface{}) error {
	campos := []string{
		"nome",
		"autor",
		"reponome",
		"path",
	}
	if err := VerifCamposBaseMeta(metaData, campos); err != nil {
		return err
	}

	// Conversão de map[string]interface{}, para a struct correta
	meta := resolvedschema.FicheiroMetaDataParaStruct(&metaData)

	if err := VerifMetaNomeELower(&meta); err != nil {
		return err
	}
	if err := VerifPathMinLen(&meta); err != nil {
		return err
	}
	if err := VerifPathValido(&meta); err != nil {
		return err
	}

	// Define o filtro a usar na procura de informação na BD
	filter := bson.M{"hash": meta.Hash}
	// Documento e repo onde procurar o repo
	collection := endpointfuncs.MongoClient.Database("documentacao").Collection("files-meta-data")
	cntx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	// Procura por um registo com a mesma hash (registo igual)
	err := collection.FindOne(cntx, filter, options.FindOne()).Err()
	defer cancel()
	if err == nil {
		return errors.New("não foi possivél criar o ficheiro pedido, esse ficheiro já existe")
	}

	// Procura por um ficheiro com o mesmo path, o path é praticamente o identificador absoluto do ficheiro
	if path := GetMetaDataFicheiro(map[string]interface{}{"path": meta.Path}).Path; !reflect.ValueOf(path).IsZero() {
		return errors.New("não foi possivél criar o ficheiro pedido, esse path já existe")
	}

	return nil
}

// GetMetaDataPorCampo Busca meta data de um ficheiro e devolve o mesmo na struct resolvedschema.FicheiroMetaData
// Busca a meta data através de um campo e valor do mesmo, especificado na sua chamada
func GetMetaDataFicheiro(campos map[string]interface{}) (meta resolvedschema.FicheiroMetaData) {
	// Documento e Coleção onde procurar a meta data
	collection := endpointfuncs.MongoClient.Database("documentacao").Collection("files-meta-data")
	cntx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	// Procura na BD do registo pedido
	err := collection.FindOne(cntx, campos, options.FindOne()).Decode(&meta)
	defer cancel()
	if err != nil {
		// Devolve um repo vzaio se não se encontrar o pedido
		meta = resolvedschema.FicheiroMetaData{}
		return
	}

	// Devolve meta data
	return
}

// ApagarMetaDataFicheiro Apaga o ficheiro em que a hash é a mesma que a passada nos parametros
func ApagarMetaDataFicheiro(hash string) error {
	// Documento e Coleção onde procurar a meta data
	collection := endpointfuncs.MongoClient.Database("documentacao").Collection("files-meta-data")
	cntx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	// Procura na BD do registo pedido
	err := collection.FindOneAndDelete(cntx, bson.M{"hash": hash}, options.FindOneAndDelete())
	defer cancel()
	if err.Err() != nil {
		// Devolve um repo vzaio se não se encontrar o pedido
		return err.Err()
	}

	return nil
}

func ApagarFicheiroMetaRepo(hash string, user string) error {
	// Documento e Coleção onde procurar a meta data
	collection := endpointfuncs.MongoClient.Database("documentacao").Collection("repos")
	cntx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	// Query para apagar o ficheiro que conicide com a hash: hash
	operacaoDrop := bson.M{"$pull": bson.M{"ficheiros": bson.M{"hash": hash}}}
	// Operação para atualizar o repo, o filtro da pesquisa pelo registo, é o segundo param
	resultado, err := collection.UpdateOne(cntx, bson.M{"autor": user, "ficheiros.hash": hash}, operacaoDrop)
	defer cancel()

	// Error handeling
	if err != nil {
		return err
	}
	if resultado.ModifiedCount < 1 {
		return errors.New("nenhum ficheiro foi apagado")
	}

	return nil
}

// RepoInserirMetaFileInfo Atualiza o array de ficheiros que pertence ao repo especificado
func RepoInserirMetaFileInfo(repoNome string, meta *resolvedschema.FicheiroMetaData) error {
	if meta.Path[1] != repoNome {
		return errors.New("caminho do ficheiro não coincide com o do repositorio")
	}
	// Combinação de nome do ficheiro e do seu path
	ficheiroNomePath := map[string]interface{}{"nome": meta.Nome, "path": meta.Path, "hash": meta.Hash}

	colecao := endpointfuncs.MongoClient.Database("documentacao").Collection("repos")
	cntx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	err := colecao.FindOneAndUpdate(cntx, bson.M{"nome": repoNome}, bson.M{"$push": bson.M{"ficheiros": ficheiroNomePath}})
	defer cancel()
	if err.Err() != nil {
		// Devolve um repo vzaio se não se encontrar o pedido
		return err.Err()
	}

	// Verifica se o autor deste ficheiro é diferente do autor do repo,
	// Se sim, adiciona este utilizador À lista dos contribuidores
	repoAutor := repos.GetRepoPorCampo("nome", repoNome).Autor
	if err := VerificaNovoContribuidor(meta.Autor, repoAutor, repoNome); err != nil {
		return err
	}

	return nil
}

// VerificaNovoContribuidor Se o ficheiro a insserir no repo for de autoria de um user,
//							que não é o autor do repo, adiciona esse user aos contribuidores
func VerificaNovoContribuidor(ficheiroAutor string, repoAutor string, repoNome string) error {
	if ficheiroAutor != repoAutor {
		colecao := endpointfuncs.MongoClient.Database("documentacao").Collection("repos")
		cntx, cancel := context.WithTimeout(context.Background(), time.Second*10)

		err := colecao.FindOneAndUpdate(cntx, bson.M{"nome": repoNome}, bson.M{"$push": bson.M{"contribuidores": ficheiroAutor}})
		defer cancel()
		if err.Err() != nil {
			// Devolve um repo vzaio se não se encontrar o pedido
			return err.Err()
		}
	}

	return nil
}

// CriarMetaHash Cria uma hash da meta data do ficheiro
func CriarMetaHash(metaData map[string]interface{}) (string, error) {
	// encodifica a meta data para []byte (em formato JSON)
	x, err := json.Marshal(metaData)
	if err != nil {
		return "", err
	}
	// Adiciona a hash o valor convertido para []byte
	return fmt.Sprintf("%x", sha1.Sum(x)), nil
}

// VerificarRepoExiste Verifica se o repositório com este nome existe
func VerificarRepoExiste(repoNome string) bool {
	return !reflect.ValueOf(repos.GetRepoPorCampo("nome", repoNome)).IsZero()
}

func ModificarContrbFileInRepoUsrInfo(opDef string, usrNome string, repoAutor string, nomeFicheiro string, token string) error {
	fileAddSpecific := fmt.Sprintf(`{"user": "%s","repo": "%s", "file": "%s"}`, usrNome, repoAutor, nomeFicheiro)
	// Mongodb query para atualizar o status do user
	adicionarQuery := fmt.Sprintf("\"%s\",\n%s,\n\"%s\",", opDef, fileAddSpecific, token)
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

	loggers.ResolverLogger.Printf("ModificarContribuicoes status: %v\n", string(bodyContentBytes))
	return nil
}
