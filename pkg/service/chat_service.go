package service

import (
	"crypto/md5"
	"encoding/hex"
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

// hashUserPassword hashes the user password. For tests purposes it's only a MD5 hash, and without any salt
func hashUserPassword(user *models.User) {
	hash := md5.New()
	hash.Write([]byte(user.Password))
	user.Password = hex.EncodeToString(hash.Sum([]byte(nil)))
}
