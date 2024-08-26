package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Title             string    `gorm:"size:100;not null" json:"title"`
	Author            string    `gorm:"size:100;not null" json:"author"`
	ISBN              string    `gorm:"size:100;not null" json:"isbn"`
	Publisher         string    `gorm:"size:100;not null" json:"publisher"`
	PublicationDate   string    `gorm:"size:100;not null" json:"publication_date"`
	Genre             string    `gorm:"size:100;not null" json:"genre"`
	Langage           string    `gorm:"size:100;not null" json:"language"`
	Pages             int       `gorm:"size:100;not null" json:"pages"`
	Edition           string    `gorm:"size:100;not null" json:"edition"`
	Quantity          int       `gorm:"size:100;not null" json:"quantity"`
	UUID              uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"uuid"`
	CreatedBy         string    `gorm:"size:100;not null" json:"created_by"`
	RemainingQuantity int    `gorm:"size:100;not null" json:"rem_quantity"`
	Store             string    `gorm:"size:100;not null" json:"store"`
}

func (book *Book) BeforeCreate(tx *gorm.DB) (err error) {
	book.UUID = uuid.New()
	return
}
