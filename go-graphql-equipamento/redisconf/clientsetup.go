package redisconf

import (
	"context"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

var redisLogger = log.New(os.Stdout, "Redis-Setup [*] ", log.LstdFlags)

const (
	defaultRedisAddress = "localhost"
	defaultRedisPort    = "6379"
	defaultDB           = 0
	defaultPassword     = ""
)

/*
NovoClienteRedis Cria um novo cliente redis para conectar ao serviço redis
Params:
	addres - String Endereço onde o serviçoo está a correr
	port - String Porta onde o serviço está desponível
	db - Int Indica se vai usar a data-base default do redis
*/
func NovoClienteRedis(addres, port, password string) redis.Client {
	if addres == "" {
		addres = defaultRedisAddress
	}
	if port == "" {
		port = defaultRedisPort
	}
	if password == "" {
		password = defaultPassword
	}

	client := redis.NewClient(&redis.Options{
		Addr:     string(addres + ":" + port),
		Password: password,
		DB:       0,
	})

	// verifica se o cliente está UP e funcional
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		redisLogger.Printf("Error: %v", err)
		redisLogger.Fatal()
	}

	redisLogger.Println("Cliente Redis Criado")
	return *client
}
