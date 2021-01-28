package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/TomascpMarques/dynamic-querys-go/actions"
	"github.com/gorilla/mux"
	"github.com/tomascpmarques/PAP/backend/robinservicoequipamento/endpointfuncs"
	"github.com/tomascpmarques/PAP/backend/robinservicoequipamento/redishandle"
)

var redisClientDB = redishandle.NovoClienteRedis(redishandle.AddressRed, "8080", redishandle.PasswordRed, "", 0)

func main() {
	os.Setenv("ENV_GOACTIONS_PORT", "8000")
	actions.FuncsStorage["Hello"] = endpointfuncs.Hello

	// flag setup fo graceful-shutdown
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	r := mux.NewRouter()
	r.HandleFunc("/", actions.Handler)

	server := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:" + os.Getenv("ENV_GOACTIONS_PORT"),
	}

	if err := server.ListenAndServe(); err != nil {
		log.Println("Erro: ", err)
		os.Exit(1)
	}

	// Graceful-Shutdown implementation
	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C) or SIGKILL,
	// SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt, os.Kill)

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
