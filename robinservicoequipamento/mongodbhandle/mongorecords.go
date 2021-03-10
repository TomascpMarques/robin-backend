package mongodbhandle

import "github.com/tomascpmarques/PAP/backend/robinservicoequipamento/resolvedschema"

// ParseTipoDeRegisto -
func ParseTipoDeRegisto(alvo map[string]interface{}) interface{} {

	if alvo["tipo_de_registo"] == "Item" {
		return resolvedschema.ItemParaStruct(&alvo)
	}

	return nil
}
