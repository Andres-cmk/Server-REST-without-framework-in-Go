// @title REST Task Server
// @version 1.0
// @description API REST para gestión de tareas
// @contact.name Andres
// @host localhost:8443
// @BasePath /
// @securityDefinitions.basic BasicAuth

package main

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	s "github.com/swaggo/http-swagger"
	"log"
	"net/http"
	d "restServer/docs"
	"restServer/graph"
	"restServer/internal"
	"restServer/server"
)

func main() {

	// Servidor principal
	mux := http.NewServeMux()

	// Logica de negocio
	taskServer := server.NewTaskServer()

	// Verificar que el store existe
	store := taskServer.GetStore()
	log.Printf("TaskServer store creado: %p", store)

	// GraphQL server - comparte el mismo store que REST
	graphqlServer := handler.New(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{Store: store},
	}))

	// Configurar transportes HTTP para GraphQL
	graphqlServer.AddTransport(transport.Options{})
	graphqlServer.AddTransport(transport.GET{})
	graphqlServer.AddTransport(transport.POST{})

	log.Printf("GraphQL Resolver store configurado: %p", store)

	// Swagger config
	d.SwaggerInfo.Schemes = []string{"https"}
	d.SwaggerInfo.Host = "localhost:8443"

	// Endpoints publicos (sin autenticación)
	mux.Handle("/docs/", s.WrapHandler)

	// GraphQL endpoints
	mux.Handle("/graphql", graphqlServer)
	mux.Handle("/playground", playground.Handler("GraphQL Playground", "https://localhost:8443/graphql"))
	mux.HandleFunc("POST /task/", taskServer.CreateTaskHandler)
	mux.HandleFunc("GET /task/{id}/", taskServer.GetTaskHandler)
	mux.HandleFunc("GET /tag/{tag}/", taskServer.TagHandler)
	mux.HandleFunc("GET /due/{year}/{month}/{day}/", taskServer.DueHandler)
	mux.HandleFunc("GET /task/", taskServer.GetAllTasksHandler)
	mux.HandleFunc("DELETE /task/", taskServer.DeleteAllTasksHandler)
	mux.HandleFunc("DELETE /task/{id}/", taskServer.DeleteTaskHandler)

	// Endpoints privados (con autenticación básica)
	/**
	mux.Handle("GET /task/", internal.BasicAuth("admin", "1234",
		http.HandlerFunc(taskServer.GetAllTasksHandler)))
	mux.Handle("DELETE /task/", internal.BasicAuth("admin", "1234",
		http.HandlerFunc(taskServer.DeleteAllTasksHandler)))
	mux.Handle("DELETE /task/{id}/", internal.BasicAuth("admin", "1234",
		http.HandlerFunc(taskServer.DeleteTaskHandler)))
	*/

	// Middlewares globales
	h := internal.Logging(mux)
	handlerResponseServer := internal.NameResponseServer(h, "Andres :D")

	log.Printf("Listening in https://localhost:8443\n")
	log.Fatal(http.ListenAndServeTLS(
		":8443",
		"localhost.pem", // se generan con mkcert -install localhost. Para generarlos se tiene que ir a https://github.com/FiloSottile/mkcert
		"localhost-key.pem",
		handlerResponseServer))
}

// https://github.com/swaggo/swag
