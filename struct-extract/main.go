package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// PC -
type PC struct {
	ID         string `json:"id,omitempty"`
	Nome       string `json:"nome,omitempty"`
	Info       Info   `json:"info,omitempty"`
	Manutencao Manutencao
}

// Info -
type Info struct {
	Sala       int        `json:"sala,omitempty"`
	Piso       int        `json:"piso,omitempty"`
	Notas      string     `json:"notas,omitempty"`
	Manutencao Manutencao `json:"manutencao,omitempty"`
}

// Manutencao -
type Manutencao struct {
	Status           string `json:"status,omitempty"`
	UltimaManutencao string `json:"ultimamanutencao,omitempty"`
	Ultima           Ultima
}

// Ultima -
type Ultima struct {
	ID   string `json:"id,omitempty"`
	Nome string `json:"nome,omitempty"`
	AAAA AAAA
}

// AAAA -
type AAAA struct {
	ID   string `json:"id,omitempty"`
	Nome string `json:"nome,omitempty"`
}

// CustomExtractSchema -
type CustomExtractSchema map[string][]string

func main() {
	// Example -
	var Example = PC{
		ID:   "PC1",
		Nome: "PC Super Fixe",
		Info: Info{
			Sala:  2,
			Piso:  3,
			Notas: "OMG XD",
			Manutencao: Manutencao{
				Status:           "Em circulação",
				UltimaManutencao: "24/12/19",
				Ultima: Ultima{
					ID:   "1234",
					Nome: "4321",
					AAAA: AAAA{
						ID: "456",
					},
				},
			},
		},
		Manutencao: Manutencao{
			Status: "alive",
		},
	}

	var custom = CustomExtractSchema{
		"PC":         {"PC", "Nome,ID"},
		"Info":       {"PC", "Sala"},
		"Manutencao": {"Info", "Status"},
		"Ultima":     {"Manutencao", "ID"},
		"AAAA":       {"Ultima", "ID"},
	}

	aaaa := ExtrairCamposEspecificosStruct(Example, custom)
	fmt.Println("Res: ", aaaa)

	x, err := json.MarshalIndent(aaaa, "", "   ")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Res Json: ", x)
}

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
func ExtrairCamposEspecificosStruct(estrutura interface{}, listaCampos CustomExtractSchema) (retorno map[string]interface{}) {
	retorno = make(map[string]interface{}, 0)

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
