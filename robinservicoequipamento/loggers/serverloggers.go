package loggers

import (
	"log"
	"os"
)

// ResolverLogger - loger para os resolvers da schema GraphQL
var ResolverLogger = log.New(os.Stdout, "Equipamento-Resolver (*) ", log.LstdFlags)

// RedisLogger - logger para o o tratamento e criação do cliente que liga ao serviço redis
var RedisLogger = log.New(os.Stdout, "Redis-Setup......[*] ", log.LstdFlags)

// DbFuncsLogger - logger para o handeling de funções relacionadas á bd
var DbFuncsLogger = log.New(os.Stdout, "DBIndexing.......<*> ", log.LstdFlags)

// OperacoesBDLogger - logger para as operações relacionadas á bd
var OperacoesBDLogger = log.New(os.Stdout, "Operações-BD.....|*| ", log.LstdFlags)
