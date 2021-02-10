package endpointfuncs

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/tomascpmarques/PAP/backend/robinservicoequipamento/loggers"
	"github.com/tomascpmarques/PAP/backend/robinservicoequipamento/redishandle"
)

var hmacSecret = hmac.New(sha256.New, []byte(`SUPPER_SECRET_DEVELOPMENT_KEY`)).Sum(nil)

// User define um utilizador
type User struct {
	Username   string `json:"user,omitempty"`
	Password   string `json:"passwd,omitempty,-"`
	
	Permissoes int    `json:"perms,omitempty"`
	JWT        string `json:"jwt,omitempty"`
}

// VerificarAdminFirstBoot verifica se o utilizador admin da backend robin existe, se não existir cria esse user
// com as credenciais default
func VerificarAdminFirstBoot() bool {
	// Tenta encontrar o registo do admin, se não o encontrar cria-o
	_, err := redishandle.GetRegistoBD(&RedisClientDB, "admin", 0)
	if err != nil {
		loggers.DbFuncsLogger.Println("O utilizador administrador não existe...")
		// Cria a struct de utilisador para o admin
		admin := User{
			Username:   "admin",
			Password:   "532f1f7e5e4ae1475835c4978696c1e3",
			Permissoes: 2,
		}
		registoUserJSON, err := json.Marshal(&admin)
		if err != nil {
			loggers.DbFuncsLogger.Println("Erro: ", err)
			return false
		}
		// Inssere o administrador
		redishandle.SetRegistoBD(&RedisClientDB, redishandle.RegistoRedisDB{
			Key:    admin.Username,
			Valor:  registoUserJSON,
			Expira: 0,
		}, 0)

		return true
	}
	return false
}

// VerificarTokenUser verifica se a token passada é válida, logo vê se já expirou
// se o modo de assinatura é o correto, e se o emissor é o servidor de autenticação
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

// VerificarTokenAdmin verifica tudo o que a função VerificarTokenUser verifica,
// e ainda verifica se o utilisador é o administrador
func VerificarTokenAdmin(userToken string) string {
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
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid && claims["iss"] == "Robin-Servico-Auth" &&
		claims["perms"].(float64) == 2 {
		return "OK"
	}
	return "Token inválida ou expirada"
}
