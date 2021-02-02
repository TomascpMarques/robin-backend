package redishandle

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/tomascpmarques/PAP/backend/robinservicologin/loggers"
)

var operacoesBDLogger = loggers.LoginOperacoesBDLogger

//BDContentLimitter - limitador de conteúdo da bd.
var BDContentLimitter = "[+] Conteúdo Encontrado: \n--*<INICÍO>*--\n%s\n--*<FIM>*--\n"

/*
SetRegistoBD - Inssere o registo passado na base de dados.
---
Params
	cr - redis.Client / cliente redis a usar
	registo - RegistoRedisDB / registo a insserir
*/
func SetRegistoBD(cr *redis.Client, registo RegistoRedisDB, debugg int) {
	err := cr.Set(context.Background(), registo.Key, registo.Valor, registo.Expira).Err()
	if err != nil {
		operacoesBDLogger.Panic("[!] Erro ao insserir registo na base de dados, Erro: ", err)
		return
	}
	operacoesBDLogger.Printf("[+] Registo insserido com ID: %v\n", registo.Key)
	if debugg == 1 {
		operacoesBDLogger.Printf(BDContentLimitter, registo)
	}
}

/*
DelRegistoBD - Apaga um ou mais registos passados para função como um vetor de strings
---
Params
	cr - redis.Client / cliente redis a usar
	regID - []string / Ids dos registos a apagar
*/
func DelRegistoBD(cr *redis.Client, regID ...string) error {
	pattern := regexp.MustCompile(`\d+$`)
	err := cr.Del(context.Background(), regID...)
	if pattern.FindString(err.String()) != fmt.Sprintf("%v", len(regID)) {
		operacoesBDLogger.Panicf("[!] Erro ao apagar o registo de keys: %v", regID)
		return errors.New("[!!] ID de registo inválido")
	}

	return nil
}

/*
GetRegistoBD - Busca um registo na BD, através da redis-key insserida.
---
Params
	cr - redis.Client / cliente redis a usar
	keyDoRegisto - string / key do registo a procurar
*/
func GetRegistoBD(cr *redis.Client, keyDoRegisto string, debugg int) (string, error) {
	// Escreve no ecrã o registo insserido para verificação da insserção
	// e visualização do novo registo
	registo, getErr := cr.Get(context.Background(), keyDoRegisto).Result()
	if getErr != nil {
		operacoesBDLogger.Printf("[!] Erro ao buscar pelo registo de key <%v> : %v", keyDoRegisto, getErr)
		erroNaProcura := "Sem registo para id: " + keyDoRegisto
		return "null", errors.New(erroNaProcura)
	}
	operacoesBDLogger.Printf("[$] ID do Registo <%v>:", keyDoRegisto)
	if debugg == 1 {
		operacoesBDLogger.Printf(BDContentLimitter, registo)
	}

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
		operacoesBDLogger.Printf("[!!] Aviso: Lista de keys para <%v> vaiza (nil) Valor enviado: %v0\n", tiporegisto, tiporegisto)
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
		operacoesBDLogger.Panic("[!] Erro: A Conversão falhou")
		return 0
	}
	return keyInt
}
