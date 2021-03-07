package endpointfuncs

import (
	"context"

	"github.com/tomascpmarques/PAP/backend/robinservicoequipamento/mongodbhandle")


var mongoParams = mongodbhandle.MongoConexaoParams {
	Ctx: context.Background(),
	Cancel:
}
	

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

	return result
}

// AtualizararRegistoDeItem -
func AtualizararRegistoDeItem() {}

// BuscarInfoDeItems -
func BuscarInfoDeItems() {}

// BuscarInfoDeItem -
func BuscarInfoDeItem() {}
