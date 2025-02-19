package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Community struct {
	gorm.Model
	ID         uuid.UUID        `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name       string           `json:"name"`
	About      string           `json:"about"`
	Private    bool             `json:"private"`
	Cover      string           `json:"cover"`
	Instagram  string           `json:"instagram"`
	X          string           `json:"x"`
	Facebook   string           `json:"facebook"`
	Other      string           `json:"other"`
	Website    string           `json:"website"`
	CategoryID uuid.UUID        `gorm:"type:uuid;not null" json:"category_id"`
	Category   Category         `gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE;" json:"categories"`
	UserID     uuid.UUID        `gorm:"type:uuid;not null" json:"user_id"`
	User       User             `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"users"`
	Events     []CommunityEvent `gorm:"foreignKey:CommunityID" json:"events"`
}
