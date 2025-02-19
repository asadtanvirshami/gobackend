package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a registered user in the system
type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	FirstName string    `gorm:"type:varchar(50)" json:"first_name"`
	LastName  string    `gorm:"type:varchar(50)" json:"last_name"`
	Bio       *string   `gorm:"type:text" json:"bio,omitempty"`
	Image     *string   `gorm:"type:varchar(255)" json:"image,omitempty"`
	ImageID   *string   `gorm:"type:varchar(255)" json:"image_id,omitempty"`

	// Authentication & Security
	Email        string  `gorm:"type:varchar(100);unique;not null" json:"email"`
	Password     string  `gorm:"type:varchar(255);not null" json:"-"` // Omit in JSON response
	OTP          *string `gorm:"type:varchar(10)" json:"otp,omitempty"`
	OTPExpiresAt *int64  `gorm:"type:bigint" json:"otp_expires_at,omitempty"`
	Token        *string `gorm:"type:text" json:"token,omitempty"`
	Blocked      bool    `gorm:"default:false" json:"blocked"`

	// Personal Preferences
	Pronouns    *string `gorm:"type:varchar(50)" json:"pronouns,omitempty"`
	MyersBriggs *string `gorm:"type:varchar(10)" json:"myers_briggs,omitempty"`
	Location    *string `gorm:"type:varchar(100)" json:"location,omitempty"`
	Url         *string `gorm:"type:varchar(255)" json:"url,omitempty"`

	// Social Links
	Instagram *string `gorm:"type:varchar(255)" json:"instagram,omitempty"`
	LinkedIn  *string `gorm:"type:varchar(255)" json:"linkedin,omitempty"`
	X         *string `gorm:"type:varchar(255)" json:"x,omitempty"`
	YouTube   *string `gorm:"type:varchar(255)" json:"youtube,omitempty"`
	Facebook  *string `gorm:"type:varchar(255)" json:"facebook,omitempty"`
	Other     *string `gorm:"type:varchar(255)" json:"other,omitempty"`

	CreatedAt   int64          `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   int64          `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Communities []Community    `gorm:"foreignKey:UserID" json:"communities"`
}
