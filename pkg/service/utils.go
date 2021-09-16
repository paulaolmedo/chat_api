package service

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/challenge/pkg/models"
)

const (
	// Error messages
	missingRecord    = "record not found"
	missingUser      = "user not found"
	missingSender    = "sender does not exist"
	missingRecipient = "recipient does not exist"
	userExists       = "user already exists"
)

// hashUserPassword hashes the user password. For tests purposes it's only a MD5 hash, and without any salt
func hashUserPassword(user *models.User) {
	hash := md5.New()
	hash.Write([]byte(user.Password))
	user.Password = hex.EncodeToString(hash.Sum([]byte(nil)))
}
