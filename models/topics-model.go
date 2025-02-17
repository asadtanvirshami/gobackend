package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Community Leaderboard
type CommunityTopics struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	CommunityID uuid.UUID `gorm:"type:uuid;not null" json:"community_id"`
	Points      int       `gorm:"not null" json:"points"`
	Community   Community `gorm:"foreignKey:CommunityID;constraint:OnDelete:CASCADE;" json:"community"`
}
