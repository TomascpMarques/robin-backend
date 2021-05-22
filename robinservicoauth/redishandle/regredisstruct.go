package redishandle

import (
	"time"
)

// RegistoRedisDB Estrutura o registo a insserir ou retirar da BD
type RegistoRedisDB struct {
	Key    string
	Valor  interface{}   // Conteúdo do registo
	Expira time.Duration // Expiração do registo
}
