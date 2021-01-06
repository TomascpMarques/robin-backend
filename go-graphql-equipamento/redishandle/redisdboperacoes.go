package redishandle

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
)

var operacoesBDLogger = log.New(os.Stdout, "Operações-BD.....|*| ", log.LstdFlags)
var bdContentLimitter = "[+] Conteúdo Encontrado: \n--*<INICÍO>*--\n%s\n--*<FIM>*--\n"

/*
SetRegistoBD - Inssere o registo passado na base de dados.
---
Params
	cr - redis.Client / cliente redis a usar
	registo - RegistoRedisDB / registo a insserir
*/
func SetRegistoBD(cr *redis.Client, registo RegistoRedisDB) {
	err := cr.Set(context.Background(), registo.Key, registo.Valor, registo.Expira).Err()
	if err != nil {
		operacoesBDLogger.Panic("[!] Erro ao insserir registo na base de dados")
		return
	}
	operacoesBDLogger.Printf("[+] Registo insserido com ID: %v\n", registo.Key)
	operacoesBDLogger.Printf(bdContentLimitter, registo)
}

/*
GetRegistoBD - Busca um registo na BD, através da redis-key insserida.
---
Params
	cr - redis.Client / cliente redis a usar
	keyDoRegisto - string / key do registo a procurar
*/
func GetRegistoBD(cr *redis.Client, keyDoRegisto string) (string, error) {
	// Escreve no ecrã o registo insserido para verificação da insserção
	// e visualização do novo registo
	registo, getErr := cr.Get(context.Background(), keyDoRegisto).Result()
	if getErr != nil {
		operacoesBDLogger.Panicf("[!] Erro ao buscar pelo registo de key<%v>: %v", keyDoRegisto, getErr)
		return "null", getErr
	}
	operacoesBDLogger.Printf("[$] Conteudo do Registo <%v>:", keyDoRegisto)
	operacoesBDLogger.Printf(bdContentLimitter, registo)

	return registo, nil
}

/*
BuscarKeysVerificarResultado - Procura e devolve as keys relacionadas ao tipo de registo fornecido
							   existentes na redisDB utilizada, pelo client redis fornecido.
---
Params
	contexto context.Context - contexto a executar a função
	clienteRedis *redis.Client - cliente ligado a uma redisDB a utilisar
	tiporegisto string - o tipo de registo associado ás keys procuradas
*/
func BuscarKeysVerificarResultado(contexto context.Context, clienteRedis *redis.Client, tiporegisto string) []string {
	// Retorna só a lista com as keys do tipo de registo especificado
	// se o valor passado for "" retorna todas as keys
	keys := clienteRedis.Keys(contexto, tiporegisto+`*`).Val()
	// Veridfica se têm algum registo na DB alvo
	if len(keys) == 0 {
		dbFuncsLogger.Printf("[!!] Aviso: Lista de keys para <%v> vaiza (nil) Valor enviado: %v0\n", tiporegisto, tiporegisto)
		return append(make([]string, 1), tiporegisto+"0")
	}

	return keys
}

/*
ConversaoIDStringInt - Converte os digitos contidos no id para Ints
---
Params
	conteudo string - conteudo a converter
*/
func ConversaoIDStringInt(conteudo string) int {
	keyInt, err := strconv.Atoi(conteudo)
	if err != nil {
		dbFuncsLogger.Panic("[!] Erro: A Conversão falhou")
		return 0
	}
	return keyInt
}
