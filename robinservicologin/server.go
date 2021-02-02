package main

import (
	"context"
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/TomascpMarques/dynamic-querys-go/actions"
	"github.com/gorilla/mux"
	"github.com/tomascpmarques/PAP/backend/robinservicologin/loggers"
	"github.com/tomascpmarques/PAP/backend/robinservicologin/loginregistohandlers"
)

// HTTPport - porta onde o servidor http está localizado
var HTTPport = os.Getenv("LOGIN_SERV_PORT")

// DEFAULTHTTPPORT - valor default para HTTPport
var DEFAULTHTTPPORT = "8080"

func main() {
	if HTTPport == "" {
		HTTPport = DEFAULTHTTPPORT
	}

	// flag setup fo graceful-shutdown
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	// Mapeamento das funções desponíveis
	actions.FuncsStorage["VerificarTokenUser"] = loginregistohandlers.VerificarTokenUser
	actions.FuncsStorage["Login"] = loginregistohandlers.Login
	actions.FuncsStorage["Registar"] = loginregistohandlers.Registar

	router := mux.NewRouter()
	router.HandleFunc("/auth", actions.Handler)
	router.HandleFunc("/test", func(rw http.ResponseWriter, r *http.Request) {
		content, _ := ioutil.ReadAll(r.Body)
		rw.Header().Set("Content-Type", "application/json")
		rw.Write(content)
	})

	srv := &http.Server{
		Addr:         "0.0.0.0:" + HTTPport,
		Handler:      router,
		WriteTimeout: 2 * time.Second,
		ReadTimeout:  2 * time.Second,
		IdleTimeout:  4 * time.Second,
		ErrorLog:     loggers.LoginServerErrorLogger,
	}

	if err := srv.ListenAndServe(); err != nil {
		loggers.LoginServerErrorLogger.Fatal("Erro Fatal: ", err)
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
