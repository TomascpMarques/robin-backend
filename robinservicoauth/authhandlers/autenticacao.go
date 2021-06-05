package authhandlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"

	"github.com/tomascpmarques/PAP/backend/robinservicoauth/loggers"
	"github.com/tomascpmarques/PAP/backend/robinservicoauth/redishandle"
)

// Chave usada para encriptar tokens de autorização do serviço
var assinaturaSecretaServer = hmac.New(sha256.New, []byte(`SUPPER_SECRET_DEVELOPMENT_KEY`)).Sum(nil)

//RedisClientDB  Setup do cliente do serviço redisDB
var RedisClientDB = redishandle.NovoClienteRedis(
	//os.Getenv("REDISADDRESS"),
	//os.Getenv("AUTH_SERVER_REDIS_PORT"),
	"0.0.0.0",              // Endereço perferido do serviço
	"6379",                 // Porta default do redis
	"Pg+V@j+Z9gKj88=-?dSk", // Pass do cliente admin do serv. redis
	"admin",                // Nome do user do serviço redis
	0,
)

// Verifica se o utilisador admin já existe ou não
// Se não, cria o utilizador admin com as crdênciais default
var _ = VerificarAdminFirstBoot()

// Login Recebe dois parametros, o username e a passwd, cria uma token com esses dados e compara com o utilisador pedido
// devolve uma token com o tempo de expiração de time.Now().Add(time.Minute * 40).Unix()
func Login(user string, passwd string) (retorno map[string]interface{}) {
	retorno = make(map[string]interface{})

	// Busca o registo de utilisador que se está a usar para fazer login
	utilizadorPedido, err := GetUserParaValorStruct(user)
	if err != nil {
		loggers.LoginAuthLogger.Println("Erro: ", err)
		retorno["erro"] = "O utilisador pedido não existe."
		return
	}

	// Compára as credenciais com as do utilisador fornecido
	if utilizadorPedido.Password != passwd {
		loggers.LoginAuthLogger.Println("Error: ", "credeenciais inválidas")
		retorno["erro"] = "Credenciais inválidas"
		return
	}

	// Cria um token de utilisador a partir dos dados fornecidos
	novaTokenLogin, err := utilizadorPedido.CriarJWTAuth().SignedString(assinaturaSecretaServer)
	if err != nil {
		loggers.LoginAuthLogger.Println("Error: ", err)
		retorno["erro"] = err
		return
	}

	// Loga que o utilisador XXXX iniciou sessão
	// E devolve a token, em como o utilisador está logado
	loggers.LoginAuthLogger.Println("Utilizador, ", user, ", iniciou sessão")
	retorno["token"] = novaTokenLogin
	return
}

// Registar utiliza os dados de utilisador base defenidos, cria e inssere na BD um utilisador novo, antes disso
// a função verifica que quem está a fazer o pedido é o administrador do serviço, só administradores podem registar utilizadores.
// Se todas as regras forem cumpridas, a função devolve a jwt token desse novo utilizador.
func Registar(user string, password string, perms int, token string) (retorno map[string]interface{}) {
	retorno = make(map[string]interface{})

	// Se a token não for igual ao do admin não se regista um user novo
	if VerificarTokenAdmin(token) != "OK" {
		loggers.LoginAuthLogger.Println("A token não têm permissões")
		retorno["error"] = "A token não têm permissões"
		return
	}

	// Verifica se o utilisador que se quer criar já existe
	// Se já existir, não se devolve nenhuma jwt, nem se inssere nada na BD
	_, exists := redishandle.GetRegistoBD(&RedisClientDB, user, 0)
	if exists != nil {
		// Cria a struct para o novo user
		novoUser := CriarNovoUser(user, password, perms)

		// Encodifica a info relacionada ao user, para um formato json
		novoUserJSON, err := json.Marshal(&novoUser)
		loggers.LoginAuthLogger.Println("Novo user: ", user)
		if err != nil {
			loggers.LoginAuthLogger.Println("Error: ", err)
			retorno["error"] = err
			return
		}

		// Inssere o novo utilisador na bd se o utilisador não existir
		redishandle.SetRegistoBD(&RedisClientDB, redishandle.RegistoRedisDB{
			Key:    novoUser.Username,
			Valor:  novoUserJSON,
			Expira: 0,
		}, 0)
		loggers.LoginAuthLogger.Println("Registo adicionado com sucesso.")
		retorno["sucesso"] = true
		return
	}

	retorno["error"] = "Credenciais inválidas ou utilizador já existente"
	return
}
