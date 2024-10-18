package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name            string
	Email           string `gorm:"unique"`
	Address         string
	UserType        string
	PasswordHash    string
	ProfileHeadline string
	Profile         Profile
	Applications    []Application `gorm:"foreignKey:UserID"`
}

type Profile struct {
	gorm.Model
	UserID     uint
	ResumeFile string
	Skills     string
	Education  string
	Experience string
}
