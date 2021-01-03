package redisconf

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"

	"github.com/go-redis/redis/v8"
)

var indexLogger = log.New(os.Stdout, "DBIndexing.......<*> ", log.LstdFlags)

/*
IDHandlerBD Retorna o tamanho da lista de keys na base de dados para criar o id/key do novo registo
---
Params
	rc - ponteiro(*) redis.Client utiliza o clienete conectado a bd para utilisar na função
*/
func IDHandlerBD(rc *redis.Client) string {
	// Retorna só a lista com as keys
	// Não devolve o padrão utilizado
	keys := rc.Keys(context.Background(), `*`).Val()
	indexLogger.Println(keys)
	if len(keys) == 0 {
		indexLogger.Println("[!] Aviso: Lista de keys vaiza (nil) Valor enviado: 0")
		return "0"
	}

	// Cria o padrão para retirar os digitos do ID
	padrao := regexp.MustCompile(`\d+`)
	indexMaior := ""
	// Lê todos os ids e procura o maior
	for _, v := range keys {
		if padrao.FindString(v) > indexMaior {
			indexMaior = padrao.FindString(v)
		}
	}

	// converte a string com a parte numérica do id para uma integer
	resultado, err := strconv.Atoi(indexMaior)
	if err != nil {
		indexLogger.Panic("[!] Erro: A Conversão falhou")
		return ""
	}

	// Devolve o maior id +1, relacionando aos ids existentes
	// Logo resolve o problema de continuidade de IDs incrementados
	return fmt.Sprint(resultado + 1)
}
