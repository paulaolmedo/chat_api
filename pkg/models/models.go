package models

import (
	"time"
)

// Database Structures

// User represents the user for this application
//
// swagger:model
type User struct {
	// User ID of the newly created user.
	//gorm.Model
	//
	// example: 1
	UserID int64 `json:"-" gorm:"primary_key"`
	// messages associated to the user
	Message []Message `json:"-" gorm:"foreignkey:UserID"`
	// user identification. must be unique
	//
	// example: test_user
	Username string `json:"username" gorm:"index:username,unique"`
	// user password
	//
	// example: password
	Password string `json:"password"`
}

// Message represents the main structure of the sent messages
//
// swagger:model
type Message struct {
	// messages identification. autogenerated by the database
	// gorm.Model
	//
	// example: 1
	MessageID int64 `json:"id" gorm:"primary_key"`
	// User ID of sender
	//
	// example: 1
	UserID int64 `json:"sender" gorm:"index:sender"`
	// User ID of recipient.
	//
	// example: 2
	Recipient int64 `json:"recipient" gorm:"index:recipient"`
	// Message content (one of three possible types)
	MessageContent Content `json:"content" gorm:"foreignKey:MessageID"`
	// messages timestamp in UTC to avoid time-zones problems
	Timestamp time.Time `json:"timestamp"`
}

// Content is a complement of the messages structure, since it contains all the possible types
//
// swagger:model
type Content struct {
	// content identification. autogenerated by the database
	// gorm.Model
	//
	// example: 1
	ContentID int64  `json:"id" gorm:"primary_key"`
	Text      Text   `json:"text,omitempty" gorm:"foreignKey:TextID"`
	Image     Image  `json:"image,omitempty" gorm:"foreignKey:ImageID"`
	Video     Video  `json:"video,omitempty" gorm:"foreignKey:VideoID"`
	Type      string `json:"type"`
}

// Image type
//
// swagger:model
type Image struct {
	ImageID int64  `json:"id,omitempty" gorm:"primary_key"`
	Url     string `json:"url,omitempty"`
	Height  int    `json:"height,omitempty"`
	Width   int    `json:"width,omitempty"`
}

// Video type
//
// swagger:model
type Video struct {
	VideoID int64  `json:"id,omitempty" gorm:"primary_key"`
	Url     string `json:"url,omitempty"`
	Source  string `json:"source,omitempty"`
}

// Text type
//
// swagger:model
type Text struct {
	TextID int64  `json:"id,omitempty" gorm:"primary_key"`
	Text   string `json:"text,omitempty"`
}

// Support structures
type Health struct {
	Status string `json:"health"`
}

// Login represents the information of the logged user
//
// swagger:model
type Login struct {
	// User ID of the user who logged in.
	//
	// example: 1
	Id int64 `json:"id"`
	// Authentication token to use for API calls on behalf of this user.
	//
	// example: Bearer token
	Token string `json:"token"`
}

// MessageFilter represents the filter to search messages
//
// swagger:model
type MessageFilter struct {
	// Starting message ID. Messages will be returned in increasing order of message ID, starting from this value (or the next lowest value stored in the database).
	//
	// example: 1
	Start int64 `json:"start"`
	// User ID of recipient.
	//
	// example: 1
	Recipient int64 `json:"recipient"`
	// Limit the response to this many messages.
	//
	// example: 100
	Limit int `json:"limit"`
}

// MessageResponse represents the response when inserting a new message
//
// swagger:model
type MessageResponse struct {
	// Message IDs are required to be unique and increase over time; they may or may not be sequential.
	//
	// example: 1
	MessageID int64 `json:"id"`
	// Timestamp for this message, as recorded on the server.
	// example "2019-09-03T19:54:22Z"
	Timestamp time.Time `json:"timestamp"`
}

// ModelError
//
// swagger:model
type ModelError struct {
	// Error code
	Code string `json:"code"`
	//
	Message string `json:"message"`
}

// User represents identification of the created user
//
// swagger:model
type UserResponse struct {
	// User ID of the newly created user.
	//gorm.Model
	//
	// example: 1
	UserID int64
}
