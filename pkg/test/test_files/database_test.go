package test

import (
	"log"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/challenge/pkg/controller"
	"github.com/challenge/pkg/models"
	"github.com/challenge/pkg/service"
	"github.com/magiconair/properties"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var chatService service.ChatService
var testServer *httptest.Server
var testHandler controller.Handler

const propertiesFile = "../test_properties/test.properties"

func TestMain(test *testing.M) {
	InitDatabase()
	InitService()
	os.Exit(test.Run())
}

func InitDatabase() {
	p := properties.MustLoadFile(propertiesFile, properties.UTF8)
	pathDatabase := p.MustGetString("pathDatabase")

	dao, err := service.NewDAO(pathDatabase)
	if err != nil {
		log.Fatalf("failed to connect to test database: %v", err)
		return
	}
	chatService = service.NewChatService(dao)
}

func InitService() {
	p := properties.MustLoadFile(propertiesFile, properties.UTF8)
	pathEndpoints := p.MustGetString("pathEndpoints")

	testHandler = controller.Handler{}
	testHandler.InitHTTPServer(pathEndpoints) // sólo inicializo la base de datos
	testServer = httptest.NewServer(testHandler.Router)
}

func InitUser(username string, password string) models.User {
	return models.User{Username: username, Password: password}
}

func TestAddUser(test *testing.T) {
	user := InitUser("testUser", "12345678")
	responseUser, errCreatingUser := chatService.CreateUser(user)

	require.NoError(test, errCreatingUser)

	expectedPassword := "25d55ad283aa400af464c76d713c07ad"

	assert.Equal(test, user.Username, responseUser.Username)
	assert.Equal(test, expectedPassword, responseUser.Password)
}

func TestAddUserThatAlreadyExists(test *testing.T) {
	user := InitUser("testUser", "12345678")
	responseUser, errCreatingUser := chatService.CreateUser(user)

	expectedError := "user already exists"
	assert.EqualError(test, errCreatingUser, expectedError)
	assert.Empty(test, responseUser)
}

func TestGetUser(test *testing.T) {
	user := InitUser("testUser", "12345678")
	responseUser, errSearchingUser := chatService.GetUser(user)

	require.NoError(test, errSearchingUser)

	expectedPassword := "25d55ad283aa400af464c76d713c07ad"

	assert.Equal(test, user.Username, responseUser.Username)
	assert.Equal(test, expectedPassword, responseUser.Password)
}

func TestGetNonExistentUser(test *testing.T) {
	user := InitUser("nonExistentUser", "12345678")
	responseUser, errSearchingUser := chatService.GetUser(user)

	expectedError := "user not found"

	assert.EqualError(test, errSearchingUser, expectedError)
	assert.Empty(test, responseUser)
}
