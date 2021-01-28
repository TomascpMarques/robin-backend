package loginregistohandlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/tomascpmarques/PAP/backend/robinservicologin/loggers"
	"github.com/tomascpmarques/PAP/backend/robinservicologin/redishandle"
)

var hmacSecret = hmac.New(sha256.New, []byte(`SUPPER_SECRET_DEVELOPMENT_KEY`)).Sum(nil)

// RedisClientDB -
var RedisClientDB = redishandle.NovoClienteRedis(redishandle.AddressRed, os.Getenv("AUTH_SERVER_REDIS_PORT"), redishandle.PasswordRed, "", 0)

// TestLoggedUser -
func TestLoggedUser(userToken string) string {
	token, err := jwt.Parse(userToken, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return hmacSecret, nil
	})
	if err != nil {
		return fmt.Sprint(err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims)
		fmt.Println(claims["exp"] == time.Now().Unix())
		fmt.Println(claims["exp"], time.Now().Unix())

		return "OK"
	}
	fmt.Println("err: ", err)
	return "not ok ?"
}

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
	registo, err := redishandle.GetRegistoBD(&RedisClientDB, "PyroTM", 0)
	if err != nil {
		loggers.LoginAuthLogger.Println("Erro ao buscar jwt do utilisador necessário.")
		return ""
	}
	var userReg User
	err = json.Unmarshal([]byte(registo), &userReg)
	if err != nil {
		loggers.LoginAuthLogger.Println("Erro ao descodificar registo do utilizador.")
		return ""
	}
	if userReg.JWT != token {
		loggers.LoginAuthLogger.Println("A token não têm permissões")
		return `token de autorização errado`
	}

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
