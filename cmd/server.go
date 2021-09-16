// Chat API backend.
//
// Basic security and persistence solution for a chat API
//
//     Schemes: http
//     Host: 0.0.0.0:8080
//	   BasePath: /
//     Version: 1.0.0
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
// swagger:meta
package main

import (
	"log"
	"os"

	"github.com/challenge/pkg/controller"
	"github.com/magiconair/properties"
)

const (
	// app properties
	fileLocation              = "cmd/chat.properties"
	hostProperty              = "host"
	databaseFolderProperty    = "database_folder"
	databaseHostProperty      = "database_host"
	serverEnvironmentProperty = "environment"

	// possible environments
	development = "dev"
	production  = "production"
)

//go:generate swagger generate spec
func main() {
	// first, initialize the necessary properties to run the server
	p := properties.MustLoadFile(fileLocation, properties.UTF8)

	host := p.MustGetString(hostProperty)
	databaseFolder := p.MustGetString(databaseFolderProperty)
	databaseHost := p.MustGetString(databaseHostProperty)
	environment := p.MustGetString(serverEnvironmentProperty)

	if _, err := os.Stat(databaseFolder); os.IsNotExist(err) {
		err := os.Mkdir(databaseFolder, 0755)
		if err != nil {
			log.Fatalf("error creating database folder %v", err)
		}
	}

	configuration := controller.Handler{}

	switch environment {
	case development:
		configuration.SetEnvironment(false)
	case production:
		configuration.SetEnvironment(true)
	default:
		log.Fatalf("error reading environment value")
	}

	configuration.InitHTTPServer(databaseHost)
	configuration.Run(host)
}
