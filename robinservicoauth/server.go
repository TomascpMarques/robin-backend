package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/TomascpMarques/dynamic-querys-go/actions"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/tomascpmarques/PAP/backend/robinservicoauth/authhandlers"
	"github.com/tomascpmarques/PAP/backend/robinservicoauth/loggers"
)

// HTTPport - porta onde o servidor http está localizado
var HTTPport = os.Getenv("LOGIN_SERV_PORT")

// DEFAULTHTTPPORT - valor default para HTTPport
var DEFAULTHTTPPORT = "8081"

func main() {
	if HTTPport == "" {
		HTTPport = DEFAULTHTTPPORT
	}

	// flag setup fo graceful-shutdown
	var flagWait time.Duration
	flag.DurationVar(&flagWait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	// Mapeamento das funções desponíveis aos action requests
	actions.FuncsStorage["VerificarTokensParaReAuth"] = authhandlers.VerificarTokenReAuth
	actions.FuncsStorage["VerificarUserExiste"] = authhandlers.VerificarUserExiste
	actions.FuncsStorage["VerificarTokenAdmin"] = authhandlers.VerificarTokenAdmin
	actions.FuncsStorage["VerificarTokenUser"] = authhandlers.VerificarTokenUser
	actions.FuncsStorage["SessActualStatus"] = authhandlers.SessActualStatus
	actions.FuncsStorage["AtualizarUser"] = authhandlers.AtualizarUser
	actions.FuncsStorage["ApagarUser"] = authhandlers.ApagarUser
	actions.FuncsStorage["Registar"] = authhandlers.Registar
	actions.FuncsStorage["Login"] = authhandlers.Login

	router := mux.NewRouter()
	go router.HandleFunc("/", actions.Handler)

	// Defenições de partilha de recursos cruzada
	corsOptions := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080"}, // Permite só o localhost:8080 a fazer requests ao serviço
		AllowCredentials: false,
	})

	//handler := cors.Default().Handler(router)
	corsHandlers := corsOptions.Handler(router)

	// Defenições do servidor web
	srv := &http.Server{
		Handler:      corsHandlers,                   // Gestor dos requests ao entrar no servidor
		Addr:         "0.0.0.0:" + HTTPport,          // Localização do web server, ip + port combo
		WriteTimeout: 2 * time.Second,                // Se o pedido demorar mais do que 2s a escrever o conteúdo fecha a conexão
		ReadTimeout:  2 * time.Second,                // Se o servidor demorar mais que 2s a ler o request, fecha a conexão
		IdleTimeout:  4 * time.Second,                // Quando keep-alive estiver especificado, se a próxima conec, demorar mais de 2s fecha
		ErrorLog:     loggers.LoginServerErrorLogger, // Logger dos erros de servidor
	}

	// Exits se não se consseguir iniciar o servidor com as defnições necessárias
	if err := srv.ListenAndServe(); err != nil {
		loggers.LoginServerErrorLogger.Fatal("Erro Fatal: ", err)
	}

	// Graceful-Shutdown implementation
	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C) or SIGKILL,
	// SIGQUIT will not be caught.
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), flagWait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)

	actions.DQGLogger.Println("Servidor a desligar")
	os.Exit(0)
}
