package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	FirstName   string
	LastName    string
	Bio         string
	Image       string
	OTP         string
	Pronouns    string
	MyersBriggs string
	Url         string
	Location    string
	Instagram   string
	LinkedIn    string
	X           string
	YouTube     string
	Facebook    string
	Other       string
	Token       string
	Blocked     bool
	Email       string `gorm:"unique"`
	Password    string
}
