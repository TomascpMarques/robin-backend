package loggers

import (
	"log"
	"os"
)

// LoginResolverLogger - loger para os resolvers da schema GraphQL
var LoginResolverLogger = log.New(os.Stdout, "Auth-Resolver....(*) ", log.LstdFlags)

// LoginRedisLogger - logger para o o tratamento e criação do cliente que liga ao serviço redis
var LoginRedisLogger = log.New(os.Stdout, "Redis-Setup......[*] ", log.LstdFlags)

// LoginDbFuncsLogger - logger para o handeling de funções relacionadas á bd
var LoginDbFuncsLogger = log.New(os.Stdout, "DBIndexing.......<*> ", log.LstdFlags)

// LoginOperacoesBDLogger - logger para as operações relacionadas á bd
var LoginOperacoesBDLogger = log.New(os.Stdout, "Operações-BD.....|*| ", log.LstdFlags)

// LoginServerErrorLogger - Logger para erros do servidor http
var LoginServerErrorLogger = log.New(os.Stdout, "Server-Error-Log !*! ", log.LstdFlags)

// LoginAuthLogger - Logger para as funções de autenticação
var LoginAuthLogger = log.New(os.Stdout, "Login-Auth...... !*! ", log.LstdFlags)
