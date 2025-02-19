package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Community Leaderboard
type CommunityLeaderboard struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	CommunityID uuid.UUID `gorm:"type:uuid;not null" json:"community_id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Points      int       `gorm:"not null" json:"points"`
	Rank        int       `gorm:"not null" json:"rank"`

	User      User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"users"`
	Community Community `gorm:"foreignKey:CommunityID;constraint:OnDelete:CASCADE;" json:"communities"`
}
