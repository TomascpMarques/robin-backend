package loginregistohandlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/tomascpmarques/PAP/backend/robinservicologin/loggers"
	"github.com/tomascpmarques/PAP/backend/robinservicologin/redishandle"
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
		returnVal["erro"] = err
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
		userAtualizar.Permissoes = userInfo["perms"].(int)
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
func ApagarUser(userID, token string) (retorno map[string]interface{}) {
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
func SessActualStatus(usrNome string, status string) (results map[string]interface{}) {
	results = make(map[string]interface{})

	updateQuery := "{\"$set\":{\"status\": \"" + status + "\"}}"
	action := fmt.Sprintf("action:\n\"%s\":\n\"%s\",\n%s,", "UpdateInfoUtilizador", usrNome, updateQuery)

	resp, err := http.Post("http://0.0.0.0:8001", "text/plain", bytes.NewBufferString(action))
	if err != nil {
		loggers.LoginAuthLogger.Println("Error: ", err)
		results["error"] = err
		return
	}
	defer resp.Body.Close()
	bodyContentBytes, _ := ioutil.ReadAll(resp.Body)

	loggers.LoginResolverLogger.Printf("Update status: %v", string(bodyContentBytes))

	results["sucesso"] = "Campo atualizado!"
	return
}
