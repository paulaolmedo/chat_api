package service

import (
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
}

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
