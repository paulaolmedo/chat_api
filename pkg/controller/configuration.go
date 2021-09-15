package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/challenge/pkg/service"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var environment bool

// Handler provides the interface to handle different requests
type Handler struct {
	Router   *mux.Router
	Database service.ChatService
}

func (serverConfiguration *Handler) SetEnvironment(envValue bool) {
	environment = envValue
}

func (serverConfiguration *Handler) InitHTTPServer(databasepath string) {
	serverConfiguration.InitDatabase(databasepath)

	serverConfiguration.Router = mux.NewRouter()
	serverConfiguration.InitRouters()

	serverConfiguration.Router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		template, _ := route.GetPathTemplate()
		methods, _ := route.GetMethods()
		fmt.Printf("routes %s %s", methods, template)
		fmt.Println()
		return nil
	})
}

func (serverConfiguration *Handler) Run(host string) {
	fmt.Println("Listening to:", host)
	handler := cors.Default().Handler(serverConfiguration.Router)
	log.Fatal(http.ListenAndServe(host, handler))
}

func (config *Handler) InitRouters() {
	config.Router.HandleFunc("/check", SetMiddlewareWithoutAuthentication(config.Check)).Methods("POST")
	config.Router.HandleFunc("/users", SetMiddlewareWithoutAuthentication(config.CreateUser)).Methods("POST")
	config.Router.HandleFunc("/login", SetMiddlewareWithoutAuthentication(config.Login)).Methods("POST")
	config.Router.HandleFunc("/messages", SetMiddlewareWithAuthentication(config.SendMessage)).Methods("POST")
	config.Router.HandleFunc("/messages", SetMiddlewareWithAuthentication(config.GetMessages)).Methods("GET")
}

func (serverConfiguration *Handler) InitDatabase(databasepath string) {
	dao, err := service.NewDAO(databasepath)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	serverConfiguration.Database = service.NewChatService(dao)
}

func SetMiddlewareWithoutAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

func SetMiddlewareWithAuthentication(next http.HandlerFunc) http.HandlerFunc {
	if environment {
		jwtMiddleware := SetMiddlewareJWT()
		if jwtMiddleware == nil {
			log.Fatal("failed to initialize jwt middleware :(")
		}

		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			jwtMiddleware.HandlerWithNext(w, r, next)
		}
	} else { // no es ambiente productivo
		return SetMiddlewareWithoutAuthentication(next)
	}
}
