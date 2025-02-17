package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CommunityEvent struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	CommunityID uuid.UUID `gorm:"type:uuid;not null" json:"community_id"`
	EventName   string    `gorm:"type:varchar(255);not null" json:"event_name"`
	Description string    `gorm:"type:text" json:"description"`
	EventDate   time.Time `gorm:"not null" json:"event_date"`
	Cover       string    `json:"cover"`

	Community Community `gorm:"foreignKey:CommunityID;constraint:OnDelete:CASCADE;" json:"community"`
}
