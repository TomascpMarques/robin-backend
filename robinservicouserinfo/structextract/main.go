package structextract

import (
	"reflect"
	"strings"
)

// CustomExtractSchema -
type CustomExtractSchema map[string][]string

/*
	Exemplo do query passado à func ExtrairCampos...
	var custom = CustomExtractSchema{
		"PC":         {"PC", "Nome,ID"},
		"Info":       {"PC", "Sala"},
		"Manutencao": {"Info", "Status"},
		"Ultima":     {"Manutencao", "ID"},
		"AAAA":       {"Ultima", "ID"},
	}
*/

/*
	Exemplo do query escrito em go actions
	{"Item": ["Item", "Quantidade,Nome"]},
*/

/*
ExtrairCamposEspecificosStruct :
	Extrai os campos especificados da struct fornecida,
	utiliza um map[string][]string para especificar os campos a tirar da struct.

Params:
	-> estrutura interface{}, struct src para extrair os dados
	-> listaCampos CustomExtractSchema (map[string][]string), especifica os acmpos a extrair da struct

Notas:
	Utiliza-se uma interface{} como src dos valorespara se poder utilisar qualquer struct passada como param.
	A função é recurssiva para chegar a todas as structs presentes, até ás mais profundas. Os nomes dos campos da struct
		devem ser iguais aos defenidos nos campos das structs
*/
func ExtrairCamposEspecificosStruct(estrutura interface{}, listaCampos map[string][]string) (retorno map[string]interface{}) {
	retorno = make(map[string]interface{})

	// Valore refletido da estrutura passada
	estruturaReflect := reflect.ValueOf(estrutura)
	// Tipo refletido da estrutura passada
	estruturaReflectType := reflect.TypeOf(estrutura)

	// Iterar por todos os campos (à superficíe) da struct
	for i := 0; i < estruturaReflect.NumField(); i++ {
		// Verifica se o campo foi especificado na lista com os valores a extrair
		if _, existe := listaCampos[estruturaReflectType.Name()]; existe {
			// Verifica se o campo atual da estrutura foi indicado para extrasão
			if strings.Contains(listaCampos[estruturaReflect.Type().Name()][1], estruturaReflectType.Field(i).Name) {
				// Adiciona o valor extraido ao retorno da função
				retorno[estruturaReflectType.Field(i).Name] = estruturaReflect.Field(i).Interface()
			}

			// Verifica se o campo atual é mais uma estrutura, (parte recurssiva da função)
			if estruturaReflectType.Field(i).Type.Kind().String() == "struct" {
				// Reflexão dos valores presentes na estrutura embutida
				estruturaEmbutida := estruturaReflect.Field(i)
				nomeEstruturaEmbutida := estruturaReflectType.Field(i)

				// Indica o index onde o parent da struct começa na string ex: main.PC, ignora tudo até ao 1º "." +1
				dropPacoteIndex := strings.Index(estruturaReflect.Type().String(), ".") + 1
				structParent := estruturaReflect.Type().String()[dropPacoteIndex:]

				// Se a estrutura não estiver vazia, não null, e for um campo defenido na lista através dos parents de cada campo
				// Chama esta função de novo para ler e extrair os valores da struct
				if !estruturaEmbutida.IsZero() && structParent == listaCampos[nomeEstruturaEmbutida.Name][0] {
					retorno[estruturaEmbutida.Type().Name()] = ExtrairCamposEspecificosStruct(estruturaEmbutida.Interface(), listaCampos)
				}
			}
		}
	}

	return
}
