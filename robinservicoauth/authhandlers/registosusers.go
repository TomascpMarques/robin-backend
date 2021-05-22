package authhandlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/tomascpmarques/PAP/backend/robinservicoauth/loggers"
	"github.com/tomascpmarques/PAP/backend/robinservicoauth/redishandle"
)

// AtualizarUser atualiza os dados dos utilizador fornecido, depois de verificar a token fornecida
func AtualizarUser(user string, userInfo map[string]interface{}, token string) map[string]interface{} {
	returnVal := make(map[string]interface{})

	if VerificarTokenAdmin(token) != "OK" {
		loggers.LoginAuthLogger.Println("Token inválida.")
		returnVal["err"] = "Token inválida ou expirada"
		return returnVal
	}

	userAtualizar, err := GetUserParaValorStruct(user)
	if err != nil {
		loggers.LoginAuthLogger.Println("Erro: ", "Sem registo para <", user, ">")
		returnVal["erro"] = ("Sem registo para <" + user + ">")
		return returnVal
	}

	// Verifica se os dados passados nos param, são novos ou not null, só depois é que os atualiza
	if userInfo["user"] != nil && userInfo["user"] != userAtualizar.Username {
		userAtualizar.Username = userInfo["user"].(string)
	}
	if userInfo["pass"] != nil && userInfo["pass"] != userAtualizar.Password {
		userAtualizar.Password = userInfo["pass"].(string)
	}
	if userInfo["perms"] != nil && userInfo["perms"] != userAtualizar.Permissoes {
		/* Limita o numero que equival ás permissões na plataforma*/
		if int(userInfo["perms"].(float64)) <= 1 && int(userInfo["perms"].(float64)) >= 3 {
			loggers.LoginAuthLogger.Println("Error: ", "Permissões fora dos valores permitidos, entre 1 e 3")
			returnVal["error"] = "Permissões fora dos valores permitidos, entre 1 e 3"
			return returnVal
		}
		userAtualizar.Permissoes = int(userInfo["perms"].(float64))
	}

	userAtualizadoJSON, err := json.Marshal(&userAtualizar)
	if err != nil {
		loggers.LoginAuthLogger.Println("Error: ", err)
		returnVal["error"] = err
		return returnVal
	}

	// Apaga o registo, para poder ser insserido o novo registo atualizado
	err = redishandle.DelRegistoBD(&RedisClientDB, user)
	if err != nil {
		loggers.LoginAuthLogger.Println("Error: ", err)
		returnVal["error"] = err
		return returnVal
	}

	// Insserção do registo
	redishandle.SetRegistoBD(&RedisClientDB, redishandle.RegistoRedisDB{
		Key:    userAtualizar.Username,
		Valor:  userAtualizadoJSON,
		Expira: 0,
	}, 0)

	returnVal["Menssagem"] = "Sucesso ao alterar dados."
	return returnVal
}

// ApagarUser, apaga um user da bd , pelo id especificado
func ApagarUser(userID string, token string) (retorno map[string]interface{}) {
	retorno = make(map[string]interface{})

	if VerificarTokenAdmin(token) != "OK" {
		loggers.LoginAuthLogger.Println("Token inválida.")
		retorno["err"] = "Token inválida ou expirada"
		return
	}

	// Apaga o registo, para poder ser insserido o novo registo atualizado
	err := redishandle.DelRegistoBD(&RedisClientDB, userID)
	if err != nil {
		loggers.LoginAuthLogger.Println("Error: ", err)
		retorno["error"] = err
		return
	}

	retorno["status"] = "Sucesso!"
	return
}

// SessActualStatus Atualiza a mensagem de status
func SessActualStatus(usrNome string, status string) (retorno map[string]interface{}) {
	retorno = make(map[string]interface{})

	// Mongodb query para atualizar o status do user
	updateQuery := "{\"status\": \"" + status + "\"}"
	// DynamicGoQuery body para conssumir o endpoint do serviço userinfo
	action := fmt.Sprintf("action:\n\"%s\":\n\"%s\",\n%s,", "UpdateInfoUtilizador", usrNome, updateQuery)

	// Utilização do endpoint UpdateInfoUtilizador, exposto em http://0.0.0.0:8001
	resp, err := http.Post("http://0.0.0.0:8001", "text/plain", bytes.NewBufferString(action))
	if err != nil {
		loggers.LoginAuthLogger.Println("Error: ", err)
		retorno["error"] = err
		return
	}
	defer resp.Body.Close()
	bodyContentBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		loggers.LoginAuthLogger.Println("Error: ", err)
		retorno["error"] = "Erro ao ler o conteudo da response do seviço userinfo"
		return
	}

	loggers.LoginResolverLogger.Printf("Update status: %v", string(bodyContentBytes))
	retorno["sucesso"] = "Campo atualizado!"
	return
}

func VerificarUserExiste(userName string, token string) (retorno map[string]interface{}) {
	retorno = make(map[string]interface{})

	if VerificarTokenUser(token) != "OK" {
		loggers.LoginAuthLogger.Println("Token inválida ou expirada.")
		retorno["err"] = "Token inválida ou expirada"
		return
	}

	if !redishandle.BuscarKeysVerificarResultado(context.Background(), &RedisClientDB, userName) {
		loggers.LoginOperacoesBDLogger.Println("Sem registo para a key fornecida, pode ser usada")
		retorno["existe"] = false
		return
	}

	loggers.LoginOperacoesBDLogger.Println("Já existes um registo para a key fornecida")
	retorno["existe"] = true
	return
}
