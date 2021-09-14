package service

import (
	"errors"

	"github.com/challenge/pkg/models"
)

//Objeto que contiene a la interfaz propia de la base de datos
type serviceProperties struct {
	repository ChatRepository
}

type ChatService interface {
	CreateUser(userInfo models.User) (*models.User, error)
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
