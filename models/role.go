package models

type Role struct {
	ID   uint   `gorm:"primaryKey" json:"-"`
	Name string `gorm:"uniqueIndex;size:64" json:"name"`
}
