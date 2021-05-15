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
	if defs.VerificarRepoParaContribuicao(repo) {
		return errors.New("o repositório pedido não existe nas contribuições")
	}

	// Dá inssert no array ficheiros do objeto "$", que está dentro do array  contribuicoes
	contribuicao := bson.M{"$push": bson.M{"contribuicoes.$.ficheiros": ficheiro}}

	// Atualização da info do user
	resultado, err := defs.Colecao.UpdateOne(defs.Cntxt, defs.Filter, contribuicao)
	defs.CancelFunc()

	if err != nil {
		return err
	}
	// Se existir um erro na operação, o nº de files modificados é menor que 1
	if resultado.ModifiedCount < 1 {
		return errors.New("nenhum ficheiro foi modificado")
	}
	return nil
}

// RemoverContribuicaoFile Remove contribuições (nomes dos ficheiros), da info do user
func (defs *MongoDBOperation) RemoverContribuicaoFile(repo string, ficheiro string) error {
	// Contribuição a inserir
	contribuicao := bson.M{"$pull": bson.M{"contribuicoes.$.ficheiros": ficheiro}}

	// Inserção da contribuição na info do user
	resultado, err := defs.Colecao.UpdateOne(defs.Cntxt, defs.Filter, contribuicao)
	defs.CancelFunc()

	// Error handeling
	if err != nil {
		return err
	}
	if resultado.ModifiedCount < 1 {
		return errors.New("nenhum ficheiro foi modificado")
	}
	return nil
}

// RemoverRepoContribuicao Remove contribuições (repos), da info do user
func (defs *MongoDBOperation) RemoverRepoContribuicao(repo string) error {
	// Contribuição a inserir
	operacaoDrop := bson.M{"$pull": bson.M{"contribuicoes": bson.M{"reponome": repo}}}

	// Inserção da contribuição na info do user
	resultado, err := defs.Colecao.UpdateOne(defs.Cntxt, defs.Filter, operacaoDrop)
	defs.CancelFunc()

	// Error handeling
	if err != nil {
		return err
	}
	if resultado.ModifiedCount < 1 {
		return errors.New("nenhum repo foi largado")
	}
	return nil
}

// VerificarRepoParaContribuicao Verifica se o repo existe antes de inserir o nome do ficheiro no mesmo
func (defs *MongoDBOperation) VerificarRepoParaContribuicao(repoNome string) bool {
	return defs.Colecao.FindOne(defs.Cntxt, defs.Filter).Err() != nil
}

// CriarContribuicaoStruct Inicializa os valores das contribuições dos repos a zero, evita bugs
func CriarContribuicaoStruct(nome string) (cntrb resolvedschema.Contribuicoes) {
	cntrb.RepoNome = nome
	cntrb.Ficheiros = make([]string, 0)
	return
}

// CriarRepoContribuicoes Cria o repo no user-info para armazenar contribuições, devolve um erro se falhar
func (defs *MongoDBOperation) CriarRepoContribuicoes(repoNome string) error {
	// Inicializa as contribuições a zeros
	contrib := CriarContribuicaoStruct(repoNome)
	// Query Mongo db para inserir em ultimo lugar a contribuição
	atualizacao := bson.M{"$push": bson.M{"contribuicoes": contrib}}

	// Atualiza a lista das contribuições com os conteudos das linhas anteriores
	_, err := defs.Colecao.UpdateOne(defs.Cntxt, defs.Filter, atualizacao)
	if err != nil {
		return err
	}
	return nil
}
