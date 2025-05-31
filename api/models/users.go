package models

import (
	"time"
)

type User struct {
	Id            	uint    `gorm:"primaryKey;autoIncrement"`
	Name          	string  `gorm:"size:25"`
	Email         	string  `gorm:"unique;size:25"`
	Password      	string  `gorm:"type:text"`
	Type			string 	`gorm:"size:25;default:'USER'"`
	Photo 			*string `gorm:"size:25"`
	CreatedAt 		*time.Time
	UpdatedAt 		*time.Time
}
