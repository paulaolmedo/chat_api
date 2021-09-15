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
func NewChatService(repository ChatRepository) ChatService {
	return &serviceProperties{repository}
}

func (service *serviceProperties) CreateUser(userInfo models.User) (*models.User, error) {
	hashUserPassword(&userInfo)

	err := service.repository.AddUser(&userInfo)
	if err != nil && err.Error() == "UNIQUE constraint failed: users.username" {
		return &models.User{}, errors.New("user already exists")
	} else if err != nil {
		return &models.User{}, err
	}

	return &userInfo, nil
}

func (service *serviceProperties) GetUser(user models.User) (models.User, error) {
	hashUserPassword(&user)

	userInfo, err := service.repository.GetUser(&user)
	if err != nil && err.Error() == "record not found" {
		return models.User{}, errors.New("user not found")
	} else if err != nil {
		return models.User{}, err
	}

	return userInfo, err
}

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

func (service *serviceProperties) GetMessages(filter models.MessageFilter) ([]models.Message, error) {
	msgs, err := service.repository.GetMessages(filter)
	if err != nil {
		return nil, err
	}

	return msgs, err
}
