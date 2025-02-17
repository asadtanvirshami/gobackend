package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CommunityCourse struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	CommunityID uuid.UUID `gorm:"type:uuid;not null" json:"community_id"`
	CreatedByID uuid.UUID `gorm:"type:uuid;not null" json:"created_by_id"`
	Title       string    `gorm:"type:varchar(255);not null" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	VideoURL    string    `json:"video_url"`
	Thumbnail   string    `json:"thumbnail"`
	Duration    int       `json:"duration"` // In minutes

	CreatedBy User      `gorm:"foreignKey:CreatedByID;constraint:OnDelete:CASCADE;" json:"created_by"`
	Community Community `gorm:"foreignKey:CommunityID;constraint:OnDelete:CASCADE;" json:"community"`
}
