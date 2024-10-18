package models

import (
	"gorm.io/gorm"
)

type Job struct {
	gorm.Model
	Title             string
	Description       string
	PostedOn          string
	TotalApplications int
	CompanyName       string
	UserID            uint
	PostedBy          User          `gorm:"foreignKey:UserID"`
	Applications      []Application `gorm:"foreignKey:JobID"`
}

type Application struct {
	gorm.Model
	UserID uint
	JobID  uint
	User   User `gorm:"foreignKey:UserID"`
	Job    Job  `gorm:"foreignKey:JobID"`
}
