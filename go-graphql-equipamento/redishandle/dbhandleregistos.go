package redishandle

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/go-redis/redis/v8"
)

// RegistoRedisDB Estrutura o registo a insserir ou retirar da BD
type RegistoRedisDB struct {
	Key    string
	Valor  interface{}   // Conteúdo do registo
	Expira time.Duration // Expiração do registo
}

var dbFuncsLogger = log.New(os.Stdout, "DBIndexing.......<*> ", log.LstdFlags)

/*
extrairIDMaisRecente - Extrai o id mais recente da lista fornecida
---
Params
	listaIDs *[]string - lista de ids a usar para extração do maior index
*/
func extrairIDMaisRecente(listaIDs *[]string) string {
	// Cria o padrão para retirar os digitos do ID
	padrao := regexp.MustCompile(`\d+$`)
	indexMaior := ""
	// Lê todos os ids e procura o maior
	for _, v := range *listaIDs {
		if padrao.FindString(v) > indexMaior {
			indexMaior = padrao.FindString(v)
		}
	}
	return indexMaior
}

/*
DBCriadorID Retorna o tamanho da lista de keys na base de dados para criar o id/key do novo registo
---
Params
	rc - ponteiro(*) redis.Client utiliza o cliente conectado a bd para utilisar na função
*/
func DBCriadorID(clienteRedis *redis.Client, tiporegisto string) string {
	// Retorna só a lista com as keys
	// Não devolve o padrão utilizado
	ids := BuscarKeysVerificarResultado(context.Background(), clienteRedis, tiporegisto)

	// Extrai o ID mais recente da base de dados
	indexMaior := extrairIDMaisRecente(&ids)

	// converte o indice do id para uma integer
	resultado := ConversaoIDStringInt(indexMaior)

	// Devolve o maior id +1, relacionando aos ids existentes
	return string(tiporegisto + fmt.Sprint(resultado+1))
}

/*
FormatarValorParaJSON Recebe uma struct e devolve o equivalente em JSON
---
Params
	conteudo - interface{} Conteúdo a ser traduzido
	(interface para ser compativél com todo o tipo de estrutura de dados)
*/
func FormatarValorParaJSON(conteudo interface{}) []byte {
	// Tradução do input da mutation para json através da indentação do mesmo
	valorTraduzido, err := json.MarshalIndent(conteudo, "", "\t")
	if err != nil {
		dbFuncsLogger.Panic("[!] Erro ao converter valor insserido para JSON")
		return nil
	}
	return valorTraduzido
}

/*
CriaEstruturaRegisto -
---
Params
	redisClienteDB - *redis.Client / Ponteiro ao cliente redis a utilizar
	regsitoCorpo - interface{} / corpo/conteúdo do registo a insserir
*/
func (registo *RegistoRedisDB) CriaEstruturaRegisto(redisClienteDB *redis.Client, regsitoCorpo interface{}, tipoRegisto string) {
	// Criação de chave/index relacional aos existentes
	registo.Key = DBCriadorID(redisClienteDB, tipoRegisto)
	// Tradução do input da mutation para json através da indentação do mesmo
	registo.Valor = FormatarValorParaJSON(regsitoCorpo)
	// Defenição do tempo de expiração da key do registo (0 = não expira)
	registo.Expira = time.Duration(0)
	dbFuncsLogger.Printf("[.] Registo Criado\n")
}
