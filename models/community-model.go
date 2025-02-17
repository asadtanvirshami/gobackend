package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Community struct {
	gorm.Model
	ID          uuid.UUID        `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name        string           `json:"name"`
	About       string           `json:"about"`
	Private     bool             `json:"private"`
	Cover       string           `json:"cover"`
	Members     []CommunityUser  `gorm:"foreignKey:CommunityID" json:"members"`
	Events      []CommunityEvent `gorm:"foreignKey:CommunityID" json:"events"`
	CreatedByID uuid.UUID        `gorm:"type:uuid;not null" json:"created_by_id"`
	CreatedBy   User             `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"created_by"`
	Instagram   string           `json:"instagram"`
	X           string           `json:"x"`
	Facebook    string           `json:"facebook"`
	Other       string           `json:"other"`
	Website     string           `json:"website"`
}
