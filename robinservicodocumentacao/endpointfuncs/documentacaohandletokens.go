package endpointfuncs

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

var hmacSecret = hmac.New(sha256.New, []byte(`SUPPER_SECRET_DEVELOPMENT_KEY`)).Sum(nil)

// User define um utilizador
type User struct {
	Username   string `json:"user,omitempty"`
	Password   string `json:"passwd,omitempty"`
	Permissoes int    `json:"perms,omitempty"`
	JWT        string `json:"jwt,omitempty"`
}

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

// VerificarTokenUserSpecif verifica se a token passada é válida, logo vê se já expirou
// se o modo de assinatura é o correto, e se o emissor é o servidor de autenticação, e se o user especificado é igual ao da token
func VerificarTokenUserSpecif(userToken string, user string) string {
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
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid && claims["iss"] == "Robin-Servico-Auth" && claims["user"] == user {
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

// DevolveTokenClaims Devolve os valores do body da token (claims)
func DevolveTokenClaims(userToken string) map[string]interface{} {
	token, err := jwt.Parse(userToken, func(token *jwt.Token) (interface{}, error) {
		// valida o metodo de assinatura da key
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, nil
		}
		// hmacSampleSecret é o []byte que contem o segredo de assinatura
		return hmacSecret, nil
	})
	// Se a token for assinada por outro metodo ou a key for diferente dá erro
	if err != nil {
		return nil
	}

	return token.Claims.(jwt.MapClaims)
}

// VerificarTokenReAuth Verifica a token de reload de autenticação do user
func VerificarTokenReAuth(reAuthToken string) string {
	token, err := jwt.Parse(reAuthToken, func(token *jwt.Token) (interface{}, error) {
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
		claims["typ"] == "reauth" {
		return "OK"
	}
	return "Token inválida ou expirada"
}
