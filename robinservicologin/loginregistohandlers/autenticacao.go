package loginregistohandlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/tomascpmarques/PAP/backend/robinservicologin/loggers"
	"github.com/tomascpmarques/PAP/backend/robinservicologin/redishandle"
)

var hmacSecret = hmac.New(sha256.New, []byte(`SUPPER_SECRET_DEVELOPMENT_KEY`)).Sum(nil)

//RedisClientDB -
var RedisClientDB = redishandle.NovoClienteRedis(
	os.Getenv("REDISADDRESS"),
	os.Getenv("AUTH_SERVER_REDIS_PORT"),
	"Pg+V@j+Z9gKj88=-?dSk",
	"admin",
	0,
)

// Se não, cria o utilizador admin com as crdênciais default
// Verifica se o utilisador admin já existe ou não
var _ = VerificarAdminFirstBoot()

// VerificarTokenUser -
func VerificarTokenUser(userToken string) string {
	token, err := jwt.Parse(userToken, func(token *jwt.Token) (interface{}, error) {
		// valida o metodo de assinatura da key
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Metodo de assinatura inesperado: %v", token.Header["alg"])
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

// Login - Recebe dois parametros, o username e a passwd, cria uma token com esses dados e compara
func Login(user string, passwd string, perms int) map[string]interface{} {
	returnVal := make(map[string]interface{}, 0)

	// Cria uma struc user, com os dados passados nos parametros
	novoUser := CriarNovoUser(user, passwd, perms)

	// Cria um token de utilisador a partir dos dados fornecidos
	novoUserJWTToken, err := novoUser.CriarUserJWT().SignedString(hmacSecret)
	if err != nil {
		loggers.LoginAuthLogger.Println("Error: ", err)
		returnVal["erro"] = err
		return returnVal
	}

	// Busca o registo de utilisador que se está a usar para fazer login
	// E compara as credenciais com as do utilisador fornecido
	if reg := GetUserParaValorStruct(user); reg.Password != passwd && reg.Permissoes != perms {
		loggers.LoginAuthLogger.Println("Error: ", "credeenciais inválidas")
		returnVal["erro"] = "Credenciais inválidas"
		return returnVal
	}

	// Loga que o utilisador XXXX iniciou sessão
	loggers.LoginAuthLogger.Println("Utilizador, ", user, ", iniciou sessão")
	returnVal["token"] = novoUserJWTToken
	return returnVal
}

// Registar -
func Registar(user string, password string, token string, perms int) map[string]interface{} {
	returnVal := make(map[string]interface{}, 0)

	// Busca o registo do admin
	registoAdmin, err := redishandle.GetRegistoBD(&RedisClientDB, "admin", 0)
	if err != nil {
		loggers.LoginAuthLogger.Println("Erro ao buscar jwt do utilisador necessário.")
		returnVal["error"] = "não foi possível encontrar o utilizador necessário"
		return returnVal
	}
	// Transforma o contéudo do registo do admin em uma estrutura User
	var userRegAdmin User
	err = json.Unmarshal([]byte(registoAdmin), &userRegAdmin)
	if err != nil {
		loggers.LoginAuthLogger.Println("Erro ao descodificar registo do utilizador.")
		returnVal["error"] = "Erro interno"
		return returnVal
	}
	// Se a token não for igual ao do admin não se regista nenhumuser novo
	if userRegAdmin.JWT != token {
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
		loggers.LoginAuthLogger.Println("Error: ", err)
		returnVal["token"] = jwtTokenSigned
		return returnVal
	}

	returnVal["error"] = "Credenciais inválidas"
	return returnVal
}
