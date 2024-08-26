package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Client struct {
	gorm.Model
	Name      string    `gorm:"size:100;not null" json:"name"`
	UUID      uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"uuid"`
	CreatedBy string    `gorm:"size:100;not null" json:"created_by"`
	Phone     string    `gorm:"size:100;not null" json:"phone"`
	Email     string    `gorm:"size:100;not null" json:"email"`
	Store     string    `gorm:"size:100;not null" json:"store"`
	City      string    `gorm:"size:100;not null" json:"city"`
}

func (client *Client) BeforeCreate(tx *gorm.DB) (err error) {
	client.UUID = uuid.New()
	return
}
