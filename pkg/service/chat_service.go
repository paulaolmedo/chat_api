package service

import (
	"errors"
	"time"

	"github.com/challenge/pkg/models"
)

//Objeto que contiene a la interfaz propia de la base de datos
type serviceProperties struct {
	repository ChatRepository
}

type ChatService interface {
	CreateUser(userInfo models.User) (*models.User, error)
	GetUser(user models.User) (models.User, error)
	AddMessage(message *models.Message) (models.MessageResponse, error)
	GetMessages(filter models.MessageFilter) ([]models.Message, error)
}

// NewChatService initializes the service that communicates with the database
// El propósito de esta "clase" es realizar operaciones que no se supone que haga la base de datos
func NewChatService(repository ChatRepository) ChatService {
	return &serviceProperties{repository}
}

//CreateUser adds a new user
func (service *serviceProperties) CreateUser(userInfo models.User) (*models.User, error) {
	hashUserPassword(&userInfo)

	err := service.repository.AddUser(&userInfo)
	if err != nil && err.Error() == "UNIQUE constraint failed: users.username" {
		return &models.User{}, errors.New(userExists)
	} else if err != nil {
		return &models.User{}, err
	}

	return &userInfo, nil
}

//GetUser verify if the pair username/password belongs to an user
func (service *serviceProperties) GetUser(user models.User) (models.User, error) {
	hashUserPassword(&user)

	userInfo, err := service.repository.GetUser(&user)
	if err != nil && err.Error() == missingRecord {
		return models.User{}, errors.New(missingUser)
	} else if err != nil {
		return models.User{}, err
	}

	return userInfo, err
}

//AddMessage adds a new message after making sure that the sender and the recipients are valid users
func (service *serviceProperties) AddMessage(message *models.Message) (models.MessageResponse, error) {
	timestamp := time.Now().UTC()
	message.Timestamp = timestamp

	err := service.repository.AddMessage(message)
	if err != nil {
		return models.MessageResponse{}, err
	}
	response := models.MessageResponse{MessageID: message.MessageID, Timestamp: timestamp}
	return response, nil
}

// GetMessages given a starting id and a recipient, retrieves a certain amount of messages (default: 100)
func (service *serviceProperties) GetMessages(filter models.MessageFilter) ([]models.Message, error) {
	msgs, err := service.repository.GetMessages(filter)
	if err != nil {
		return nil, err
	}

	return msgs, err
}
