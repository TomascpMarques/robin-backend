package main

import (
	"fmt"
	"go-graphql-equipamento/graph"
	"go-graphql-equipamento/graph/generated"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultPort = "8080"

var serverloger = log.New(os.Stdout, "ServerGate.......{*} ", log.LstdFlags)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("---------- START ----------")
	serverloger.Println("O servidor está a iniciar ... ")
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	serverloger.Printf("Connect to http://localhost:%s/ for GraphQL playground", port)
	serverloger.Fatal(http.ListenAndServe(":"+port, nil))

	<-exit
	serverloger.Println("O servidor parou. A tentar Gracefull Shutdown")
	err := graph.RedisClienteDB.Close()
	if err != nil {
		serverloger.Println("ERRO de saida: ", err)
	}
	defer serverloger.Println("O Servidor Saiu...\n ---------- END ----------")
}
