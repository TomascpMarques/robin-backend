package schemaresolvehelper

import (
	"reflect"
)

/*
OperacoesMemoriaStructsValores -
---
Params
	-
*/
func OperacoesMemoriaStructsValores(localStructInput, localStructModificar interface{}) {
	inputReflect := reflect.ValueOf(localStructInput).Elem()
	structAlvo := reflect.ValueOf(localStructModificar).Elem()

	for i := 0; i < structAlvo.NumField(); i++ {
		if inputReflect.Field(i).IsZero() == false {
			structCampoString := inputReflect.Type().Field(i).Name

			tipoVariavel := structAlvo.FieldByName(structCampoString).Addr().Type().String()
			enderecoValor := structAlvo.FieldByName(structCampoString).Addr().Interface()

			if tipoVariavel == "*string" {
				//log.Println("string: ", tipoVariavel, inputReflect.Type().Field(i).Name, inputReflect.Field(i).String())
				enderecoValorCampo := enderecoValor.(*string)
				varMudaValor := enderecoValorCampo
				*varMudaValor = inputReflect.Field(i).String()
				continue
			}
			if tipoVariavel == "**string" {
				//log.Println("**string: ", tipoVariavel, inputReflect.Type().Field(i).Name, inputReflect.Field(i).String())
				enderecoValorCampo := structAlvo.FieldByName(structCampoString).Addr().Interface().(**string)
				varMudaValor := enderecoValorCampo
				**varMudaValor = inputReflect.Field(i).Elem().String()
				continue
			}

			if tipoVariavel == "*int" {
				//log.Println("*int: ", tipoVariavel, inputReflect.Type().Field(i).Name, inputReflect.Field(i).String())
				enderecoValorCampo := structAlvo.FieldByName(structCampoString).Addr().Interface().(*int)
				varMudaValor := enderecoValorCampo
				*varMudaValor = int(inputReflect.Field(i).Elem().Int())
				continue
			}
			if tipoVariavel == "**int" {
				//log.Println("**int: ", tipoVariavel, inputReflect.Type().Field(i).Name, inputReflect.Field(i).String())
				enderecoValorCampo := structAlvo.FieldByName(structCampoString).Addr().Interface().(**int)
				varMudaValor := enderecoValorCampo
				**varMudaValor = int(inputReflect.Field(i).Elem().Int())
				continue
			}

			if tipoVariavel == "*float64" {
				//log.Println("*int: ", tipoVariavel, inputReflect.Type().Field(i).Name, inputReflect.Field(i).String())
				enderecoValorCampo := structAlvo.FieldByName(structCampoString).Addr().Interface().(*float64)
				varMudaValor := enderecoValorCampo
				*varMudaValor = inputReflect.Field(i).Elem().Float()
				continue
			}
			if tipoVariavel == "**float64" {
				//log.Println("**int: ", tipoVariavel, inputReflect.Type().Field(i).Name, inputReflect.Field(i).String())
				enderecoValorCampo := structAlvo.FieldByName(structCampoString).Addr().Interface().(**float64)
				varMudaValor := enderecoValorCampo
				**varMudaValor = inputReflect.Field(i).Elem().Float()
				continue
			}

			if tipoVariavel == "*float64" {
				//log.Println("*int: ", tipoVariavel, inputReflect.Type().Field(i).Name, inputReflect.Field(i).String())
				enderecoValorCampo := structAlvo.FieldByName(structCampoString).Addr().Interface().(*float64)
				varMudaValor := enderecoValorCampo
				*varMudaValor = inputReflect.Field(i).Elem().Float()
				continue
			}
			if tipoVariavel == "**float64" {
				//log.Println("**int: ", tipoVariavel, inputReflect.Type().Field(i).Name, inputReflect.Field(i).String())
				enderecoValorCampo := structAlvo.FieldByName(structCampoString).Addr().Interface().(**float64)
				varMudaValor := enderecoValorCampo
				**varMudaValor = inputReflect.Field(i).Elem().Float()
				continue
			}
		}
	}
}
