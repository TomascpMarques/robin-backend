package loginregistohandlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/tomascpmarques/PAP/backend/robinservicologin/loggers"
	"github.com/tomascpmarques/PAP/backend/robinservicologin/redishandle"
)

var hmacSecret = hmac.New(sha256.New, []byte(`SUPPER_SECRET_DEVELOPMENT_KEY`)).Sum(nil)

//RedisClientDB -
var RedisClientDB = redishandle.NovoClienteRedis(
	//os.Getenv("REDISADDRESS"),
	//os.Getenv("AUTH_SERVER_REDIS_PORT"),
	"0.0.0.0",
	"6379",
	"Pg+V@j+Z9gKj88=-?dSk",
	"admin",
	0,
)

// Verifica se o utilisador admin já existe ou não
// Se não, cria o utilizador admin com as crdênciais default
var _ = VerificarAdminFirstBoot()

// VerificarTokenUser verifica se a token passada é válida, logo vê se já expirou
// se o modo de assinatura é o correto, e se o emissor é o servidor de autenticação
func VerificarTokenUser(userToken string) string {
	token, err := jwt.Parse(userToken, func(token *jwt.Token) (interface{}, error) {
		// valida o metodo de assinatura da key
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("metodo de assinatura inesperado: %v", token.Header["alg"])
		}

		// hmacSampleSecret é o []byte que contem o segredo de assinatura
		return hmacSecret, nil
	})
	// Se a token for assinada por outro metodo ou a key for diferente dá erro
	if err != nil {
		return fmt.Sprint(err)
	}

	// Verifica que a token é válida e assinada pelo servidor de login
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid && claims["iss"] == "Robin-Servico-Auth" {
		return "OK"
	}
	return "Token inválida ou expirada"
}

// VerificarTokenAdmin verifica tudo o que a função VerificarTokenUser verifica,
// e ainda verifica se o utilisador é o administrador
func VerificarTokenAdmin(userToken string) string {
	token, err := jwt.Parse(userToken, func(token *jwt.Token) (interface{}, error) {
		// valida o metodo de assinatura da key
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("metodo de assinatura inesperado: %v", token.Header["alg"])
		}

		// hmacSampleSecret é o []byte que contem o segredo de assinatura
		return hmacSecret, nil
	})
	// Se a token for assinada por outro metodo ou a key for diferente dá erro
	if err != nil {
		return fmt.Sprint(err)
	}

	// Verifica que a token é válida e assinada pelo servidor de login
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid && claims["iss"] == "Robin-Servico-Auth" &&
		claims["perms"].(float64) == 2 {
		return "OK"
	}
	return "Token inválida ou expirada"
}

// Login Recebe dois parametros, o username e a passwd, cria uma token com esses dados e compara com o utilisador pedido
// devolve uma token com o tempo de expiração de time.Now().Add(time.Minute * 40).Unix()
func Login(user string, passwd string) map[string]interface{} {
	returnVal := make(map[string]interface{})

	// Busca o registo de utilisador que se está a usar para fazer login
	utilizadorPedido, err := GetUserParaValorStruct(user)
	if err != nil {
		loggers.LoginAuthLogger.Println("Erro: ", err)
		returnVal["erro"] = err
		return returnVal
	}

	// Compára as credenciais com as do utilisador fornecido
	if utilizadorPedido.Password != passwd {
		loggers.LoginAuthLogger.Println("Error: ", "credeenciais inválidas")
		returnVal["erro"] = "Credenciais inválidas"
		return returnVal
	}

	// Cria um token de utilisador a partir dos dados fornecidos
	UserNewJWTToken, err := utilizadorPedido.CriarUserJWT().SignedString(hmacSecret)
	if err != nil {
		loggers.LoginAuthLogger.Println("Error: ", err)
		returnVal["erro"] = err
		return returnVal
	}

	// Loga que o utilisador XXXX iniciou sessão
	// E devolve a token, em como o utilisador está logado
	loggers.LoginAuthLogger.Println("Utilizador, ", user, ", iniciou sessão")
	returnVal["token"] = UserNewJWTToken
	return returnVal
}

