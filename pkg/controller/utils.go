package controller

import (
	"errors"

	"github.com/challenge/pkg/models"
)

//validateUserData validates that the given input data is not empty
func validateUserData(user models.User) error {
	if user.Username == "" {
		return errors.New("username cannot be empty")
	}
	if user.Password == "" || len(user.Password) < 8 {
		return errors.New("password should be at least 8 characters")
	}
	return nil
}

func validateMessageData(message models.Message) error {
	if message.UserID == 0 {
		return errors.New("sender cannot be empty")
	}
	if message.Recipient == 0 {
		return errors.New("recipient cannot be empty")
	}

	return nil
}
