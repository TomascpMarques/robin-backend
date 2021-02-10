package endpointfuncs

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/tomascpmarques/PAP/backend/robinservicoequipamento/redishandle"
	"github.com/tomascpmarques/PAP/backend/robinservicoequipamento/resolvedschema"
)

//RedisClientDB -
var RedisClientDB = redishandle.NovoClienteRedis(
	os.Getenv("REDISADDRESS"),
	"8080",
	"",
	"",
	0,
)
var _ = VerificarAdminFirstBoot()

// Hello - HELLOs back
func Hello(name map[string]interface{}) map[string]interface{} {
	return name
}

// AdicionarRegistoDeItem -
func AdicionarRegistoDeItem(token string, tipoDeItem string, item []byte) map[string]interface{} {
	result := make(map[string]interface{}, 0)

	if VerificarTokenUser(token) != "OK" {
		fmt.Println("Erro: A token fornecida é inválida ou expirou")
		result["erro"] = "A token fornecida é inválida ou expirou"
		return result
	}

	var test resolvedschema.Item
	bts := json.Unmarshal(item, &test)
	fmt.Println(bts, test)

	return result
}

// ApagarRegistoDeItem -
func ApagarRegistoDeItem() {}

// AtualizararRegistoDeItem -
func AtualizararRegistoDeItem() {}

// BuscarInfoDeItems -
func BuscarInfoDeItems() {}

// BuscarInfoDeItem -
func BuscarInfoDeItem() {}
