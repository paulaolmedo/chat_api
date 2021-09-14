package models

import (
	"time"
)

type Health struct {
	Status string `json:"health"`
}

type Login struct {
	// TODO: Implement Login model
}

type User struct {
	UserID   int64     `json:"-" gorm:"primary_key"`
	Message  []Message `json:"-" gorm:"foreignkey:UserID"`
	Username string    `json:"username" gorm:"index:username,unique"`
	Password string    `json:"password"`
}

type Message struct {
	MessageID      int64     `json:"id" gorm:"primary_key"`
	UserID         int64     `json:"sender" gorm:"index:sender"`
	Recipient      int64     `json:"recipient" gorm:"index:recipient"`
	MessageContent Content   `json:"content" gorm:"foreignKey:MessageID"`
	Timestamp      time.Time `json:"-"`
}

type MessageFilter struct {
	Start     int64 `json:"start"`
	Recipient int64 `json:"recipient"`
	Limit     int   `json:"limit"`
}

type MessageResponse struct {
	MessageID int64     `json:"id"`
	Timestamp time.Time `json:"timestamp"`
}

type Content struct {
	ContentID int    `json:"id" gorm:"primary_key"`
	Text      Text   `json:"text,omitempty" gorm:"foreignKey:ContentID"`
	Image     Image  `json:"image,omitempty" gorm:"foreignKey:ContentID"`
	Video     Video  `json:"video,omitempty" gorm:"foreignKey:ContentID"`
	Type      string `json:"type"`
}

type Image struct {
	ImageID int    `json:"id,omitempty" gorm:"primary_key"`
	Url     string `json:"url,omitempty"`
	Height  int    `json:"height,omitempty"`
	Width   int    `json:"width,omitempty"`
}

type Video struct {
	VideoID int    `json:"id,omitempty" gorm:"primary_key"`
	Url     string `json:"url,omitempty"`
	Source  string `json:"source,omitempty"`
}

type Text struct {
	TextID int    `json:"id,omitempty" gorm:"primary_key"`
	Text   string `json:"text,omitempty"`
}
