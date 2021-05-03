package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/TomascpMarques/dynamic-querys-go/actions"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/tomascpmarques/PAP/backend/robinservicoequipamento/endpointfuncs"
	"github.com/tomascpmarques/PAP/backend/robinservicoequipamento/loggers"
)

func main() {
	os.Setenv("ENV_ROBINEQUIPAMENTO_PORT", "8000")

	// flag setup fo graceful-shutdown
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "gracefully wait for existing connections to finish in 15s")
	flag.Parse()

	actions.FuncsStorage["BuscarRegistoCustom"] = endpointfuncs.QueryRegistosCustom
	actions.FuncsStorage["ApagarRegistoDeItem"] = endpointfuncs.ApagarRegistoPorID
	actions.FuncsStorage["AtualizarRegistoID"] = endpointfuncs.AtualizararRegistoDeItem
	actions.FuncsStorage["AdicionarRegisto"] = endpointfuncs.AdicionarRegisto
	actions.FuncsStorage["BuscarRegistoID"] = endpointfuncs.QueryRegistoObjID
	actions.FuncsStorage["BuscarInfoItems"] = endpointfuncs.BuscarInfoItems
	actions.FuncsStorage["BuscarInfoItem"] = endpointfuncs.BuscarInfoItemQuery
	actions.FuncsStorage["PingServico"] = endpointfuncs.PingServico

	router := mux.NewRouter()
	router.HandleFunc("/", actions.Handler)

	corsOptions := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080"},
		AllowCredentials: false,
	})

	//handler := cors.Default().Handler(router)
	corsHandlers := corsOptions.Handler(router)

	server := &http.Server{
		Handler:      corsHandlers,
		Addr:         "0.0.0.0:" + os.Getenv("ENV_ROBINEQUIPAMENTO_PORT"),
		IdleTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 3,
		ReadTimeout:  time.Second * 2,
		ErrorLog:     loggers.ServerErrorLogger,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Println("Erro: ", err)
		os.Exit(1)
	}

	// Graceful-Shutdown implementation
	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C) or SIGKILL,
	// SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	server.Shutdown(ctx)

	actions.DQGLogger.Println("Servidor a desligar")
	os.Exit(0)
}
