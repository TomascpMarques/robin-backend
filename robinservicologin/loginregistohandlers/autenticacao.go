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

// RedisClientDB -
var RedisClientDB = redishandle.NovoClienteRedis(
	redishandle.AddressRed,
	os.Getenv("AUTH_SERVER_REDIS_PORT"),
	os.Getenv("REDIS_USER1_PASS"),
	os.Getenv("REDIS_USER1_NAME"), 0,
)

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
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid && claims["iss"] == "loginServer" {
		return "OK"
	}
	return "Token inválida ou expirada"
}

// Login - Recebe dois parametros, o username e a passwd, cria uma token com esses dados e compara
//
func Login(user string, pass string) map[string]interface{} {
	returnVal := make(map[string]interface{}, 0)

	// cria uma struc user, com os dados passados nos parametros
	novoUser := CriarNovoUser(user, pass)

	// Cria um token para esse utilisador
	novoUserJWTToken, err := novoUser.CriarUserJWT().SignedString(hmacSecret)
	if err != nil {
		loggers.LoginAuthLogger.Println("Error: ", err)
		returnVal["erro"] = err
		return returnVal
	}
	// Atribui a token ao utilizador
	novoUser.JWT = novoUserJWTToken
	loggers.LoginAuthLogger.Println("Token: ", novoUserJWTToken[:16], "...")

	// Busca o registo correspondente ao user passado nos parametros
	userCompare, err := redishandle.GetRegistoBD(&RedisClientDB, user, 0)
	if err != nil {
		loggers.LoginAuthLogger.Println("O utilisador não existe")
		returnVal["erro"] = "O utilisador não existe"
		return returnVal
	}
	// Cria a estrutura User para o registo, descodifica o conteúdo de valores json
	var registo User
	err = json.Unmarshal([]byte(userCompare), &registo)
	if err != nil {
		loggers.LoginAuthLogger.Println("Error: ", err)
		returnVal["erro"] = "Erro Interno"
		return returnVal
	}

	// Verifica se a plavra-passe insserida é a correta para o utilisador fornecido
	if registo.Password != pass {
		loggers.LoginAuthLogger.Println("Error: ", "credeenciais inválidas")
		returnVal["erro"] = "Credenciais inválidas"
		return returnVal
	}

	// Inssére o registo atualizado, com a token atualizada, na bd de uma forma assincrona
	go func() {
		// Encodifica os valores de novoUSer para um formato json
		registo, err := json.Marshal(novoUser)
		if err != nil {
			loggers.LoginServerErrorLogger.Println("Erro: ", err)
			panic("Erro ao criar access token, tente de novo.")
		}

		// Inssére o registo do utlizador atualizado na base de dados
		redishandle.SetRegistoBD(&RedisClientDB, redishandle.RegistoRedisDB{
			Key:    novoUser.Username,
			Valor:  registo,
			Expira: 0,
		}, 0)

		loggers.LoginAuthLogger.Println("Token atualizada.")
	}()

	// Loga que o utilisador XXXX iniciou sessão
	loggers.LoginAuthLogger.Println("Utilizador, ", user, ", iniciou sessão")
	returnVal["token"] = novoUserJWTToken
	return returnVal
}

// Registar -
func Registar(user string, password string, token string) map[string]interface{} {
	returnVal := make(map[string]interface{}, 0)

	// Busca o registo do admin
	registoAdmin, err := redishandle.GetRegistoBD(&RedisClientDB, "PyroTM", 0)
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
	jwtTokenSigned, err := CriarNovoUser(user, password).CriarUserJWT().SignedString(hmacSecret)
	if err != nil {
		loggers.LoginAuthLogger.Println("Error: ", err)
		returnVal["error"] = err
		return returnVal
	}
	loggers.LoginAuthLogger.Println("Token: ", jwtTokenSigned[34:54], "...")

	// Cria a struct para o novo user
	novoUser := CriarNovoUser(user, password)
	// Atribui a token gerada ao utilisador
	novoUser.JWT = jwtTokenSigned
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
		// Inssere o novo utilisador na bd se u utilisador não existir
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
