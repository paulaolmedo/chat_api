package main

import (
	"log"
	"os"

	"github.com/challenge/pkg/controller"
	"github.com/magiconair/properties"
)

const (
	ServerPort       = "8080"
	CheckEndpoint    = "/check"
	UsersEndpoint    = "/users"
	LoginEndpoint    = "/login"
	MessagesEndpoint = "/messages"
)

func main() {
	p := properties.MustLoadFile("cmd/chat.properties", properties.UTF8)

	host := p.MustGetString("host")
	databaseFolder := p.MustGetString("database_folder")
	databaseHost := p.MustGetString("database_host")
	environment := p.MustGetString("environment")

	if _, err := os.Stat(databaseFolder); os.IsNotExist(err) {
		err := os.Mkdir(databaseFolder, 0755)
		if err != nil {
			log.Fatalf("error creating database folder %v", err)
		}
	}

	configuration := controller.Handler{}

	switch environment {
	case "dev":
		configuration.SetEnvironment(false)
	case "production":
		configuration.SetEnvironment(true)
	default:
		log.Fatalf("error reading environment value")
	}

	configuration.InitHTTPServer(databaseHost)
	configuration.Run(host)
}
