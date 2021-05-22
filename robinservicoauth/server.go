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

	corsOptions := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080"},
		AllowCredentials: true,
	})

	//handler := cors.Default().Handler(router)
	corsHandlers := corsOptions.Handler(router)

	srv := &http.Server{
		Handler:      corsHandlers,
		Addr:         "0.0.0.0:" + HTTPport,
		WriteTimeout: 2 * time.Second,
		ReadTimeout:  2 * time.Second,
		IdleTimeout:  4 * time.Second,
		ErrorLog:     loggers.LoginServerErrorLogger,
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
