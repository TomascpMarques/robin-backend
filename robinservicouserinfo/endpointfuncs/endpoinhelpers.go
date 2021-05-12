package endpointfuncs

import (
	"context"
	"errors"
	"time"

	"github.com/tomascpmarques/PAP/backend/robinservicouserinfo/resolvedschema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

// AdicionarContribuicao Metodo que adiciona contribuicoes à info do user
func (defs *MongoDBOperation) AdicionarContribuicao(repo string, ficheiro string) error {
	// Dá inssert no array ficheiros do objeto "$", que está dentro do array  contribuicoes
	contribuicao := bson.M{"$push": bson.M{"contribuicoes.$.ficheiros": ficheiro}}
	resultado, err := defs.Colecao.UpdateOne(defs.Cntxt, defs.Filter, contribuicao)
	defs.CancelFunc()
	if err != nil {
		return err
	}
	if resultado.ModifiedCount < 1 {
		return errors.New("nenhum ficheiro foi modificado")
	}
	return nil
}

func (defs *MongoDBOperation) RemoverContribuicao(repo string, ficheiro string) error {
	contribuicao := bson.M{"$pull": bson.M{"contribuicoes.$.ficheiros": ficheiro}}
	resultado, err := defs.Colecao.UpdateOne(defs.Cntxt, defs.Filter, contribuicao)
	defs.CancelFunc()
	if err != nil {
		return err
	}
	if resultado.ModifiedCount < 1 {
		return errors.New("nenhum ficheiro foi modificado")
	}
	return nil
}

func (defs *MongoDBOperation) VerificarRepoParaContribuicao(repoNome string) bool {
	return defs.Colecao.FindOne(defs.Cntxt, defs.Filter).Err() != nil
}

func CriarContribuicaoStruct(nome string) (cntrb resolvedschema.Contribuicoes) {
	cntrb.RepoNome = nome
	cntrb.Ficheiros = make([]string, 0)
	return
}

func (defs *MongoDBOperation) CriarRepoContribuicoes(repoNome string) error {
	contrib := CriarContribuicaoStruct(repoNome)
	atualizacao := bson.M{"$push": bson.M{"contribuicoes": contrib}}

	_, err := defs.Colecao.UpdateOne(defs.Cntxt, defs.Filter, atualizacao)
	if err != nil {
		return err
	}
	return nil
}
