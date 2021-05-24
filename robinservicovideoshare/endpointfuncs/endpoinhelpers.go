package endpointfuncs

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"time"

	"github.com/tomascpmarques/PAP/backend/robinservicovideoshare/loggers"
	"github.com/tomascpmarques/PAP/backend/robinservicovideoshare/resolvedschema"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDBOperation Struct com o setup minímo para fazer uma oepração na BDs
type MongoDBOperation struct {
	Colecao    *mongo.Collection
	Cntxt      context.Context
	CancelFunc context.CancelFunc
	Filter     interface{}
}

// Setup Evita mais lihas desnecessárias e repetitivas para poder-se usar a coleção necessaria
func SetupColecao(dbName, collName string) (defs MongoDBOperation) {
	defs.Colecao = MongoClient.Database(dbName).Collection(collName)
	defs.Cntxt, defs.CancelFunc = context.WithTimeout(context.Background(), time.Second*10)
	return
}

// VerificarCamposBase Verifica se os vlores para keys do mapa <params> se encontram no array <campos>
func VerificarCamposBase(params map[string]interface{}, campos []string) error {
	// Verifica se têm os campos obrigatórios defenidos
	for _, v := range campos {
		if _, existe := params[v]; !existe {
			return fmt.Errorf("o campo <%s>, não está presente na info da videoshare", v)
		}
	}

	// Retorna nil se não houver erros
	return nil
}

func VerificarSearchParams(params map[string]interface{}) error {
	// Verifica se as keys obrigatorias existem em params
	if err := VerificarCamposBase(params, []string{"quanti", "params"}); err != nil {
		return err
	}

	return nil
}

func VerificarVideoShareMetaData(videoShare map[string]interface{}) error {
	camposObrgt := []string{
		"url",
		"tema",
		"titulo",
		"criador",
	}

	// Verifica os campos base
	if err := VerificarCamposBase(videoShare, camposObrgt); err != nil {
		return err
	}

	// Verifica se o tema do video é válido
	if len(videoShare["tema"].(string)) < 3 {
		return errors.New("o tema do video é demasiado curto, deve ter no minímo 3 caracteres")
	}

	// Verifica o tamanho do título
	if len(videoShare["titulo"].(string)) < 4 {
		return errors.New("o título do video é demasiado curto, deve ter no minímo 4 caracteres")
	}

	// Verifica que o url é válido
	url := (videoShare["url"].(string))
	regex := regexp.MustCompile(`(?m)https://youtu\.be/[a-zA-Z0-9_]+`).FindAllString(url, -1)
	if reflect.ValueOf(regex).IsZero() {
		return errors.New("o link fornecido não é válido")
	}

	// If all good returb nil (no error)
	return nil
}

// TrimURL Larga a parte desnecessária do url, >https://youtu\.be/<
func TrimURL(url string) string {
	return url[18:]
}

// Adiciona os dados da video-share à base de dados
func AdicionarVideoShareDB(videoShare *resolvedschema.Video) error {
	colecao := SetupColecao("videoshares", "videos")

	// Inserção do registo na BD
	result, err := colecao.Colecao.InsertOne(colecao.Cntxt, videoShare, options.InsertOne())
	defer colecao.CancelFunc()
	if err != nil {
		loggers.DbFuncsLogger.Println("Não foi possivél adicionar a video-share na base de dados: ", err.Error())
		return err
	}

	// Verifica se o id inserido é != de nil, para verificar a inserção do registo
	if result.InsertedID == "" {
		loggers.DbFuncsLogger.Println("Não foi possivél criar o registo na base de dados")
		return errors.New("não foi possivél criar o registo da video-share na base de dados")
	}

	return nil
}

// GetVideoShareWithParams Procura todos os registos de videoshares que igualem aos params passados
func GetVideoShareWithParams(params *resolvedschema.VideoSearchParams) ([]resolvedschema.Video, error) {
	// Setup da coleção a usar & search params
	colecao := SetupColecao("videoshares", "videos")
	colecao.Filter = params.Params
	opcoesBusca := options.Find().SetAllowPartialResults(true).SetLimit(int64(params.Quanti))

	cursor, err := colecao.Colecao.Find(colecao.Cntxt, colecao.Filter, opcoesBusca)
	defer colecao.CancelFunc()
	if err != nil {
		return nil, err
	}

	// Obtêm todos os documentos retornados e colocaos num array de maps
	var results []resolvedschema.Video
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}

	// Verifica se encontrou algum resultado
	if reflect.ValueOf(results).IsZero() {
		return nil, errors.New("sem registos para esse valor")
	}

	return results, nil
}
