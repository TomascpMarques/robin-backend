package redishandle

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/tomascpmarques/PAP/backend/robinservicologin/loggers"
)

const (
	defaultRedisPort = "6379" // porta base onde o serviçio redis está exposto
	defaultDB        = 0
	defaultPassword  = "Pg+V@j+Z9gKj88=-?dSk" // Pg+V@j+Z9gKj88=-?dSk
	defaultUsername  = "admin"                // user para a conexão á base de dados, não o utilisador admin do sistema
)

// DefClienteRedis -
type DefClienteRedis struct {
	Addres   string // Endereço do serviço redis
	Port     string // Porta onde o serviço redis está exposto
	Password string // Password do utilizador a usar na autenticação
	User     string // Username do utilizador a usar na autenticação
	DB       int    // Base de dados a usar para as operações de base de dados
}

//AddressRed endereço do serviço redis
var AddressRed = os.Getenv("AUTH_SERVER_REDIS_PORT")

//PortRed porta onde o serviço redis está a correr
var PortRed = os.Getenv("REDISPORT")

//PasswordRed password para a autenticação na redis bd
var PasswordRed = os.Getenv("REDIS_USER1_PASS")

// UserRed Utilisador a utilizar na autenticação no serviço redis
var UserRed = os.Getenv("REDIS_USER1_NAME")

// DBRed a base de dados a concetar pelo cliente
var DBRed, _ = strconv.Atoi(os.Getenv("REDISDB"))

var redisLogger = loggers.LoginResolverLogger

/*
NovoClienteRedis Cria um novo cliente redis para conectar ao serviço redis
---
Params:
	addres - String Endereço onde o serviçoo está a correr
	port - String Porta onde o serviço está desponível
	db - Int Indica se vai usar a data-base default do redis
*/
func NovoClienteRedis(addres, port, password, username string, db int) redis.Client {
	// verifica as variaveis env passadas
	// e define valores default se não forem defenidos valores por vars env
	if port == "" {
		port = defaultRedisPort
	}
	if password == "" {
		password = defaultPassword
	}
	if username == "" {
		password = defaultUsername
	}

	// Sleep ajuda a que a redis-bd inicie antes de tentar conectar
	time.Sleep(time.Second * 10)
	// aplica as defenições passadas nos argumentos da função
	client := redis.NewClient(&redis.Options{
		Addr:     string(addres + ":" + port),
		Password: password,
		Username: username,
		DB:       db,
	})

	// Define o tempo máximo que o client-setup deve esperar pela BD
	cntx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second*40))
	// verifica se o cliente está UP e funcional
	_, err := client.Ping(cntx).Result()
	defer cancel()
	if err != nil {
		redisLogger.Printf("[!] Erro: %v", err)
		redisLogger.Fatal()
	}

	redisLogger.Println("[$] Cliente Redis Criado")
	return *client
}
