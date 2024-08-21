package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Store struct {
	gorm.Model

	Name                     string    `gorm:"size:100;not null" json:"name"`
	Location                 string    `gorm:"size:100;not null" json:"location"`
	City                     string    `gorm:"size:100;not null" json:"city"`
	PhoneNumber              string    `gorm:"size:100;not null" json:"phone_number"` // Changed to match JSON field "phone"
	Email                    string    `gorm:"size:100;not null" json:"email"`
	ContactPerson            string    `gorm:"size:100;not null" json:"contact_person"`
	ContactPersonEmail       string    `gorm:"size:100;not null" json:"contact_person_email"`
	ContactPersonPhoneNumber string    `gorm:"size:100;not null" json:"contact_person_phone"` // Changed to match JSON field "contact_person_phone"
	Description              string    `gorm:"size:100;not null" json:"description"`
	UUID                     uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"uuid"`
	CreatedBy                string    `gorm:"size:100;not null" json:"created_by"`
}

func (store *Store) BeforeCreate(tx *gorm.DB) (err error) {
	store.UUID = uuid.New()
	return
}
