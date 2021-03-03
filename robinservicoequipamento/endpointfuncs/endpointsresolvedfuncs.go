package endpointfuncs

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/tomascpmarques/PAP/backend/robinservicoequipamento/redishandle"
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

// AdicionarRegisto -
func AdicionarRegisto(tipoItem string, item []byte, token string) map[string]interface{} {
	result := make(map[string]interface{}, 0)

	// if VerificarTokenUser(token) != "OK" {
	// 	fmt.Println("Erro: A token fornecida é inválida ou expirou")
	// 	result["erro"] = "A token fornecida é inválida ou expirou"
	// 	return result
	// }

	var surfaceLevelKeys map[string]interface{}

	err := json.Unmarshal(item, &surfaceLevelKeys)
	if err != nil {
		fmt.Println("Erro: Erro ao descodificar valor json")
		result["erro"] = "Erro ao descodificar valor json"
		return result
	}

	var registo redishandle.RegistoRedisDB
	registo.CriaEstruturaRegisto(&RedisClientDB, item, tipoItem)
	redishandle.SetRegistoBD(&RedisClientDB, registo, 0)

	result["keys"] = surfaceLevelKeys
	return result
}

// ApagarRegistoDeItem -
func ApagarRegistoDeItem(idItem string, token string) map[string]interface{} {
	result := make(map[string]interface{}, 0)

	// if VerificarTokenUser(token) != "OK" {
	// 	fmt.Println("Erro: A token fornecida é inválida ou expirou")
	// 	result["erro"] = "A token fornecida é inválida ou expirou"
	// 	return result
	// }

	err := redishandle.DelRegistoBD(&RedisClientDB, idItem)
	if err != nil {
		result["id_registo"] = idItem
		result["apagado"] = false
		result["erro"] = err
		return result
	}

	result["id_registo"] = idItem
	result["apagado"] = true
	return result
}

// AtualizararRegistoDeItem -
func AtualizararRegistoDeItem() {}

// BuscarInfoDeItems -
func BuscarInfoDeItems() {}

// BuscarInfoDeItem -
func BuscarInfoDeItem() {}
