package service

import (
	"errors"
	"log"

	"github.com/challenge/pkg/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type ChatDAO struct {
	db *gorm.DB
}

// ChatRepository funciones mínimas para utilizar la base de datos
type ChatRepository interface {
	AddUser(user *models.User) error
	GetUser(user *models.User) (models.User, error)
	AddMessage(message *models.Message) error
	checkUserExistence(userID int64) bool
}

// NewDAO given a path where to store the database, initializes the DAO with the API models
func NewDAO(databasePath string) (*ChatDAO, error) {
	db, err := gorm.Open(sqlite.Open(databasePath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	var args []interface{}
	args = append(args, &models.User{})
	args = append(args, &models.Message{})
	args = append(args, &models.Content{})
	args = append(args, &models.Text{})
	args = append(args, &models.Image{})
	args = append(args, &models.Video{})

	err = db.AutoMigrate(args...)
	if err != nil {
		log.Fatalf("error migrating database models %v", err)
		log.Fatal(err)
		return nil, err
	}

	dao := ChatDAO{db}

	return &dao, nil
}

//AddUser adds a new user
func (dao *ChatDAO) AddUser(user *models.User) error {
	currentConnection := dao.db

	result := currentConnection.Create(&user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

//GetUser verify if the pair username/password belongs to an user
func (dao *ChatDAO) GetUser(user *models.User) (models.User, error) {
	currentConnection := dao.db

	query := "username = ? AND password = ?"
	var args []interface{}
	args = append(args, user.Username)
	args = append(args, user.Password)

	result := currentConnection.Where(query, args...).Find(&user)
	if result.Error != nil {
		return models.User{}, result.Error
	} else if result.RowsAffected == 0 {
		return models.User{}, errors.New("record not found")

	}
	return *user, nil
}

func (dao *ChatDAO) AddMessage(message *models.Message) error {
	checkSender := dao.checkUserExistence(message.UserID)
	if !checkSender {
		return errors.New("sender does not exist")
	}

	checkRecipient := dao.checkUserExistence(message.Recipient)
	if !checkRecipient {
		return errors.New("recipient does not exist")
	}

	currentConnection := dao.db
	result := currentConnection.Create(&message)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (dao *ChatDAO) checkUserExistence(userID int64) bool {
	currentConnection := dao.db
	var user models.User

	result := currentConnection.First(&user, userID)

	return result.Error == nil
}