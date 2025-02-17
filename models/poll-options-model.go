package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CommunityPollOption struct {
	gorm.Model
	ID     uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	PollID uuid.UUID `gorm:"type:uuid;not null" json:"poll_id"`
	Option string    `gorm:"type:varchar(255);not null" json:"option"`
	Votes  int       `gorm:"default:0" json:"votes"`

	Poll CommunityPoll `gorm:"foreignKey:PollID;constraint:OnDelete:CASCADE;" json:"poll"`
}
