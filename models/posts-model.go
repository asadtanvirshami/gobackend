package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CommunityPost struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	CommunityID uuid.UUID `gorm:"type:uuid;not null" json:"community_id"`
	AuthorID    uuid.UUID `gorm:"type:uuid;not null" json:"author_id"`
	Content     string    `gorm:"type:text" json:"content"`
	HasPoll     bool      `json:"has_poll"`

	Author    User           `gorm:"foreignKey:AuthorID;constraint:OnDelete:CASCADE;" json:"author"`
	Community Community      `gorm:"foreignKey:CommunityID;constraint:OnDelete:CASCADE;" json:"community"`
	Poll      *CommunityPoll `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE;" json:"poll"`
}
