package loginregistohandlers

import "github.com/dgrijalva/jwt-go"

const (
	// ROOT -
	ROOT = iota
	// ADMIN -
	ADMIN
	// USER -
	USER
)

// User -
type User struct {
	Username   string `json:"user,omitempty"`
	Password   string `json:"passwd,omitempty,-"`
	Permissoes int    `json:"perms,omitempty"`
	JWT        string `json:"jwt,omitempty"`
}

// CriarNovoUser -
func CriarNovoUser(user string, password string) User {
	return User{
		Username:   user,
		Password:   password,
		Permissoes: USER,
	}
}

// CriarUserJWT -
func (user User) CriarUserJWT() *jwt.Token {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": user.Username,
		"pass": user.Password,
		"iss":  "loginServer",
	})
	return jwtToken
}
