package loginregistohandlers

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	// ROOT -
	ROOT = iota
	// ADMIN -
	ADMIN
	// USER -
	USER
)

// User - Epecifica os dados que definem um utilizador
type User struct {
	Username   string `json:"user,omitempty"`
	Password   string `json:"passwd,omitempty,-"`
	Permissoes int    `json:"perms,omitempty"`
	Logged     bool   `json:"logged"`
	JWT        string `json:"jwt,omitempty"`
}

// CriarNovoUser -
func CriarNovoUser(user string, password string) User {
	return User{
		Username:   user,
		Password:   password,
		Logged:     false,
		Permissoes: USER,
	}
}

// CriarUserJWT -
func (user User) CriarUserJWT() *jwt.Token {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": user.Username,
		"pass": user.Password,
		"iss":  "loginServer",
		"exp":  time.Now().Add(time.Second * 240).Unix(),
	})
	return jwtToken
}
