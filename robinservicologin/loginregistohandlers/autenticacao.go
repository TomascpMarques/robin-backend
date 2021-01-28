package loginregistohandlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"

	"github.com/tomascpmarques/PAP/backend/robinservicologin/loggers"
	"github.com/tomascpmarques/PAP/backend/robinservicologin/redishandle"
)

var hmacSecret = hmac.New(sha256.New, []byte(`SUPPER_SECRET_DEVELOPMENT_KEY`)).Sum(nil)

// RedisClientDB -
var RedisClientDB = redishandle.NovoClienteRedis(redishandle.AddressRed, os.Getenv("AUTH_SERVER_REDIS_PORT"), redishandle.PasswordRed, "", 0)

// Login -
func Login(user string, pass string) string {
	novoUser := CriarNovoUser(user, pass)

	jwtTokenSigned, err := novoUser.CriarUserJWT().SignedString(hmacSecret)
	if err != nil {
		loggers.LoginAuthLogger.Println("Error: ", err)
		return fmt.Sprint(err)
	}
	loggers.LoginAuthLogger.Println("Token: ", jwtTokenSigned[:16], "...")
	//LoginAuthLogger
	_, err = redishandle.GetRegistoBD(&RedisClientDB, user, 0)
	if err != nil {
		loggers.LoginAuthLogger.Println(err)
		return "Credeenciais inválidas"
	}

	loggers.LoginAuthLogger.Println("Utilizador, ", user, ", iniciou sessão")
	return jwtTokenSigned
}

// Registar -
func Registar(user string, password string, token string) string {
	fmt.Println("token: >", token)

	novoUser := CriarNovoUser(user, password)

	jwtToken := CriarNovoUser(user, password).CriarUserJWT()
	jwtTokenSigned, err := jwtToken.SignedString(hmacSecret)
	if err != nil {
		loggers.LoginAuthLogger.Println("Error: ", err)
		return fmt.Sprint(err)
	}

	loggers.LoginAuthLogger.Println("Token: ", jwtTokenSigned[34:54], "...")

	novoUser.JWT = jwtTokenSigned
	loggers.LoginAuthLogger.Println("Novo user: ", user)

	novoUserJSON, err := json.Marshal(&novoUser)
	if err != nil {
		loggers.LoginAuthLogger.Println("Error: ", err)
		return fmt.Sprint(err)
	}

	_, exists := redishandle.GetRegistoBD(&RedisClientDB, user, 0)
	if exists != nil {
		redishandle.SetRegistoBD(&RedisClientDB, redishandle.RegistoRedisDB{
			Key:    novoUser.Username,
			Valor:  novoUserJSON,
			Expira: 0,
		}, 0)
		return jwtTokenSigned
	}

	return ""
}
