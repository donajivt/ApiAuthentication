package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID          uuid.UUID      `gorm:"type:uniqueidentifier;primaryKey" json:"id"`
	Email       string         `gorm:"uniqueIndex;size:128" json:"email"`
	Name        string         `gorm:"size:128" json:"name"`
	PhoneNumber string         `gorm:"size:32" json:"phoneNumber"`
	Password    string         `gorm:"size:255" json:"-"`
	RoleID      uint           `json:"-"`
	Role        Role           `gorm:"foreignKey:RoleID" json:"role"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
