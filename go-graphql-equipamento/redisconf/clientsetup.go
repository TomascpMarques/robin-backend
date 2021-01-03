package redisconf

import (
	"context"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

var redisLogger = log.New(os.Stdout, "Redis-Setup......[*] ", log.LstdFlags)

const (
	defaultRedisAddress = "localhost"
	defaultRedisPort    = "6379"
	defaultDB           = 0
	defaultPassword     = ""
)

//AddressRed endereço do serviço redis
var AddressRed = os.Getenv("REDISADDRESS")

//PortRed porta onde o serviço redis está a correr
var PortRed = os.Getenv("REDISPORT")

//PasswordRed password para a autenticação na redis bd
var PasswordRed = os.Getenv("REDISPASSWORD")

/*
NovoClienteRedis Cria um novo cliente redis para conectar ao serviço redis
---
Params:
	addres - String Endereço onde o serviçoo está a correr
	port - String Porta onde o serviço está desponível
	db - Int Indica se vai usar a data-base default do redis
*/
func NovoClienteRedis(addres, port, password string) redis.Client {
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

	// aplica as defenições passadas nos argumentos da função
	client := redis.NewClient(&redis.Options{
		Addr:     string(addres + ":" + port),
		Password: password,
		DB:       0,
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
