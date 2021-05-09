package ficheiros

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/tomascpmarques/PAP/backend/robinservicodocumentacao/endpointfuncs"
	"github.com/tomascpmarques/PAP/backend/robinservicodocumentacao/resolvedschema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MetaDataBaseValida(metaData map[string]interface{}) error {
	// Campos que a metaData têm que validar
	camposBase := []string{
		"nome",
		"autor",
		"repo_nome",
		"path",
	}
	// Verifica se contêm os campos básicos
	for _, campo := range camposBase {
		if valor, existe := metaData[campo]; !(reflect.ValueOf(valor).IsValid()) || !existe {
			return errors.New("a meta data fornecida não é sufeciente para a criação do ficheiro")
		}
	}

	// O path têm de conter mais de 2 elementos
	meta := resolvedschema.FicheiroMetaDataParaStruct(&metaData)
	if len(meta.Path) < 2 {
		return errors.New("não foi possivél criar o ficheiro pedido, erros no seu path")
	}

	// Se o path não corresponder ao seguinte formato: "repo/<file_repo_name>/.../<file_name>"
	if meta.Path[1] != meta.RepoNome || meta.Path[0] != "repo" || meta.Path[len(meta.Path)-1] != meta.Nome {
		return errors.New("não foi possivél criar o ficheiro pedido, erros no armazenamento descrito na meta data")
	}

	// Define o filtro a usar na procura de informação na BD
	filter := bson.M{"hash": meta.Hash}
	// Documento e repo onde procurar o repo
	collection := endpointfuncs.MongoClient.Database("documentacao").Collection("files-meta-data")
	cntx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	// Procura pot um registo com a mesma hash (registo igual)
	err := collection.FindOne(cntx, filter, options.FindOne()).Err()
	defer cancel()
	if err == nil {
		return errors.New("não foi possivél criar o ficheiro pedido, esse ficheiro já existe")
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
	fmt.Println("meta: ", meta, "\nerr: ", err)
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

func CriarMetaHash(metaData map[string]interface{}) (string, error) {
	// encodifica a meta data para []byte (em formato JSON)
	x, err := json.Marshal(metaData)
	if err != nil {
		return "", err
	}
	// Adiciona a hash o valor convertido para []byte
	return fmt.Sprintf("%x", sha1.Sum(x)), nil
}