// Registar utiliza os dados de utilisador base defenidos, cria e inssere na BD um utilisador novo, antes disso
// ela verifica se quem está a fazer o pedido é o administrador, só administradores podem registar utilisadores.
// Se todas as regras forem cumpridas, a função devolve a jwt token desse novo utilizador.
func Registar(user string, password string, perms int, token string) map[string]interface{} {
	returnVal := make(map[string]interface{})

	// Se a token não for igual ao do admin não se regista nenhumuser novo
	if VerificarTokenAdmin(token) != "OK" {
		loggers.LoginAuthLogger.Println("A token não têm permissões")
		returnVal["error"] = "A token não têm permissões"
		return returnVal
	}

	// Cria a jwt para o novo utilisador
	jwtTokenSigned, err := CriarNovoUser(user, password, perms).CriarUserJWT().SignedString(hmacSecret)
	if err != nil {
		loggers.LoginAuthLogger.Println("Error: ", err)
		returnVal["error"] = err
		return returnVal
	}

	// Cria a struct para o novo user
	novoUser := CriarNovoUser(user, password, perms)
	loggers.LoginAuthLogger.Println("Novo user: ", user)

	// Encodifica a info relacionada ao user, para um formato json
	novoUserJSON, err := json.Marshal(&novoUser)
	if err != nil {
		loggers.LoginAuthLogger.Println("Error: ", err)
		returnVal["error"] = err
		return returnVal
	}

	// Verifica se o utilisador que se quer criar já existe
	// Se já existir, não se devolve nenhuma jwt, nem se inssere nada na BD
	_, exists := redishandle.GetRegistoBD(&RedisClientDB, user, 0)
	if exists != nil {
		// Inssere o novo utilisador na bd se o utilisador não existir
		redishandle.SetRegistoBD(&RedisClientDB, redishandle.RegistoRedisDB{
			Key:    novoUser.Username,
			Valor:  novoUserJSON,
			Expira: 0,
		}, 0)
		loggers.LoginAuthLogger.Println("Registo adicionado com sucesso.")
		returnVal["token"] = jwtTokenSigned
		return returnVal
	}

	returnVal["error"] = "Credenciais inválidas ou utilizador já existente"
	return returnVal
}

// AtualizarUser atualiza os dados dos utilizador fornecido, depois de verificar a token fornecida
func AtualizarUser(user string, userInfo map[string]interface{}, token string) map[string]interface{} {
	returnVal := make(map[string]interface{})

	if VerificarTokenAdmin(token) != "OK" {
		loggers.LoginAuthLogger.Println("Token inválida.")
		returnVal["err"] = "Token inválida ou expirada"
		return returnVal
	}

	userAtualizar, err := GetUserParaValorStruct(user)
	if err != nil {
		loggers.LoginAuthLogger.Println("Erro: ", "Sem registo para <", user, ">")
		returnVal["erro"] = err
		return returnVal
	}

	if userInfo["user"] != nil && userInfo["user"] != userAtualizar.Username {
		userAtualizar.Username = userInfo["user"].(string)
	}
	if userInfo["pass"] != nil && userInfo["pass"] != userAtualizar.Password {
		userAtualizar.Password = userInfo["pass"].(string)
	}
	if userInfo["perms"] != nil && userInfo["perms"] != userAtualizar.Permissoes {
		userAtualizar.Permissoes = userInfo["perms"].(int)
	}

	userAtualizadoJSON, err := json.Marshal(&userAtualizar)
	if err != nil {
		loggers.LoginAuthLogger.Println("Error: ", err)
		returnVal["error"] = err
		return returnVal
	}

	err = redishandle.DelRegistoBD(&RedisClientDB, user)
	if err != nil {
		loggers.LoginAuthLogger.Println("Error: ", err)
		returnVal["error"] = err
		return returnVal
	}

	redishandle.SetRegistoBD(&RedisClientDB, redishandle.RegistoRedisDB{
		Key:    userAtualizar.Username,
		Valor:  userAtualizadoJSON,
		Expira: 0,
	}, 0)

	returnVal["Menssagem"] = "Sucesso ao alterar dados."
	return returnVal
}
