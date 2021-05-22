package endpointfuncs

import (
	"fmt"
	"reflect"

	"github.com/tomascpmarques/PAP/backend/robinservicovideoshare/loggers"
	"github.com/tomascpmarques/PAP/backend/robinservicovideoshare/mongodbhandle"
	"github.com/tomascpmarques/PAP/backend/robinservicovideoshare/resolvedschema"
)

var mongoParams = mongodbhandle.MongoConexaoParams{
	URI: "mongodb://0.0.0.0:27022/",
}

// MongoClient cliente com a conexão à instancia mongo
var MongoClient = mongodbhandle.CriarConexaoMongoDB(mongoParams)

// PingServico responde que o serviço está online
func PingServico(name string) (retorno map[string]interface{}) {
	retorno = make(map[string]interface{})
	retorno["status"] = fmt.Sprintf("Hello %s, I'm alive and OK", name)
	return
}

func CriarVideoShare(videoMetaData map[string]interface{}, token string) (retorno map[string]interface{}) {
	retorno = make(map[string]interface{})

	// Verifica se o request está a ser efetuado por um user logado
	// if VerificarTokenUser(token) != "OK" {
	// 	loggers.ResolverLogger.Println("Token inválida.")
	// 	retorno["err"] = "Token inválida ou expirada"
	// 	return
	// }

	// Verifica se o criador do video é o mesmo que o que fez o request
	// Assim evita outros users criarem videos à passarem-se por outros users
	// if VerificarTokenUserSpecif(token, videoMetaData["criador"].(string)) != "OK" {
	// 	loggers.ResolverLogger.Println("O criador deste vídeo não é o autor do request.")
	// 	retorno["err"] = "O criador deste vídeo não é o autor do request"
	// 	return
	// }

	// Verifica os dados fornecidos
	err := VerificarVideoShareMetaData(videoMetaData)
	if err != nil {
		loggers.ResolverLogger.Println("A info fornecida para a videoshare não é válida")
		retorno["err"] = err.Error()
		return
	}

	// Transforma a info da videoshare numa struct
	videoSahre := resolvedschema.VideoParaStruct(&videoMetaData)
	if reflect.ValueOf(videoSahre).IsZero() {
		loggers.ResolverLogger.Println("Erro ao converter as informações dadas para uma struct.")
		retorno["err"] = "Erro ao converter as informações dadas para uma struct"
		return
	}

	// Retira a parte desnecessária do URL
	videoSahre.URL = TrimURL(videoSahre.URL)

	// Adicona a videoshare à bd
	err = AdicionarVideoShareDB(&videoSahre)
	if err != nil {
		loggers.ResolverLogger.Println("Erro ao criar registo da videoshare na BD")
		retorno["err"] = err.Error()
		return
	}

	loggers.ResolverLogger.Println("Criado registo da videoshare na BD")
	retorno["sucesso"] = true
	return
}
