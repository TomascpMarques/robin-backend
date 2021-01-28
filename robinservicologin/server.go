package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/TomascpMarques/dynamic-querys-go/actions"
	"github.com/gorilla/mux"
	"github.com/tomascpmarques/PAP/backend/robinservicologin/loggers"
	"github.com/tomascpmarques/PAP/backend/robinservicologin/loginregistohandlers"
)

// HTTPport -
var HTTPport = os.Getenv("LOGIN_SERV_PORT")

// DEFAULTHTTPPORT -
var DEFAULTHTTPPORT = "5600"

func main() {
	if HTTPport == "" {
		HTTPport = DEFAULTHTTPPORT
	}

	// flag setup fo graceful-shutdown
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	actions.FuncsStorage["Login"] = loginregistohandlers.Login
	actions.FuncsStorage["Registar"] = loginregistohandlers.Registar

	router := mux.NewRouter()
	router.HandleFunc("/auth", actions.Handler)

	srv := &http.Server{
		Addr:         "127.0.0.1:" + HTTPport,
		Handler:      router,
		WriteTimeout: 2 * time.Second,
		ReadTimeout:  2 * time.Second,
		IdleTimeout:  4 * time.Second,
		ErrorLog:     loggers.LoginServerErrorLogger,
	}

	if err := srv.ListenAndServe(); err != nil {
		srv.ErrorLog.Fatal("Erro: ", err)
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
	srv.Shutdown(ctx)

	actions.DQGLogger.Println("Servidor a desligar")
	os.Exit(0)
}
