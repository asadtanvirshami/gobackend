package models  

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)
type User struct {
	gorm.Model
	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	FirstName string 
	LastName string 
	Token string 
	Blocked bool
	Email string `gorm:"unique"`
	Password string 
}