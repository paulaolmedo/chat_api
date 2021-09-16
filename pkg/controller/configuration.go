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

// SetEnvironment perhaps it's not optimal, but sets the enviroment to determine whether the tokens will be necessary or not
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
	config.Router.HandleFunc(CheckEndpoint, SetMiddlewareWithoutAuthentication(config.Check)).Methods("POST")

	// swagger:operation POST /users CreateUser
	//
	// Add a new user
	//
	// Create a user in the system.
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - in: body
	//   name: user
	//   description: User Information
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/User"
	// responses:
	//   '201':
	//     description: User Information
	//     schema:
	//       type: array
	//       items:
	//         "$ref": "#/definitions/UserResponsee"
	//   '409':
	//     description: User already exists
	//     schema:
	//       type: array
	//       items:
	//         "$ref": "#/definitions/ModelError"
	//   '400':
	//     description: Bad Request
	//     schema:
	//       type: array
	//       items:
	//         "$ref": "#/definitions/ModelError"
	//   '500':
	//     description: Internal Server Error
	//     schema:
	//       type: array
	//       items:
	//         "$ref": "#/definitions/ModelError"
	config.Router.HandleFunc(UsersEndpoint, SetMiddlewareWithoutAuthentication(config.CreateUser)).Methods("POST")

	// swagger:operation POST /login Login
	//
	// Login
	//
	// Log in as an existing user.
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - in: body
	//   name: user
	//   description: User Information
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/User"
	// responses:
	//   '200':
	//     description: Token Information
	//     schema:
	//       type: array
	//       items:
	//         "$ref": "#/definitions/Login"
	//   '409':
	//     description: User does not exist
	//     schema:
	//       type: array
	//       items:
	//         "$ref": "#/definitions/ModelError"
	//   '400':
	//     description: Bad Request
	//     schema:
	//       type: array
	//       items:
	//         "$ref": "#/definitions/ModelError"
	//   '500':
	//     description: Internal Server Error
	//     schema:
	//       type: array
	//       items:
	//         "$ref": "#/definitions/ModelError"
	config.Router.HandleFunc(LoginEndpoint, SetMiddlewareWithoutAuthentication(config.Login)).Methods("POST")

	// swagger:operation POST /messages SendMessage
	//
	// Send a new message
	//
	// Send a message from one user to another.
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - in: body
	//   name: user
	//   description: Message Information
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/Message"
	// responses:
	//   '201':
	//     description: Message Information
	//     schema:
	//       type: array
	//       items:
	//         "$ref": "#/definitions/MessageResponse"
	//   '409':
	//     description: Sender or recipient doesn't exist
	//     schema:
	//       type: array
	//       items:
	//         "$ref": "#/definitions/ModelError"
	//   '400':
	//     description: Bad Request
	//     schema:
	//       type: array
	//       items:
	//         "$ref": "#/definitions/ModelError"
	//   '500':
	//     description: Internal Server Error
	//     schema:
	//       type: array
	//       items:
	//         "$ref": "#/definitions/ModelError"
	config.Router.HandleFunc(MessagesEndpoint, SetMiddlewareWithAuthentication(config.SendMessage)).Methods("POST")

	// swagger:operation GET /messages GetMessages
	//
	// Get Messages
	//
	// Fetch all existing messages to a given recipient, within a range of message IDs.
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: recipient
	//   in: query
	//   description: User ID of recipient.
	//   required: true
	//   type: string
	// - name: start
	//   in: query
	//   description: Starting message ID. Messages will be returned in increasing order of message ID, starting from this value (or the next lowest value stored in the database).
	//   required: true
	//   type: string
	// - name: limit
	//   in: query
	//   description: Limit the response to this many messages.
	//   required: false
	//   type: string
	// responses:
	//   '200':
	//     description: Messages
	//     schema:
	//       type: array
	//       items:
	//         "$ref": "#/definitions/Message"
	//   '404':
	//     description: Records not found
	//     schema:
	//       type: array
	//       items:
	//         "$ref": "#/definitions/ModelError"
	//   '400':
	//     description: Bad Request
	//     schema:
	//       type: array
	//       items:
	//         "$ref": "#/definitions/ModelError"
	//   '500':
	//     description: Internal Server Error
	//     schema:
	//       type: array
	//       items:
	//         "$ref": "#/definitions/ModelError"
	config.Router.HandleFunc(MessagesEndpoint, SetMiddlewareWithAuthentication(config.GetMessages)).Methods("GET")
}

func (serverConfiguration *Handler) InitDatabase(databasepath string) {
	dao, err := service.NewDAO(databasepath)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	serverConfiguration.Database = service.NewChatService(dao)
}

// notice that there are 2 different middlewares, since there are endpoints that requires no authentication
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
