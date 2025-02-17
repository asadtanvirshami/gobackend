package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CommunityPoll struct {
	gorm.Model
	ID      uuid.UUID             `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	PostID  uuid.UUID             `gorm:"type:uuid;not null" json:"post_id"`
	Options []CommunityPollOption `gorm:"foreignKey:PollID" json:"options"`
}
