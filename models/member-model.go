package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CommunityUser struct {
	gorm.Model
	UserID      uuid.UUID `gorm:"type:uuid;primaryKey" json:"user_id"`
	CommunityID uuid.UUID `gorm:"type:uuid;primaryKey" json:"community_id"`
	Role        string    `gorm:"type:varchar(50);not null" json:"role"` // e.g., "admin", "moderator", "member"

	User      User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"user"`
	Community Community `gorm:"foreignKey:CommunityID;constraint:OnDelete:CASCADE;" json:"community"`
}
