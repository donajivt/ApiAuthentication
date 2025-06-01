package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID          uuid.UUID `gorm:"type:uniqueidentifier;primaryKey" json:"id"`
	Email       string    `gorm:"uniqueIndex;size:128" json:"email"`
	Name        string    `gorm:"size:128" json:"name"`
	PhoneNumber string    `gorm:"size:32" json:"phoneNumber"`
	Password    string    `gorm:"size:255" json:"-"`
	Roles       []Role    `gorm:"many2many:user_roles" json:"roles"`
	gorm.Model
}
