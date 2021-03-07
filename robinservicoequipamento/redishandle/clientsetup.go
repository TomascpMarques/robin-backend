package redishandle

import (
	"context"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/tomascpmarques/PAP/backend/robinservicoequipamento/loggers"
)

const (
	defaultRedisAddress = "0.0.0.0"
	defaultRedisPort    = "6379"
	defaultDB           = 0
	defaultPassword     = "" // Pg+V@j+Z9gKj88=-?dSk
	defaultUsername     = "" // admin
)

// DefClienteRedis -
type DefClienteRedis struct {
	Addres   string
	Port     string
	Password string
	User     string
	DB       int
}

//AddressRed endereço do serviço redis
var AddressRed = os.Getenv("REDISADDRESS")

//PortRed porta onde o serviço redis está a correr
var PortRed = os.Getenv("REDISPORT")

//PasswordRed password para a autenticação na redis bd
var PasswordRed = os.Getenv("REDISPASSWORD")

// UserRed Utilisador a utilizar na autenticação no serviço redis
var UserRed = os.Getenv("REDISUSER")

// DBRed a base de dados a concetar pelo cliente
var DBRed, _ = strconv.Atoi(os.Getenv("REDISDB"))

var redisLogger = loggers.RedisLogger

/*
NovoClienteRedis Cria um novo cliente redis para conectar ao serviço redis
---
Params:
	addres - String Endereço onde o serviçoo está a correr
	port - String Porta onde o serviço está desponível
	db - Int Indica se vai usar a data-base default do redis
*/
func NovoClienteRedis(addres, port, password, username string, db int) redis.Client {
	// checks for passed env variables
	// and sets default if none are passed
	if addres == "" {
		addres = defaultRedisAddress
	}
	if port == "" {
		port = defaultRedisPort
	}
	if password == "" {
		password = defaultPassword
	}
	if username == "" {
		password = defaultUsername
	}

	// aplica as defenições passadas nos argumentos da função
	client := redis.NewClient(&redis.Options{
		Addr:     string(addres + ":" + port),
		Password: password,
		Username: username,
		DB:       db,
	})

	// verifica se o cliente está UP e funcional
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		redisLogger.Printf("[!] Erro: %v", err)
		redisLogger.Fatal()
	}

	redisLogger.Println("[$] Cliente Redis Criado")
	return *client
}
